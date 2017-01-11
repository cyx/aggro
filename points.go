package main

import "math"

type Points map[string]*Point

func (p Points) Observe(e Event) {
	if p[e.Host] == nil {
		p[e.Host] = &Point{min: math.MaxInt64}
	}
	p[e.Host].Observe(e)
}

type Point struct {
	sum   int64
	min   int64
	max   int64
	count int64
}

func (p *Point) Observe(e Event) {
	p.count++
	p.sum += e.Service
	if e.Service < p.min {
		p.min = e.Service
	}
	if e.Service > p.max {
		p.max = e.Service
	}
}
