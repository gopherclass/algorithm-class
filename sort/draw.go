package main

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type drawingData struct {
	sortName  string
	inputName string
	samples   []sortStat
}

func newDrawingData(record benchmarkRecord) *drawingData {
	return &drawingData{
		sortName:  record.sortName,
		inputName: record.inputName,
		samples:   record.sizedStats,
	}
}

type drawing struct {
	imageWidth         vg.Length
	imageHeight        vg.Length
	titleMaxsize       vg.Length
	xLabelSize         vg.Length
	yLabelSize         vg.Length
	importantLineStyle draw.LineStyle
	auxiliaryLineStyle draw.LineStyle
}

func documentDrawing() *drawing {
	return &drawing{
		imageWidth:   130 * vg.Millimeter,
		imageHeight:  46.4 * vg.Millimeter,
		titleMaxsize: 6 * vg.Millimeter,
		xLabelSize:   6 * vg.Millimeter,
		yLabelSize:   6 * vg.Millimeter,
		importantLineStyle: draw.LineStyle{
			Width: vg.Points(2),
			Color: plotutil.Color(0),
		},
		auxiliaryLineStyle: draw.LineStyle{
			Width: vg.Points(1),
		},
	}
}

// func largeDrawing() *drawing {
// 	return &drawing{
// 		imageWidth:  imageWidth,
// 		imageHeight: imageHeight,
// 	}
// }

func (w *drawing) storeRecord(record benchmarkRecord) error {
	// TODO: storeRecord() 함수는 설계가 나쁘다.
	data := newDrawingData(record)
	var err0 error
	setErr := func(err error) {
		if err0 != nil {
			return
		}
		err0 = err
	}
	draw := func(server serveY) {
		pl, err := w.drawSamples(data, serveCompare{})
		setErr(err)
		if err == nil {
			name := fmt.Sprintf("%s - %s - %s.jpg",
				record.sortName,
				record.inputName,
				server.label())
			setErr(w.storePlot(pl, name))
		}
	}
	draw(serveCompare{})
	draw(serveCompare{})
	draw(serveAccess{})
	draw(serveMicrosecondLapse{})
	draw(&serveTotal{
		compareWeight: 1.0,
		swapWeight:    3.0,
		accessWeight:  1.0,
	})
	return err0
}

func (w *drawing) storePlot(pl *plot.Plot, name string) error {
	return pl.Save(w.imageWidth, w.imageHeight, name)
}

func (w *drawing) drawSamples(data *drawingData, server serveY) (*plot.Plot, error) {
	pl, err := plot.New()
	if err != nil {
		return nil, err
	}

	// TODO: size scaling
	set := serveSet(server, data.samples)
	line, err := plotter.NewLine(set)
	if err != nil {
		return nil, err
	}
	line.LineStyle = w.importantLineStyle

	title := fmt.Sprintf("%s - %s", data.sortName, data.inputName)
	w.setTitle(pl, title)

	pl.X.Label.Text = "Size"
	pl.X.Label.Font.Size = w.xLabelSize
	pl.Y.Label.Text = server.label()
	pl.Y.Label.Font.Size = w.yLabelSize

	pl.Add(line)
	pl.Legend.Add(server.label(), line)

	w.drawFunction(pl, "y = x", identifyFunc, plotutil.Color(1))
	w.drawFunction(pl, "y = x²", quadraticFunc, plotutil.Color(2))
	w.drawFunction(pl, "y = xlog(x)", xlogxFunc, plotutil.Color(3))
	return pl, nil
}

func (w *drawing) setTitle(pl *plot.Plot, title string) {
	pl.Title.Font.Size = w.titleMaxsize
	for {
		size := pl.Title.Font.Width(title)
		if size <= w.imageWidth {
			return
		}
		pl.Title.Font.Size -= 0.5
	}
}

func (w *drawing) drawFunction(pl *plot.Plot, legend string, fn func(float64) float64, c color.Color) {
	v := plotter.NewFunction(fn)
	v.LineStyle = w.auxiliaryLineStyle
	v.LineStyle.Color = c
	pl.Add(v)
	pl.Legend.Add(legend, v)
}

type serveY interface {
	serveY(sample sortStat) float64
	label() string
}

func serveSet(server serveY, samples []sortStat) plotter.XYs {
	set := make(plotter.XYs, len(samples))
	for i, sample := range samples {
		set[i].X = float64(i)
		set[i].Y = server.serveY(sample)
	}
	return set
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

type serveTotal struct {
	compareWeight float64
	swapWeight    float64
	accessWeight  float64
}

func (serveTotal) label() string { return "total" }

func (w *serveTotal) serveY(s sortStat) float64 {
	tot := w.compareWeight * s.averageLess
	tot += w.swapWeight * s.averageSwap
	tot += w.accessWeight * s.averagePeek
	return tot
}

func identifyFunc(x float64) float64 { return x }

func quadraticFunc(x float64) float64 { return x * x }

func xlogxFunc(x float64) float64 {
	if x <= 1 {
		return 0
	}
	return x * math.Log2(x)
}
