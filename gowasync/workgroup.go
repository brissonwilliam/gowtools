package gowasync

import (
	"fmt"
	"sync/atomic"
	"time"
)

// WorkGroup is to run a given 'global' (to the work group) process on data
// in parallel on an X amount of workers
type WorkGroup interface {
	AwaitExecute() ErrorSyncGroup
	SetVerbose(v bool)
}

type defaultWorkGroup[T any] struct {
	workers   uint32
	workUnits []workUnit[T]
	verbose   bool
}

type workUnit[T any] struct {
	data    []T
	process func([]T) error
}

// NewWorkGroup returns a work group that can be executed to process data asynchronously

func NewWorkGroup[T any](workers uint32, dataChunked [][]T, process func([]T) error) WorkGroup {
	units := prepareWorkUnits[T](dataChunked, process)
	return &defaultWorkGroup[T]{
		workers:   workers,
		workUnits: units,
		verbose:   false,
	}
}

func prepareWorkUnits[T any](dataChunked [][]T, p func([]T) error) []workUnit[T] {
	wu := make([]workUnit[T], len(dataChunked))
	for i := range dataChunked {
		wu[i] = workUnit[T]{
			data:    dataChunked[i],
			process: p,
		}
	}
	return wu
}

func (t *defaultWorkGroup[T]) SetVerbose(v bool) {
	t.verbose = v
}

func (t *defaultWorkGroup[T]) AwaitExecute() ErrorSyncGroup {
	var count int64
	start := time.Now()

	workChan := make(chan workUnit[T], 1000)
	done := make(chan bool, t.workers) // channel to signal completion
	errs := ErrorSyncGroup{}

	// Spawn worker goroutines
	for i := 0; uint32(i) < t.workers; i++ {
		go func(i int) {
			// each worker waits for a new chunk then processes it
			for wu := range workChan {
				if err := wu.process(wu.data); err != nil {
					errs.Add(err)
				}

				if t.verbose {
					cur := atomic.AddInt64(&count, int64(len(wu.data)))
					rps := float64(cur) / time.Since(start).Seconds()
					fmt.Printf("worker %d consumed %d objects. Group speed at speed %.2f/s \n", i, cur, rps)
				}
			}

			// Signal completion to the main function when the worker goroutine finishes
			done <- true
		}(i)
	}

	// Dispatch each work unit to workers
	for _, unitOfWork := range t.workUnits {
		workChan <- unitOfWork
	}

	close(workChan) // close() causes for range to exit once all has been processed in the channel

	// Wait for all workers' signal of finished processing
	for i := 0; uint32(i) < t.workers; i++ {
		<-done
	}

	if t.verbose {
		fmt.Printf("Work group done processing %d messages in %s", count, time.Since(start).String())
	}

	return errs
}
