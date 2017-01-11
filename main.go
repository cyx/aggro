package main

import (
	"os"
	"time"
)

func main() {
	scanner := NewScanner(os.Stdin)
	aggro := NewAggro(time.Minute, os.Stdout)
	for scanner.Scan() {
		aggro.Observe(scanner.Event())
	}

	// Potentially flush anything currently stored in memory.
	aggro.Close()
}
