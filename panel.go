package tchart

import ui "github.com/s-westphal/termui/v3"

// panel parameters
const (
	dataListWidth = 31
)

type panel struct {
	container *hContainer
	stats     []*statsWidget
	barCharts []*barChartWidget
	widget    widget
}

func getChartColors() []ui.Color {
	return []ui.Color{ui.ColorRed, ui.ColorGreen, ui.ColorYellow, ui.ColorBlue, ui.ColorCyan}
}

func newPanel(storages []*Storage, widget widget) *panel {
	container := &hContainer{}
	stats := make([]*statsWidget, len(storages))
	barCharts := make([]*barChartWidget, len(storages))
	colors := getChartColors()
	for i, storage := range storages {
		frame := newFrameWidget(storage.title, colors[i])
		dataStatsContainer := newVContainer(frame)
		dataStatsContainer.setWidth(dataListWidth)
		statsWidget := newStatsWidget(storage, "")
		statsContainer := newHContainer(statsWidget)
		stats[i] = statsWidget
		barChartWidget := newBarChartWidget(storage, "", 7)
		barChartContainer := newHContainer(barChartWidget)
		barCharts[i] = barChartWidget
		dataStatsContainer.putContainers(statsContainer, barChartContainer)
		container.putContainers(dataStatsContainer)
	}

	widgetContainer := newVContainer(widget)
	container.putContainers(widgetContainer)

	return &panel{
		container: container,
		stats:     stats,
		barCharts: barCharts,
		widget:    widget,
	}
}

func (p *panel) update() {
	for _, s := range p.stats {
		s.update()
	}
	for _, s := range p.barCharts {
		s.update()
	}
	p.widget.update()
}
