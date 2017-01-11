package main

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestAggro(t *testing.T) {
	input := `2017-01-09T22:44:10.235447+00:00 heroku[router]: at=info method=GET path="/v2/all" host="graphservice.com" request_id=5b5dcecc-0298-2f20-6252-f087ea02d89a fwd="52.70.203.72" dyno=web.2 connect=0ms service=8ms status=200 bytes=27
2017-01-09T22:44:24.780606+00:00 heroku[router]: at=info method=GET path="/events/alerts" host="eventthing.com" request_id=9d9f1740-0c2f-d1e1-9131-69d5deb4b1ac fwd="52.70.203.72" dyno=web.1 connect=0ms service=894ms status=200 bytes=27
2017-01-09T22:44:48.998716+00:00 heroku[router]: at=info method=GET path="/events/restarts" host="eventthing.com" request_id=a35f8a25-dc83-ddc0-35c4-90e419fbcb96 fwd="52.70.203.72" dyno=web.4 connect=0ms service=4ms status=200 bytes=27
2017-01-09T22:44:49.126254+00:00 heroku[router]: at=info method=GET path="/metrics/load" host="graphservice.com" request_id=df583531-b56a-8b6a-3aa3-f3d348fd92f2 fwd="52.70.203.72" dyno=web.1 connect=0ms service=4ms status=200 bytes=27
2017-01-09T22:44:55.287440+00:00 heroku[router]: at=info method=GET path="/comics/search?q=wolverine" host="comicstore.com" request_id=32264680-0f73-7f61-3174-88f4dd42c9b3 fwd="52.70.203.72" dyno=web.3 connect=0ms service=9ms status=200 bytes=27`

	scanner := NewScanner(strings.NewReader(input))

	buf := &bytes.Buffer{}
	aggro := NewAggro(time.Minute, buf)
	for scanner.Scan() {
		aggro.Observe(scanner.Event())
	}
	aggro.Close()

	expected :=
		`2017-01-09T22:44:00,comicstore.com,1,9,9,9
2017-01-09T22:44:00,eventthing.com,2,898,4,894
2017-01-09T22:44:00,graphservice.com,2,12,4,8
`

	if buf.String() != expected {
		t.Fatalf("Expected: \n\n%s\n\nActual: %s", expected, buf.String())
	}
}
