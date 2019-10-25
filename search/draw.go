package main

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const imageFormat = "jpeg"

// 160 x 240
var defaultStyle = &presentationStyle{
	imageWidth:    160 * vg.Millimeter,
	imageHeight:   113 * vg.Millimeter,
	titleSize:     8 * vg.Millimeter,
	xAxisTextSize: 5.5 * vg.Millimeter,
	yAxisTextSize: 5.5 * vg.Millimeter,
}

type DrawRunner interface {
	Runner
	Illusts() []Illust
}

type Illust interface {
	Fx(perfStat) float64
	Legend(perfClass) string
	Tag() string
}

func toXYs(cls perfClass, il Illust) plotter.XYs {
	if cls.Size() == 0 {
		return nil
	}
	xylist := make(plotter.XYs, cls.Size())
	for i, perf := range cls.stats {
		xylist[i] = plotter.XY{
			X: float64(i),
			Y: il.Fx(perf),
		}
	}
	return xylist
}

type presentation struct {
	title     string
	xAxisText string
	yAxisText string
}

type presentationStyle struct {
	imageWidth    vg.Length
	imageHeight   vg.Length
	titleSize     vg.Length
	xAxisTextSize vg.Length
	yAxisTextSize vg.Length
}

func presentPlot(pr presentation, st *presentationStyle) (*plot.Plot, error) {
	pl, err := plot.New()
	if err != nil {
		return nil, err
	}
	pl.Title.Text = pr.title
	pl.Title.Font.Size = st.titleSize
	pl.Title.Font = calibrateFontWidth(st.imageWidth, pr.title, pl.Title.Font)
	pl.X.Label.Text = "Size"
	pl.Y.Label.Text = pr.yAxisText
	pl.X.Label.Font.Size = st.xAxisTextSize
	pl.Y.Label.Font.Size = st.yAxisTextSize
	return pl, nil
}

func calibrateFontWidth(width vg.Length, text string, font vg.Font) vg.Font {
	font.Size = binarySearchLength(0, font.Size, func(size vg.Length) bool {
		font.Size = size
		return font.Width(text) <= width
	})
	return font
}

func binarySearchLength(a, b vg.Length, pred func(vg.Length) bool) vg.Length {
	if pred(b) {
		return b
	}
	const e = 1e-4
	for b-a > e {
		c := (a + b) / 2
		if pred(c) {
			a = c
		} else {
			b = c
		}
	}
	if pred(b) {
		return b
	}
	return a
}

type perfdraw struct {
	perfcls perfClass
	illust  Illust
}

type drawLine struct {
	pl           *plot.Plot
	currentColor int
}

func newDrawLine(pl *plot.Plot) *drawLine {
	return &drawLine{pl: pl}
}

func (v *drawLine) nextColor() color.Color {
	c := plotutil.Color(v.currentColor)
	v.currentColor++
	return c
}

func (v *drawLine) draw(draw perfdraw) error {
	xylist := toXYs(draw.perfcls, draw.illust)
	line, err := plotter.NewLine(xylist)
	if err != nil {
		return err
	}
	line.Width = vg.Points(2)
	line.Color = v.nextColor()
	pl := v.pl
	pl.Add(line)
	legend := draw.illust.Legend(draw.perfcls)
	if legend != "" {
		pl.Legend.Add(legend, line)
	}
	return nil
}

func (v *drawLine) drawAux() {
	pl := v.pl
	drawfx := func(legend string, fx func(float64) float64) {
		o := plotter.NewFunction(fx)
		o.Width = vg.Points(1)
		o.Color = v.nextColor()
		pl.Add(o)
		pl.Legend.Add(legend, o)
	}
	drawfx("y = x", identityFunc)
	drawfx("y = xlog(x)", xlogxFunc)
	drawfx("y = x²", quadraticFunc)
	// drawfx("y = 2xlog(x)", bigO(2, xlogxFunc))
}

