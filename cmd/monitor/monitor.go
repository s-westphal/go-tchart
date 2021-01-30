package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
	tchart "github.com/s-westphal/go-tchart"
	ui "github.com/s-westphal/termui/v3"
	cli "gopkg.in/alecthomas/kingpin.v2"
)

const (
	default_config    string = "monitor_config.json"
	histogramBinCount int    = 50
	bufferSizeScatter int    = 100
	bufferSizeLine    int    = 500
)

var (
	config = cli.Flag(
		"config",
		"path to config file",
	).Short('c').Default(default_config).String()
)

type chartEntry struct {
	Title     string `json:"title"`
	Frequency int    `json:"frequency"`
	Command   string `json:"command"`
	Delimiter string `json:"delimiter"`
	PlotSpec  string `json:"plotSpec"`
}

type rowEntry struct {
	Entries []chartEntry `json:"charts"`
	Height  int          `json:"height"`
}

func parseJSON(filename string) map[string][]rowEntry {

	jsonFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	var content = make(map[string][]rowEntry)
	if err := json.Unmarshal(bytes, &content); err != nil {
		panic(err)
	}
	return content
}

type CommandParser struct {
	command   string
	delimiter string
}

func NewCommandParser(command string, delimiter string) *CommandParser {
	return &CommandParser{
		command:   command,
		delimiter: delimiter,
	}
}

func (cp *CommandParser) runCommand() []float64 {
	commandResult := cp.execute()
	strSplit := strings.Split(commandResult, cp.delimiter)
	result := make([]float64, 0, len(strSplit))
	for _, res := range strSplit {
		if len(res) == 0 {
			// trailing delimiter or missing entry
			continue
		}
		if x, err := strconv.ParseFloat(res, 64); err == nil {
			result = append(result, x)
		} else {
			result = append(result, 0)
		}
	}
	return result
}

func (cp *CommandParser) execute() string {
	out, err := exec.Command("bash", "-c", cp.command).Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

type IntervalExecuter struct {
	interval      time.Duration
	commandParser *CommandParser
	storages      []*tchart.Storage
}

func NewIntervalExecuter(timeDelta int, commandParser *CommandParser, storages []*tchart.Storage) *IntervalExecuter {
	return &IntervalExecuter{
		interval:      time.Duration(timeDelta) * time.Second,
		commandParser: commandParser,
		storages:      storages,
	}
}

func (ie *IntervalExecuter) run() {
	go func() {
		tick := time.Tick(ie.interval)
		for {
			select {
			case <-tick:
				ie.addEvaluation()
			}
		}
	}()
}

func (ie *IntervalExecuter) addEvaluation() {

	now := time.Now()
	label := fmt.Sprintf("%02d:%02d:%02d.%03d",
		now.Hour(), now.Minute(), now.Second(), now.Nanosecond()/1e6)

	var values []float64 = ie.commandParser.runCommand()
	for i, value := range values {
		if i >= len(ie.storages) {
			err := fmt.Errorf("command returns more values than specified, skip")
			fmt.Println(err.Error())
			continue
		}
		// fmt.Println(value)
		ie.storages[i].Add(value, label)

	}
}

func main() {
	cli.Parse()

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	app := tchart.NewApp()

	conf := parseJSON(*config)
	for _, row := range conf["monitor"] {
		widgets := make([]tchart.Widget, 0, len(row.Entries))
		for _, chart := range row.Entries {
			chartType := string(chart.PlotSpec[0])
			numChartColumns := len(chart.PlotSpec)
			dataStorages := make([]*tchart.Storage, numChartColumns)

			var bufferSize int
			switch chartType {
			case "L":
				bufferSize = bufferSizeLine
			case "S":
				bufferSize = bufferSizeScatter
			case "P", "G":
				bufferSize = 1
			}

			for i := 0; i < numChartColumns; i++ {
				dataStorages[i] = tchart.NewStorage(chart.Title, bufferSize, histogramBinCount)
			}
			widget, err := tchart.CreateWidget(chartType, chart.Title, dataStorages)

			if err != nil {
				panic(err)
			}
			widgets = append(widgets, widget)

			intervalExecuter := NewIntervalExecuter(chart.Frequency, NewCommandParser(chart.Command, chart.Delimiter), dataStorages)
			intervalExecuter.run()
		}

		app.AddWidgetRow(widgets, row.Height)

	}

	app.Render()
	evt := make(chan termbox.Event)
	pause := make(chan bool)
	go func() {
		for {
			evt <- termbox.PollEvent()
		}
	}()
	go func() {
		tick := time.Tick(time.Duration(100) * time.Millisecond)
		var p bool = false
		for {
			select {
			case p = <-pause:

			case <-tick:
				if !p {
					app.Update()
				}
			}

		}
	}()

	for {
		select {
		case e := <-evt:

			if e.Type == termbox.EventKey && e.Ch == 'q' {
				return
			}

			if e.Type == termbox.EventResize {
				pause <- true
				app.Render()
				pause <- false
			}
		}
	}

}
