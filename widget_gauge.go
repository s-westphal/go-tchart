package tchart

import (
	"errors"

	"github.com/s-westphal/termui/v3"
	ui "github.com/s-westphal/termui/v3"
	"github.com/s-westphal/termui/v3/widgets"
)

type GaugeWidget struct {
	Gauge   *widgets.Gauge
	storage *Storage
	width   int
	height  int
}

func NewGaugeWidget(title string, storage *Storage, colors []ui.Color) *GaugeWidget {
	gauge := widgets.NewGauge()
	gauge.Label = title
	gauge.BarColor = colors[0]

	gw := &GaugeWidget{
		Gauge:   gauge,
		storage: storage,
		width:   0,
		height:  0,
	}
	return gw
}

func (gw *GaugeWidget) setDimension(x, y, w, h int) {
	gw.width = w
	gw.height = h
	gw.Gauge.SetRect(x, y, x+w, y+h)
}

func (gw *GaugeWidget) Buffer() (*termui.Buffer, error) {
	if gw.storage.Data.length > 0 {
		buffer := termui.NewBuffer(gw.Gauge.GetRect())

		gw.Gauge.Lock()
		gw.Gauge.Draw(buffer)
		gw.Gauge.Unlock()
		return buffer, nil
	}
	return nil, errors.New("not enough data in storage")
}

func (gw *GaugeWidget) update() {
	if gw.storage.Data.Len() > 0 {
		data := gw.storage.Data.Tail(1)[0].(float64)
		gw.Gauge.Percent = int(data * 100)
	}
}
