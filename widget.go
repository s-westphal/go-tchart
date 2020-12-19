package tchart

import "github.com/s-westphal/termui/v3"

type widget interface {
	setDimension(x, y, w, h int)
	Buffer() (*termui.Buffer, error)
	update()
}
