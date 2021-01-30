package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/s-westphal/go-tchart"

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
	delimiter = flag.String(
		"delimiter",
		"\t",
		"Delimiter of input",
	)

	panelTypes = flag.String(
		"panels",
		"",
		"Panel type of each column ("+
			"'L': line chart, "+
			"'S': scatter plot, "+
			"'P': pie chart, "+
			"'.': add to previous panel, "+
			"'x': skip column"+
			")",
	)

	dataLabelMode = flag.String(
		"data-label",
		"None",
		"Data labels used by line chart (None, first, time)",
	)

	dataSpeed = flag.String(
		"speed",
		"medium",
		"Speed at which to scroll though the data (slow, medium, fast)",
	)

	bufferSizeScatter = flag.Int(
		"num-samples",
		defaultBufferSizeScatter,
		"Number of samples displayed in scatter plot",
	)
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

	flag.Parse()

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
