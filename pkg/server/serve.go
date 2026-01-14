package server

import (
	"fmt"
	"runtime"

	"github.com/ameyw07/go-kv-store/pkg/shard"
)

var NumWorkers = -1

func Serve() {
	numCPU := runtime.NumCPU()
	if numCPU <= 1 {
		fmt.Println("System has 1 or fewer CPU cores. Spawning 1 goroutine.")
		numCPU = 2 // Set to 2 so numWorkers is at least 1
	}

	// Calculate the number of workers (CPU cores - 1)
	NumWorkers := numCPU - 1
	// The Go runtime defaults GOMAXPROCS to the number of CPUs.
	// We can explicitly set it, though it's often not needed for this pattern.

	sc := shard.NewShardController(NumWorkers)

	// TODO: switch to waitGroup
	for i := 0; i < NumWorkers; i++ {
		iomx := NewIOMultiplexer()
		go sc.Shards[i].ExecuteCommand()
		go iomx.StartPollWorker(sc)

	}
}
