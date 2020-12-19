package tchart

import (
	"errors"

	"github.com/s-westphal/termui/v3"
	"github.com/s-westphal/termui/v3/widgets"
)

// ScatterPlotWidget line chart widget
type ScatterPlotWidget struct {
	ScatterPlot *widgets.Plot
	storages    []*Storage
	width       int
	height      int
}

// NewScatterPlotWidget create line chart widget
func NewScatterPlotWidget(title string, storages []*Storage) *ScatterPlotWidget {
	scatterPlot := widgets.NewPlot()
	scatterPlot.Title = title
	scatterPlot.PlotType = widgets.ScatterPlot
	scatterPlot.Marker = widgets.MarkerDot
	// scatterPlot.DotMarkerRune = 'x'

	sp := &ScatterPlotWidget{
		ScatterPlot: scatterPlot,
		storages:    storages,
		width:       0,
		height:      0,
	}
	return sp
}

func (sp *ScatterPlotWidget) setDimension(x, y, w, h int) {
	sp.width = w
	sp.height = h
	sp.ScatterPlot.SetRect(x, y, x+w, y+h)
}

// Buffer buffer for rendering
func (sp *ScatterPlotWidget) Buffer() (*termui.Buffer, error) {
	if sp.storages[0].Data.length > 1 {
		buffer := termui.NewBuffer(sp.ScatterPlot.GetRect())

		sp.ScatterPlot.Lock()
		sp.ScatterPlot.Draw(buffer)
		sp.ScatterPlot.Unlock()
		return buffer, nil
	}
	return nil, errors.New("not enough data in storage")
}

func (sp *ScatterPlotWidget) update() {
	sp.ScatterPlot.Data = make([][]float64, len(sp.storages))
	for i, storage := range sp.storages {
		var ringBufferData = storage.Data.Tail(storage.Data.length)
		var data = make([]float64, 0, len(ringBufferData))
		for _, s := range ringBufferData {
			data = append(data, s.(float64))
		}
		if len(data) > 0 {
			sp.ScatterPlot.Data[i] = data
		}
	}
}