func (v *drawLine) plot() *plot.Plot {
	return v.pl
}

func identityFunc(x float64) float64 { return x }

func quadraticFunc(x float64) float64 { return x * x }

func xlogxFunc(x float64) float64 {
	if x <= 1 {
		return 0.0
	}
	return x * math.Log2(x)
}

func bigO(alpha float64, fx func(float64) float64) func(float64) float64 {
	return func(x float64) float64 {
		return alpha * fx(x)
	}
}

type classTag struct {
	inputClass string
	tag        string
}

type drawRunner struct {
	st       *presentationStyle
	size     int
	runners  []DrawRunner
	classMap map[classTag]*drawLine
}

func newDrawRunner(st *presentationStyle, size int) *drawRunner {
	return &drawRunner{
		st:       st,
		size:     size,
		classMap: make(map[classTag]*drawLine),
	}
}

type drawOptions struct {
	scale plot.Normalizer
}

func (o *drawRunner) draw(runner DrawRunner, iteration uint, options *drawOptions) error {
	o.runners = append(o.runners, runner)
	for _, cls := range timeitAll(runner, o.size, iteration) {
		for _, il := range runner.Illusts() {
			err := o.drawClass(runner, cls, il, options)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (o *drawRunner) drawClass(runner Runner, cls perfClass, il Illust, options *drawOptions) error {
	tag := il.Tag()
	key := classTag{cls.inputClass, tag}
	draw := o.classMap[key]
	if draw == nil {
		var err error
		draw, err = o.newPlot(cls, tag, options)
		if err != nil {
			return err
		}
		o.classMap[key] = draw
	}
	return draw.draw(perfdraw{
		perfcls: cls,
		illust:  il,
	})
}

func (o *drawRunner) newPlot(cls perfClass, tag string, options *drawOptions) (*drawLine, error) {
	pr := presentation{
		title:     presentationTitle(cls.inputClass, tag),
		xAxisText: "Size",
		yAxisText: tag,
	}
	pl, err := presentPlot(pr, o.st)
	if err != nil {
		return nil, err
	}
	if options != nil && options.scale != nil {
		pl.Y.Scale = options.scale
	}
	return newDrawLine(pl), nil
}

func (o *drawRunner) drawAux() {
	for _, draw := range o.classMap {
		draw.drawAux()
	}
}

func (o *drawRunner) setScale(scale plot.Normalizer) {
	for _, draw := range o.classMap {
		pl := draw.plot()
		pl.Y.Scale = scale
	}
}

func (o *drawRunner) store() error {
	dir := directoryName(o.runners)
	os.Mkdir(dir, os.ModePerm)

	for classTag, draw := range o.classMap {
		name := presentationName(classTag.inputClass, classTag.tag)
		path := filepath.Join(dir, name)
		err := draw.pl.Save(o.st.imageWidth, o.st.imageHeight, path)
		if err != nil {
			return err
		}
	}
	return nil
}

func presentationTitle(inputClass, tag string) string {
	return fmt.Sprintf("%s - %s", inputClass, tag)
}

func presentationName(inputClass, tag string) string {
	return fmt.Sprintf("%s - %s.%s", inputClass, tag, imageFormat)
}

func directoryName(runners []DrawRunner) string {
	s := make([]string, len(runners))
	for i, runner := range runners {
		s[i] = runner.Name()
	}
	return strings.Join(s, ", ")
}

type squareScale struct{}

func (squareScale) Normalize(min, max, x float64) float64 {
	a := (x - min) / (max - min)
	b := (x + min) / (max + min)
	return a * b
}

type logScale struct{}

func (logScale) Normalize(min, max, x float64) float64 {
	if min <= 0 || max <= 0 || x <= 0 {
		return 0.0
	}
	logMin := math.Log(min)
	return (math.Log(x) - logMin) / (math.Log(max) - logMin)
}
