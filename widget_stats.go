package tchart

import (
	"fmt"
	"math"

	"github.com/s-westphal/termui/v3"
	"github.com/s-westphal/termui/v3/widgets"
)

type statsWidget struct {
	list    *widgets.List
	storage *Storage
}

func newStatsWidget(storage *Storage, title string) *statsWidget {
	list := widgets.NewList()
	list.Title = title
	list.Border = false
	return &statsWidget{
		list:    list,
		storage: storage,
	}
}

// Buffer buffer for rendering
func (sw *statsWidget) Buffer() (*termui.Buffer, error) {
	buffer := termui.NewBuffer(sw.list.GetRect())
	sw.list.Lock()
	sw.list.Draw(buffer)
	sw.list.Unlock()
	return buffer, nil
}

func (sw *statsWidget) setDimension(x, y, w, h int) {
	sw.list.SetRect(x, y, x+w, y+h)
}

func (sw *statsWidget) update() {
	if sw.storage.Data.Len() == 0 {
		// Histogram does not work for empty data
		return
	}
	var lastValue float64
	sliceLastValue := sw.storage.Data.Tail(1)
	if len(sliceLastValue) == 1 {
		lastValue = sliceLastValue[0].(float64)
	}
	sw.list.Rows = []string{
		fmt.Sprintf("Count : %.0d", sw.storage.Count),
		fmt.Sprintf("Min   : %.2f", sw.storage.Min),
		fmt.Sprintf("Max   : %.2f", sw.storage.Max),
		fmt.Sprintf("Mean  : %.2f", sw.storage.Histogram.Mean()),
		fmt.Sprintf("Std   : %.2f", math.Sqrt(sw.storage.Histogram.Variance())),
		fmt.Sprintf("Last  : %.2f", lastValue),
	}
}
