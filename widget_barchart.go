package tchart

import (
	"fmt"
	"math"

	"github.com/s-westphal/termui/v3"
	ui "github.com/s-westphal/termui/v3"
	"github.com/s-westphal/termui/v3/widgets"
)

type barChartWidget struct {
	barChart *widgets.BarChart
	storage  *Storage
	numBins  int
}

func newBarChartWidget(storage *Storage, title string, numBins int) *barChartWidget {
	bc := widgets.NewBarChart()
	bc.Title = title
	bc.BarColors = []ui.Color{ui.ColorWhite}
	bc.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorGreen), ui.NewStyle(ui.ColorYellow)}
	bc.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}
	bc.BarWidth = 3
	bc.Border = false
	return &barChartWidget{
		barChart: bc,
		storage:  storage,
		numBins:  numBins,
	}
}

func (bc *barChartWidget) Buffer() (*termui.Buffer, error) {
	buffer := termui.NewBuffer(bc.barChart.GetRect())
	bc.barChart.Lock()
	bc.barChart.Draw(buffer)
	bc.barChart.Unlock()
	return buffer, nil
}

func (bc *barChartWidget) setDimension(x, y, w, h int) {
	bc.barChart.SetRect(x, y, x+w, y+h)
}

func (bc *barChartWidget) update() {
	if bc.storage.Data.Len() == 0 {
		// Histogram does not work for empty data
		return
	}
	stepsize := (bc.storage.Max - bc.storage.Min) / float64(bc.numBins)
	data := make([]float64, bc.numBins)
	labels := make([]string, bc.numBins)
	var prevValue float64 = 0
	for i := 0; i < bc.numBins; i++ {
		cummulativeValue := bc.storage.Histogram.CDF(bc.storage.Min + stepsize*float64(i+1))
		data[i] = math.Round((cummulativeValue - prevValue) * 100)
		prevValue = cummulativeValue
		labels[i] = fmt.Sprintf("%.1f", bc.storage.Min+stepsize*(float64(i)+0.5))
	}
	bc.barChart.Labels = labels
	bc.barChart.Data = data
}
