package tchart

import (
	"math"

	"github.com/VividCortex/gohistogram"
)

type Storage struct {
	title      string
	DataLabels *RingBuffer
	Data       *RingBuffer
	Histogram  *gohistogram.NumericHistogram
	NumBins    int
	Min        float64
	Max        float64
	Sum        float64
	Count      int
}

func NewStorage(title string, bufsize int, histogramBinCount int) *Storage {
	storage := &Storage{
		title:      title,
		DataLabels: NewRingBuffer(bufsize),
		Data:       NewRingBuffer(bufsize),
		Histogram:  gohistogram.NewHistogram(histogramBinCount),
		NumBins:    histogramBinCount,
		Min:        math.Inf(1),
		Max:        math.Inf(-1),
		Sum:        0,
		Count:      0,
	}
	return storage
}

func (storage *Storage) SetBuffersize(bufsize int) {
	storage.DataLabels = NewRingBuffer(bufsize)
	storage.Data = NewRingBuffer(bufsize)

}

func (storage *Storage) Add(x float64, dataLabel string) {
	storage.Data.Add(x)
	storage.Histogram.Add(x)

	if x < storage.Min {
		storage.Min = x
	}
	if x > storage.Max {
		storage.Max = x
	}
	storage.Sum += x
	storage.Count++

	if dataLabel != "" {
		storage.DataLabels.Add(dataLabel)
	}
}
