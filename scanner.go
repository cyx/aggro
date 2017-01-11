package main

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"time"
)

const timeFormat = "2006-01-02T15:04:05+00:00"

type Event struct {
	Timestamp time.Time
	Host      string
	Service   int64
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{sc: bufio.NewScanner(r)}
}

type Scanner struct {
	sc    *bufio.Scanner
	err   error
	event Event
}

func (s *Scanner) Event() Event {
	return s.event
}

func (s *Scanner) Err() error {
	return s.err
}

func (s *Scanner) Scan() bool {
	if s.err != nil {
		return false
	}

	if !s.sc.Scan() {
		s.err = s.sc.Err()
		return false
	}

	s.event, s.err = extract(s.sc.Text())
	return s.err == nil
}

func extract(txt string) (Event, error) {
	s := bufio.NewScanner(strings.NewReader(txt))
	s.Split(bufio.ScanWords)

	e := Event{}

	s.Scan()
	t, err := time.Parse(timeFormat, s.Text())
	if err != nil {
		return e, err
	}
	e.Timestamp = t.UTC()

	// Discard the next token
	s.Scan()

	for s.Scan() {
		t := s.Text()
		if strings.HasPrefix(t, "host=") {
			e.Host = strings.Trim(t[5:], `"`)
		} else if strings.HasPrefix(t, "service=") {
			svc, err := strconv.ParseInt(strings.TrimSuffix(t[8:], "ms"), 10, 64)
			if err != nil {
				return e, err
			}
			e.Service = svc
		}
	}
	return e, nil
}
