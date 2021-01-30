package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/s-westphal/go-tchart"
	cli "gopkg.in/alecthomas/kingpin.v2"

	"github.com/andrew-d/go-termutil"
	"github.com/nsf/termbox-go"
	ui "github.com/s-westphal/termui/v3"
)

const (
	histogramBinCount        int    = 50
	defaultBufferSizeScatter int    = 100
	bufferSizeLine           int    = 500
	defaultPanelType         string = "L"
	readFewLines             int    = 1
	readDefaultLines         int    = 5
	readManyLines            int    = 10
)

var (
	delimiter = cli.Flag(
		"delimiter",
		"Delimiter of input",
	).Short('d').Default("\t").String()

	panelTypes = cli.Flag(
		"panels",
		"Panel type of each column ("+
			"'L': line chart, "+
			"'S': scatter plot, "+
			"'P': pie chart, "+
			"'.': add to previous panel, "+
			"'x': skip column"+
			")",
	).Short('p').String()

	dataLabelMode = cli.Flag(
		"data-label",
		"Data labels used by line chart (None, first, time)",
	).Short('l').Default("None").String()

	dataSpeed = cli.Flag(
		"speed",
		"Speed at which to scroll though the data (slow, medium, fast)",
	).Short('s').Default("medium").String()

	bufferSizeScatter = cli.Flag(
		"num-samples",
		"Number of samples displayed in scatter plot",
	).Short('n').Default(fmt.Sprintf("%v", defaultBufferSizeScatter)).Int()
)

func storeRecords(storages []*tchart.Storage, dataLabelMode string, records [][]string) {
	for _, record := range records {
		label := ""
		switch dataLabelMode {
		case "first":
			label = record[0]
			record = record[1:]
		case "time":
			now := time.Now()
			label = fmt.Sprintf("%02d:%02d:%02d.%03d",
				now.Hour(), now.Minute(), now.Second(), now.Nanosecond()/1e6)
		}

		for i := 0; i < len(record); i++ {
			if x, err := strconv.ParseFloat(record[i], 64); err == nil {
				storages[i].Add(x, label)
			}
		}
	}
}

func main() {

	cli.Parse()

	var reader *csv.Reader
	if !termutil.Isatty(os.Stdin.Fd()) {
		reader = csv.NewReader(os.Stdin)
	} else {
		return
	}
	reader.Comma = []rune(*delimiter)[0]

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	app := tchart.NewApp()

	// read header
	titles, err := reader.Read()
	if err != nil {
		panic(err)
	}

	if *dataLabelMode == "first" {
		titles = titles[1:]
	}

	pTypes := strings.Split(*panelTypes, "")
	currentPanelType := defaultPanelType
	columnStorages := make([]*tchart.Storage, len(titles))
	panelStorages := make([]*tchart.Storage, 0, len(titles))
	for i, title := range titles {
		panelType := defaultPanelType
		if len(pTypes) > i {
			panelType = pTypes[i]
		}
		switch panelType {
		case "L", "S", "P":
			if len(panelStorages) > 0 {
				err := app.AddPanel(currentPanelType, panelStorages)
				if err != nil {
					panic(err)
				}
			}
			currentPanelType = panelType
			panelStorages = nil
		}

		var bufferSize int
		switch currentPanelType {
		case "L":
			bufferSize = bufferSizeLine
		case "S":
			bufferSize = *bufferSizeScatter
		case "P", "G":
			bufferSize = 1
		}
		columnStorages[i] = tchart.NewStorage(title, bufferSize, histogramBinCount)
		if panelType != "x" {
			panelStorages = append(panelStorages, columnStorages[i])
		}
	}
	if len(panelStorages) > 0 {
		err := app.AddPanel(currentPanelType, panelStorages)
		if err != nil {
			panic(err)
		}
	}
	app.AddInstructions()
	app.Render()

	pause := make(chan bool)
	evt := make(chan termbox.Event)
	dataChan := make(chan [][]string)
	var readLines int
	switch *dataSpeed {
	case "slow":
		readLines = readFewLines
	case "medium":
		readLines = readDefaultLines
	case "fast":
		readLines = readManyLines
	}
	go func() {
		for {
			evt <- termbox.PollEvent()
		}
	}()
	go func() {
		p := false
		tick := time.Tick(time.Duration(100) * time.Millisecond)
		for {
			select {
			case p = <-pause:

			case <-tick:
				if !p {
					lines := make([][]string, 0, readLines)
					for i := 0; i < readLines; i++ {
						line, err := reader.Read()
						if err != nil {
							if err == io.EOF {
								return
							}
							panic(err)
						}
						lines = append(lines, line)
					}
					dataChan <- lines

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
			if e.Type == termbox.EventKey && e.Ch == 'p' {
				pause <- true
			}
			if e.Type == termbox.EventKey && e.Ch == 'c' {
				pause <- false
			}
			if e.Type == termbox.EventResize {
				app.Render()
			}
		case records := <-dataChan:
			storeRecords(columnStorages, *dataLabelMode, records)
			app.Update()
		}
	}
}
