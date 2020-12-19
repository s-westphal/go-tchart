package tchart

import (
	"errors"

	"github.com/nsf/termbox-go"
)

// App app
type App struct {
	*vContainer
	panels []panel
}

// NewApp new app
func NewApp() *App {
	return &App{
		vContainer: newVContainer(),
		panels:     []panel{},
	}
}

// AddPanel add panel to app
func (app *App) AddPanel(panelType string, storages []*Storage) error {
	var widget widget
	switch panelType {
	case "L":
		widget = NewLineChartWidget("Line Chart", storages, getChartColors())
	case "S":
		if len(storages) != 2 {
			return errors.New("scatter plot needs 2 columns")
		}
		widget = NewScatterPlotWidget("Scatter Plot", storages)
	case "P":
		widget = NewPieChartWidget("Pie Chart", storages, getChartColors())
	}
	var panel = newPanel(storages, widget)

	app.vContainer.putContainers(panel.container)
	app.panels = append(app.panels, *panel)
	return nil
}

// AddInstructions add instractions
func (app *App) AddInstructions() {
	instructionsWidget := newInstructionsWidget("")
	instructionsContainer := newHContainer(instructionsWidget)
	instructionsContainer.setHeight(3)
	app.vContainer.putContainers(instructionsContainer)
}

// Update update panels
func (app *App) Update() {

	for _, panel := range app.panels {
		panel.update()

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
