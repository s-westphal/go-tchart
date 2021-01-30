package tchart

import (
	"errors"

	"github.com/nsf/termbox-go"
	ui "github.com/s-westphal/termui/v3"
)

func GetDefaultChartColors() []ui.Color {
	return []ui.Color{ui.ColorRed, ui.ColorGreen, ui.ColorYellow, ui.ColorBlue, ui.ColorCyan}
}

type App struct {
	*vContainer
	panels  []panel
	widgets []Widget
}

func NewApp() *App {
	return &App{
		vContainer: newVContainer(),
		panels:     []panel{},
		widgets:    []Widget{},
	}
}

func (app *App) AddPanel(panelType string, storages []*Storage) error {
	widget, err := CreateWidget(panelType, "", storages)
	if err != nil {
		panic(err)
	}
	var panel = newPanel(storages, widget)

	app.vContainer.putContainers(panel.container)
	app.panels = append(app.panels, *panel)
	return nil
}

func (app *App) AddInstructions() {
	instructionsWidget := newInstructionsWidget("")
	instructionsContainer := newHContainer(instructionsWidget)
	instructionsContainer.setHeight(3)
	app.vContainer.putContainers(instructionsContainer)
}

func (app *App) AddWidgetRow(widgets []Widget, height int) {
	rowContainer := newHContainer()
	for _, w := range widgets {
		chartContainer := newVContainer(w)
		rowContainer.putContainers(chartContainer)
	}
	app.widgets = append(app.widgets, widgets...)
	if height != 0 {
		rowContainer.setHeight(height)
	}
	app.vContainer.putContainers(rowContainer)
}

// Update update panels
func (app *App) Update() {

	for _, panel := range app.panels {
		panel.update()

	}
	for _, widget := range app.widgets {
		widget.update()
	}
	w, h := termbox.Size()

	app.vContainer.render(0, 0, w, h)
	termbox.Flush()
}

// Render render tchart
func (app *App) Render() {
	w, h := termbox.Size()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	app.vContainer.render(0, 0, w, h)
	termbox.Flush()
}

func CreateWidget(widgetType string, title string, storages []*Storage) (Widget, error) {
	var widget Widget
	switch widgetType {
	case "L":
		if title == "" {
			title = "Line Chart"
		}
		widget = NewLineChartWidget(title, storages, GetDefaultChartColors())
	case "S":
		if len(storages) != 2 {
			return nil, errors.New("scatter plot needs 2 columns")
		}
		if title == "" {
			title = "Scatter Plot"
		}
		widget = NewScatterPlotWidget(title, storages)
	case "P":
		if title == "" {
			title = "Pie Chart"
		}
		widget = NewPieChartWidget(title, storages, GetDefaultChartColors())
	case "G":
		if len(storages) != 1 {
			return nil, errors.New("gauge needs 1 column")
		}

		widget = NewGaugeWidget(title, storages[0], GetDefaultChartColors())
	}
	return widget, nil

}
