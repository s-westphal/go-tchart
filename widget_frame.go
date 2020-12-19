package tchart

import (
	"github.com/s-westphal/termui/v3"
	ui "github.com/s-westphal/termui/v3"
)

type frameWidget struct {
	frame *ui.Block
}

func newFrameWidget(title string, color ui.Color) *frameWidget {
	block := ui.NewBlock()
	block.Title = title
	block.TitleStyle = ui.NewStyle(color)
	return &frameWidget{
		frame: block,
	}
}

// Buffer buffer for rendering
func (fw *frameWidget) Buffer() (*termui.Buffer, error) {
	buffer := termui.NewBuffer(fw.frame.GetRect())
	fw.frame.Lock()
	fw.frame.Draw(buffer)
	fw.frame.Unlock()
	return buffer, nil
}

func (fw *frameWidget) setDimension(x, y, w, h int) {
	fw.frame.SetRect(x, y, x+w, y+h)
}

func (fw *frameWidget) update() {
}
