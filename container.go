package tchart

import (
	tb "github.com/nsf/termbox-go"
)

type container interface {
	getWidth() int
	getHeight() int
	setWidth(w int)
	setHeight(h int)
	render(x, y, w, h int)
}

type containerBase struct {
	width  int
	height int
}

func (bb *containerBase) getWidth() int {
	return bb.width
}

func (bb *containerBase) getHeight() int {
	return bb.height
}

func (bb *containerBase) setWidth(w int) {
	bb.width = w
}

func (bb *containerBase) setHeight(h int) {
	bb.height = h
}

type layoutContainer struct {
	containerBase
	innerContainers []container
	widgets         []Widget
}

func (lb *layoutContainer) putContainers(containers ...container) {
	lb.innerContainers = append(lb.innerContainers, containers...)
}

func (lb *layoutContainer) renderWidget(x, y, w, h int) (xNew, yNew, wNew, hNew int) {

	for _, widget := range lb.widgets {
		widget.setDimension(x, y, w, h)
		buf, err := widget.Buffer()
		if err != nil {
			return
		}
		for point, cell := range buf.CellMap {
			if point.In(buf.Rectangle) {
				tb.SetCell(
					point.X, point.Y,
					cell.Rune,
					tb.Attribute(cell.Style.Fg+1)|tb.Attribute(cell.Style.Modifier), tb.Attribute(cell.Style.Bg+1),
				)
			}
		}
		// render innerContainers inside border of this widget
		x += 1
		y += 1
		w -= 2
		h -= 2
	}

	return x, y, w, h

}

type vContainer struct {
	layoutContainer
}

func newVContainer(widgets ...Widget) *vContainer {
	vb := &vContainer{}
	vb.innerContainers = []container{}
	vb.widgets = widgets
	return vb
}

func (vb *vContainer) addWidgets(widgets ...Widget) {
	vb.widgets = widgets
}

func (vb *vContainer) render(x, y, w, h int) {
	x, y, w, h = vb.renderWidget(x, y, w, h)
	reservedHeight := 0
	numFlexContainers := 0
	for _, container := range vb.innerContainers {
		if container.getHeight() <= 0 {
			numFlexContainers++
		} else {
			reservedHeight += container.getHeight()
		}
	}
	containerFlexHeight := 0
	remainingHeight := h - reservedHeight
	if numFlexContainers > 0 && remainingHeight > 0 {
		containerFlexHeight = remainingHeight / numFlexContainers
	}

	for _, container := range vb.innerContainers {
		containerHeight := container.getHeight()
		if containerHeight <= 0 {
			containerHeight = containerFlexHeight
		}
		if containerHeight > 0 {
			container.render(x, y, w, containerHeight)
			y += containerHeight
		}
	}
}

type hContainer struct {
	layoutContainer
}

func newHContainer(widgets ...Widget) *hContainer {
	hb := &hContainer{}
	hb.innerContainers = []container{}
	hb.widgets = widgets
	return hb
}

func (hb *hContainer) addWidgets(widgets ...Widget) {
	hb.widgets = widgets
}

func (hb *hContainer) render(x, y, w, h int) {
	x, y, w, h = hb.renderWidget(x, y, w, h)
	reservedWidth := 0
	numFlexContainers := 0
	for _, container := range hb.innerContainers {
		if container.getWidth() <= 0 {
			numFlexContainers++
		} else {
			reservedWidth += container.getWidth()
		}
	}

	containerFlexWidth := 0
	remainingWidth := w - reservedWidth
	if numFlexContainers > 0 && remainingWidth > 0 {
		containerFlexWidth = remainingWidth / numFlexContainers
	}

	for _, container := range hb.innerContainers {
		containerWidth := container.getWidth()
		if containerWidth <= 0 {
			containerWidth = containerFlexWidth
		}
		if containerWidth > 0 {
			container.render(x, y, containerWidth, h)
			x += containerWidth
		}
	}
}
