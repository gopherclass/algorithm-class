package main

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type serveY interface {
	serveY(sample sortStat) float64
	label() string
}

func serveXYs(server serveY, samples []sortStat) plotter.XYs {
	set := make(plotter.XYs, len(samples))
	for i, sample := range samples {
		set[i].X = float64(i)
		set[i].Y = server.serveY(sample)
	}
	return set
}

type drawingLine struct {
	label string
	line  *plotter.Line
}

func (l drawingLine) draw(pl *plot.Plot) {
	pl.Add(l.line)
	pl.Legend.Add(l.label, l.line)
}

func serveLine(server serveY, samples []sortStat) (*plotter.Line, error) {
	return plotter.NewLine(serveXYs(server, samples))
}

type drawingData struct {
	label   string
	samples []sortStat
}

type drawingStyle struct {
	imageWidth         vg.Length
	imageHeight        vg.Length
	titleSize          vg.Length
	xLabelSize         vg.Length
	yLabelSize         vg.Length
	importantLineStyle draw.LineStyle
	auxiliaryLineStyle draw.LineStyle
}

func (w *drawingStyle) drawPlot(title string, server serveY, all ...drawingData) (*plot.Plot, error) {
	pl, err := plot.New()
	if err != nil {
		return nil, err
	}
	w.setTitle(pl, title, w.titleSize)
	pl.X.Label.Text = "Size"
	pl.X.Label.TextStyle.Font.Size = w.xLabelSize
	pl.Y.Label.Text = strings.Title(server.label())
	pl.Y.Label.TextStyle.Font.Size = w.yLabelSize

	var pal palette
	drawFunc := func(label string, fx func(float64) float64) {
		fn := plotter.NewFunction(fx)
		fn.LineStyle.Width = vg.Points(1)
		fn.LineStyle.Color = pal.color()
		pl.Add(fn)
		pl.Legend.Add(label, fn)
	}
	drawFunc("y = x", identifyFunc)
	drawFunc("y = x²", quadraticFunc)
	drawFunc("y = xlog(x)", xlogxFunc)

	for _, data := range all {
		line, err := serveLine(server, data.samples)
		if err != nil {
			return nil, err
		}
		line.LineStyle.Width = vg.Points(3)
		line.LineStyle.Color = pal.color()
		pl.Add(line)
		pl.Legend.Add(data.label, line)
	}
	return pl, nil
}

type palette struct {
	i int
}

func (p *palette) color() color.Color {
	c := plotutil.Color(p.i)
	p.i++
	return c
}

func (w *drawingStyle) setTitle(pl *plot.Plot, title string, size vg.Length) {
	v := &pl.Title
	v.Text = title
	v.Font.Size = size
	v.Font = calibrateFontWidth(w.imageWidth, title, v.Font)
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

func (w *drawingStyle) savePlot(pl *plot.Plot, path string) error {
	return pl.Save(w.imageWidth, w.imageHeight, path)
}

func (w *drawingStyle) drawRecord(r benchmarkRecord, server serveY) (*plot.Plot, error) {
	title := fmt.Sprintf("%s - %s", r.sortName, r.inputName)
	data := drawingData{
		label:   server.label(),
		samples: r.samples,
	}
	return w.drawPlot(title, server, data)
}

func (w *drawingStyle) saveResult(res benchmarkResult) error {
	os.Mkdir(res.sortName, os.ModePerm)
	for _, r := range res.records {
		for _, server := range []serveY{
			serveCompare{},
			serveCompare{},
			serveAccess{},
			serveMicrosecondLapse{},
		} {
			pl, err := w.drawRecord(r, server)
			if err != nil {
				return err
			}
			name := fmt.Sprintf("%s - %s.jpeg", r.inputName, server.label())
			err = w.savePlot(pl, filepath.Join(res.sortName, name))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type serveCompare struct{}

func (serveCompare) label() string {
	return "compare"
}

func (serveCompare) serveY(sample sortStat) float64 {
	return sample.averageLess
}

type serveSwap struct{}

func (serveSwap) label() string {
	return "swap"
}

func (serveSwap) serveY(sample sortStat) float64 {
	return sample.averageSwap
}

type serveAccess struct{}

func (serveAccess) label() string {
	return "access"
}

func (serveAccess) serveY(sample sortStat) float64 {
	return sample.averagePeek
}

type serveMicrosecondLapse struct{}

func (serveMicrosecondLapse) label() string {
	return "time"
}

func (serveMicrosecondLapse) serveY(sample sortStat) float64 {
	return float64(convMicroseconds(sample.averageLapse))
}

func convNanoseconds(d time.Duration) int64  { return int64(d) }
func convMicroseconds(d time.Duration) int64 { return int64(d) / 1e3 }
func convMilliseconds(d time.Duration) int64 { return int64(d) / 1e6 }

func identifyFunc(x float64) float64 { return x }

func quadraticFunc(x float64) float64 { return x * x }

func xlogxFunc(x float64) float64 {
	if x <= 1 {
		return 0
	}
	return x * math.Log2(x)
}

func documentStyle() *drawingStyle {
	return &drawingStyle{
		imageWidth:  130 * vg.Millimeter,
		imageHeight: 46.4 * vg.Millimeter,
		titleSize:   6 * vg.Millimeter,
		xLabelSize:  6 * vg.Millimeter,
		yLabelSize:  6 * vg.Millimeter,
	}
}
