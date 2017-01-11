package main

import (
	"strings"
	"testing"
)

func TestScanner(t *testing.T) {
	input := `2017-01-09T22:44:10.235447+00:00 heroku[router]: at=info method=GET path="/v2/all" host="graphservice.com" request_id=5b5dcecc-0298-2f20-6252-f087ea02d89a fwd="52.70.203.72" dyno=web.2 connect=0ms service=8ms status=200 bytes=27
2017-01-09T22:44:24.780606+00:00 heroku[router]: at=info method=GET path="/events/alerts" host="eventthing.com" request_id=9d9f1740-0c2f-d1e1-9131-69d5deb4b1ac fwd="52.70.203.72" dyno=web.1 connect=0ms service=894ms status=200 bytes=27
2017-01-09T22:44:48.998716+00:00 heroku[router]: at=info method=GET path="/events/restarts" host="eventthing.com" request_id=a35f8a25-dc83-ddc0-35c4-90e419fbcb96 fwd="52.70.203.72" dyno=web.4 connect=0ms service=4ms status=200 bytes=27
2017-01-09T22:44:49.126254+00:00 heroku[router]: at=info method=GET path="/metrics/load" host="graphservice.com" request_id=df583531-b56a-8b6a-3aa3-f3d348fd92f2 fwd="52.70.203.72" dyno=web.1 connect=0ms service=4ms status=200 bytes=27
2017-01-09T22:44:55.287440+00:00 heroku[router]: at=info method=GET path="/comics/search?q=wolverine" host="comicstore.com" request_id=32264680-0f73-7f61-3174-88f4dd42c9b3 fwd="52.70.203.72" dyno=web.3 connect=0ms service=9ms status=200 bytes=27`

	scanner := NewScanner(strings.NewReader(input))

	tally := make(map[string]int64)

	for scanner.Scan() {
		o := scanner.Output()
		tally[o.Host] += o.Service
	}

	if scanner.Err() != nil {
		t.Fatal(scanner.Err())
	}

	if tally["graphservice.com"] != 12 {
		t.Fatalf("%d != %d", tally["graphservice.com"], 12)
	}

	if tally["eventthing.com"] != 898 {
		t.Fatalf("%d != %d", tally["eventthing.com"], 898)
	}

	if tally["comicstore.com"] != 9 {
		t.Fatalf("%d != %d", tally["comicstore.com"], 9)
	}

}
