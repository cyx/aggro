package main

import (
	"container/list"
	"strings"
)

type SortedList struct {
	l *list.List

	// The set is an optimization for Add, drops the operation / second
	// for random operations from 10K ns to 95ns. The trade of memory
	// for speed here is a worthwhile endeavor.
	set map[string]bool
}

func NewSortedList() *SortedList {
	return &SortedList{l: list.New(), set: make(map[string]bool)}
}

func (s *SortedList) Add(name string) {
	if s.set[name] {
		return
	}

	s.set[name] = true

	for e := s.l.Front(); e != nil; e = e.Next() {
		n := e.Value.(string)
		switch cmp := strings.Compare(name, n); cmp {
		case 0:
			// It's already here, no need to double insert.
			return
		case -1:
			s.l.InsertBefore(name, e)
			return
		}
	}
	s.l.PushBack(name)
}

func (s *SortedList) Each(fn func(string)) {
	for e := s.l.Front(); e != nil; e = e.Next() {
		fn(e.Value.(string))
	}
}
