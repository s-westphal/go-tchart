package tchart

import (
	"errors"
	"fmt"
	"math"

	"github.com/s-westphal/termui/v3"
	ui "github.com/s-westphal/termui/v3"
	"github.com/s-westphal/termui/v3/widgets"
)

// PieChartWidget line chart widget
type PieChartWidget struct {
	PieChart *widgets.PieChart
	storages []*Storage
	width    int
	height   int
}

// NewPieChartWidget create line chart widget
func NewPieChartWidget(title string, storages []*Storage, colors []ui.Color) *PieChartWidget {
	pieChart := widgets.NewPieChart()
	pieChart.Title = title
	pieChart.Colors = colors
	pieChart.AngleOffset = -.5 * math.Pi
	pieChart.LabelFormatter = func(i int, v float64) string {
		return fmt.Sprintf("%.02f", v)
	}

	pc := &PieChartWidget{
		PieChart: pieChart,
		storages: storages,
		width:    0,
		height:   0,
	}
	return pc
}

func (pc *PieChartWidget) setDimension(x, y, w, h int) {
	pc.width = w
	pc.height = h
	pc.PieChart.SetRect(x, y, x+w, y+h)
}

// Buffer buffer for rendering
func (pc *PieChartWidget) Buffer() (*termui.Buffer, error) {
	if pc.storages[0].Data.length > 0 {
		buffer := termui.NewBuffer(pc.PieChart.GetRect())

		pc.PieChart.Lock()
		pc.PieChart.Draw(buffer)
		pc.PieChart.Unlock()
		return buffer, nil
	}
	return nil, errors.New("not enough data in storage")
}

func (pc *PieChartWidget) update() {
	pc.PieChart.Data = make([]float64, len(pc.storages))
	for i, storage := range pc.storages {
		pc.PieChart.Data[i] = MaxFloat64(0, storage.Sum)
	}
}
