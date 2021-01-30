package tchart

// panel parameters
const (
	dataListWidth = 31
)

type panel struct {
	container *hContainer
	stats     []*statsWidget
	barCharts []*barChartWidget
	widget    Widget
}

func newPanel(storages []*Storage, widget Widget) *panel {
	container := &hContainer{}
	stats := make([]*statsWidget, len(storages))
	barCharts := make([]*barChartWidget, len(storages))
	colors := GetDefaultChartColors()
	for i, storage := range storages {
		frame := newFrameWidget(storage.title, colors[i])
		dataStatsContainer := newVContainer(frame)
		dataStatsContainer.setWidth(dataListWidth)
		statsWidget := NewStatsWidget("", storage)
		statsContainer := newHContainer(statsWidget)
		stats[i] = statsWidget
		barChartWidget := NewBarChartWidget("", storage, 7)
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
