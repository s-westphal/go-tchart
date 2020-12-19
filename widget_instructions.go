package tchart

import (
	"fmt"

	"github.com/s-westphal/termui/v3"
	"github.com/s-westphal/termui/v3/widgets"
)

type instructionsWidget struct {
	list *widgets.List
}

func newInstructionsWidget(title string) *instructionsWidget {
	list := widgets.NewList()
	list.Title = title
	// list.Border = false

	instructionsWidget := &instructionsWidget{
		list: list,
	}
	instructionsWidget.update()
	return instructionsWidget
}

// Buffer buffer for rendering
func (sw *instructionsWidget) Buffer() (*termui.Buffer, error) {
	buffer := termui.NewBuffer(sw.list.GetRect())
	sw.list.Lock()
	sw.list.Draw(buffer)
	sw.list.Unlock()
	return buffer, nil
}

func (sw *instructionsWidget) setDimension(x, y, w, h int) {
	sw.list.SetRect(x, y, x+w, y+h)
}

func (sw *instructionsWidget) update() {
	sw.list.Rows = []string{
		fmt.Sprint("q: quit   p: pause   c: continue"),
	}
}
