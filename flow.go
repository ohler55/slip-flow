// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main

// TBD flow-flavor

type flow struct {
	name string
	// self *flavors.Instance
	// tasks     map[string]*task
	errorTask *task

	// received  atomic.Uint64
	// errors    atomic.Uint64
	// processed atomic.Uint64
	// duration  atomic.Uint64 // sum of processing times from entry to completed from box tracks
}
