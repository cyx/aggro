package main

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"time"
)

const timeFormat = "2006-01-02T15:04:05+00:00"

type Output struct {
	Timestamp time.Time
	Host      string
	Service   int64
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{sc: bufio.NewScanner(r)}
}

type Scanner struct {
	sc     *bufio.Scanner
	err    error
	output Output
}

func (s *Scanner) Output() Output {
	return s.output
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

	s.output, s.err = extract(s.sc.Text())
	return s.err == nil
}

func extract(txt string) (Output, error) {
	s := bufio.NewScanner(strings.NewReader(txt))
	s.Split(bufio.ScanWords)

	o := Output{}

	s.Scan()
	t, err := time.Parse(timeFormat, s.Text())
	if err != nil {
		return o, err
	}
	o.Timestamp = t

	// Discard the next token
	s.Scan()

	for s.Scan() {
		t := s.Text()
		if strings.HasPrefix(t, "host=") {
			o.Host = strings.Trim(t[5:], `"`)
		} else if strings.HasPrefix(t, "service=") {
			svc, err := strconv.ParseInt(strings.TrimSuffix(t[8:], "ms"), 10, 64)
			if err != nil {
				return o, err
			}
			o.Service = svc
		}
	}
	return o, nil
}
