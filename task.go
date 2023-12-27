// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main

import (
	"sync/atomic"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

var (
	taskFlavor *flavors.Flavor
)

func init() {
	taskFlavor = flavors.DefFlavor("flow-task-flavor",
		map[string]slip.Object{},
		[]string{
			"can-log-flavor",
		},
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`Tasks are objects that implement the processing nodes withing a flow. Each task
has an actor that performs the processing of the task. Processing can be completed in the current thread or
queued and processed by workers in separate threads.`),
			},
			slip.List{
				slip.Symbol(":init-keywords"),
				slip.Symbol(":name"),
				slip.Symbol(":actor"),
				slip.Symbol(":workers"),
			},
		},
	)
	taskFlavor.DefMethod(":init", "", taskInitCaller{})
	// :name
	// :workers
	// :running
	// :metrics => assoc
	//   queue-length, average-time, received, errors, processed
	// :receive
	// :transition
	// :start
	//  reset/zero metrics
	// :shutdown (&optional wait)
	// TBD
}

type link struct {
	name string
	to   *task
}

type task struct {
	name      string
	self      *flavors.Instance
	links     map[string]*link
	actors    []slip.Instance
	caller    slip.Caller
	queue     chan *flavors.Instance // must be box instances
	done      chan struct{}
	workers   int
	depth     int
	received  atomic.Uint64
	errors    atomic.Uint64
	processed atomic.Uint64
	duration  atomic.Uint64 // sum of processing times (using box.track.events)
}

func (t *task) start(s *slip.Scope) {
	t.received.Store(0)
	t.errors.Store(0)
	t.processed.Store(0)
	t.duration.Store(0)
	for _, a := range t.actors {
		if a.HasMethod(":start") {
			a.Receive(s, ":start", slip.List{t.self}, 0)
		}
	}
	if 0 < t.workers {
		t.queue = make(chan *flavors.Instance, t.depth)

		// TBD start workers looping through actors or func
	}
	// TBD
}

func (t *task) actorLoop(s *slip.Scope, actor slip.Instance) {
	for {
		bi := t.queue
		if bi == nil {
			break
		}
		// TBD receive
		//
	}
}

// TBD deal with no queue

func (t *task) funcLoop(s *slip.Scope) {
	// TBD
}

func (t *task) shutdown(s *slip.Scope) {
	if t.queue != nil {
		for i := t.workers; 0 < i; i-- {
			t.queue <- nil
		}
		close(t.queue)
		for i := t.workers; 0 < i; i-- {
			<-t.done
		}
		t.queue = nil
	}
	for _, a := range t.actors {
		if a.HasMethod(":shutdown") {
			a.Receive(s, ":shutdown", slip.List{}, 0)
		}
	}
}

func (t *task) receive(s *slip.Scope, bi *flavors.Instance) {
	// TBD scan and freeze
	// if queue then place on queue else process here
}

func (t *task) act(s *slip.Scope, actor, bi *flavors.Instance) {
	defer func() {
		if rec := recover(); rec != nil {
			// TBD if error transition then follow it else try flow error handler
		}
	}()
	// TBD
	trans, bx := a.Receive(s, ":perform", slip.List{bi}, 0)

}

type taskInitCaller struct{}

func (caller taskInitCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		args = args[0].(slip.List)
	}
	tsk := task{self: obj, links: map[string]*link{}}
	for i := 0; i < len(args)-1; i += 2 {
		switch args[i] {
		case slip.Symbol(":name"):
			switch tv := args[i+1].(type) {
			case slip.String:
				tsk.name = string(tv)
			case slip.Symbol:
				tsk.name = string(tv)
			default:
				slip.PanicType("task :init :name", args[i+1], "string", "symbol")
			}
		case slip.Symbol(":actor"):
			val := args[i+1]
		Actor:
			switch tv := val.(type) {
			case slip.Instance:
				if tv.HasMethod(":perform") {
					tsk.actors = []slip.Instance{tv}
				} else {
					slip.PanicType("task :init :actor", val, "instance with :perform method", "function", "list")
				}
			case slip.List:
				tsk.actors = make([]slip.Instance, len(tv))
				for j, sv := range tv {
					if a, ok := sv.(slip.Instance); ok && a.HasMethod(":perform") {
						tsk.actors[j] = a
					} else {
						slip.PanicType("task :init :actor", val, "instance with :perform method", "function", "list")
					}
				}
				// TBD could be lambda?
			case *slip.Lambda:
				tsk.caller = tv
			case *slip.FuncInfo:
				tsk.caller = tv.Create(nil).(slip.Funky).Caller()
			case slip.Symbol:
				val = slip.FindFunc(string(tv))
				goto Actor
			default:
				slip.PanicType("task :init :actor", val, "instance with :perform method", "function", "list")
			}
		case slip.Symbol(":workers"):
			if num, ok := args[i+1].(slip.Fixnum); ok {
				tsk.workers = int(num)
			} else {
				slip.PanicType("task :init :workers", args[i+1], "fixnum")
			}
		case slip.Symbol(":depth"):
			if num, ok := args[i+1].(slip.Fixnum); ok && 0 < num {
				tsk.depth = int(num)
			} else {
				slip.PanicType("task :init :depth", args[i+1], "fixnum greater than 0")
			}
		default:
			slip.PanicType("task :init", args[i], ":name", ":actor", ":workers")
		}
	}
	obj.Any = &tsk

	return nil
}

func (caller taskInitCaller) Docs() string {
	return `__:init__ &key _name_ _workers_ _actor_
   _:name_ [string] sets the name of the task.
   _:workers_ [fixnum] the number of workers for concurrent processing. Zero indicates no concurrent processing.
   _:depth_ [fixnum] of the work queue.
   _:actor_ [instance|function|list] if an instance that instance is used for processing and must have the
_perform_ method that expectes an instance of the _flow-box-flavor_. If the instance has a _start_ or _shutdown_
those will be called when starting or stoping a flow. If the value of _:actor_ is a function is must expect one
box argument just as the _:perform_ method does. If the actor is a list of instances those will be used as
workers.


Sets the initial value when _make-instance_ is called.
`
}

// MakeTask is only public for testing purposes.
func MakeTask(id slip.Object) (self *flavors.Instance, t *task) {
	self = taskFlavor.MakeInstance().(*flavors.Instance)
	t = &task{links: map[string]*link{}}
	self.Any = t

	return
}
