package main

import (
	"fmt"
	"io"
	"time"
)

const timeOutputFormat = "2006-01-02T15:04:05"

func NewAggro(res time.Duration, writer io.Writer) *Aggro {
	return &Aggro{resolution: res, writer: writer, sl: NewSortedList(), points: Points{}}
}

type Aggro struct {
	resolution        time.Duration
	currentWindow     time.Time
	sl                *SortedList
	points            Points
	writer            io.Writer
	lastFlushedWindow time.Time
}

func (a *Aggro) Observe(e Event) {
	ts := e.Timestamp.Truncate(a.resolution)

	if ts.After(a.currentWindow) {
		a.flush()
	}

	a.currentWindow = ts
	a.sl.Add(e.Host)
	a.points.Observe(e)
}

func (a *Aggro) Close() error {
	a.flush()
	return nil
}

func (a *Aggro) flush() {
	// Don't double flush things.
	if a.currentWindow == a.lastFlushedWindow {
		return
	}

	a.lastFlushedWindow = a.currentWindow
	ts := a.currentWindow.Format(timeOutputFormat)
	a.sl.Each(func(name string) {
		p, ok := a.points[name]
		if !ok {
			p = &Point{}
		}
		fmt.Fprintf(a.writer, "%s,%s,%d,%d,%d,%d\n", ts, name, p.count, p.sum, p.min, p.max)
	})

	// Reset points
	a.points = Points{}
}
