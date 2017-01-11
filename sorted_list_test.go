package main

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/pborman/uuid"
)

func TestSortedList(t *testing.T) {
	sl := NewSortedList()
	sl.Add("zoo.com")
	sl.Add("abc.com")
	sl.Add("def.com")
	sl.Add("zoo.com")

	var res []string

	sl.Each(func(n string) {
		res = append(res, n)
	})

	if !reflect.DeepEqual([]string{"abc.com", "def.com", "zoo.com"}, res) {
		t.Fatalf("Expected %+v to be abc.com def.com zoo.com", res)
	}
}

func BenchmarkSortedListAdd(b *testing.B) {
	N := 1000
	names := make([]string, N)
	for i := 0; i < N; i++ {
		names[i] = uuid.New() + ".com"
	}

	sl := NewSortedList()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sl.Add(names[rand.Intn(N)])
	}
	b.StopTimer()
}
