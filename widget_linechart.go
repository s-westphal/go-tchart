package tchart

import (
	"errors"

	"github.com/s-westphal/termui/v3"
	ui "github.com/s-westphal/termui/v3"
	"github.com/s-westphal/termui/v3/widgets"
)

// LineChartWidget line chart widget
type LineChartWidget struct {
	LineChart *widgets.Plot
	storages  []*Storage
	width     int
	height    int
}

// NewLineChartWidget create line chart widget
func NewLineChartWidget(title string, storages []*Storage, colors []ui.Color) *LineChartWidget {
	lineChart := widgets.NewPlot()
	lineChart.Title = title
	lineChart.LineColors = colors

	lc := &LineChartWidget{
		LineChart: lineChart,
		storages:  storages,
		width:     0,
		height:    0,
	}
	return lc
}

func (lc *LineChartWidget) setDimension(x, y, w, h int) {
	lc.width = w
	lc.height = h
	lc.LineChart.SetRect(x, y, x+w, y+h)
}

// Buffer buffer for rendering
func (lc *LineChartWidget) Buffer() (*termui.Buffer, error) {
	if lc.storages[0].Data.length > 1 {
		buffer := termui.NewBuffer(lc.LineChart.GetRect())
		lc.LineChart.Lock()
		lc.LineChart.Draw(buffer)
		lc.LineChart.Unlock()
		return buffer, nil
	}
	return nil, errors.New("not enough data in storage")
}

func (lc *LineChartWidget) update() {
	lc.LineChart.Data = make([][]float64, len(lc.storages))
	for i, storage := range lc.storages {
		var ringBufferData = storage.Data.Tail(lc.width)
		var data = make([]float64, 0, len(ringBufferData))
		for _, s := range ringBufferData {
			data = append(data, s.(float64))
		}
		if len(data) > 0 {
			lc.LineChart.Data[i] = data
		}
	}
	var ringBufferLabels = lc.storages[0].DataLabels.Tail(lc.width)
	var labels = make([]string, 0, len(ringBufferLabels))
	for _, s := range ringBufferLabels {
		labels = append(labels, s.(string))
	}
	lc.LineChart.DataLabels = labels

}
