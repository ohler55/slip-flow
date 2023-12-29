// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

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
				slip.Symbol(":depth"),
			},
		},
	)
	taskFlavor.DefMethod(":init", "", taskInitCaller{})
	taskFlavor.DefMethod(":name", "", taskNameCaller{})
	taskFlavor.DefMethod(":workers", "", taskWorkersCaller{})
	taskFlavor.DefMethod(":start", "", taskStartCaller{})
	taskFlavor.DefMethod(":shutdown", "", taskShutdownCaller{})
	taskFlavor.DefMethod(":running", "", taskRunningCaller{})
	taskFlavor.DefMethod(":receive", "", taskReceiveCaller{})
	taskFlavor.DefMethod(":metrics", "", taskMetricsCaller{})

	// :transition
	//  reset/zero metrics
	// TBD
}

type task struct {
	name      string
	self      *flavors.Instance
	flow      *flow
	links     map[string]*task
	actors    []slip.Instance
	caller    slip.Caller
	queue     chan *flavors.Instance // must be box instances
	done      chan struct{}
	workers   int
	depth     int
	qmu       sync.Mutex
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
		t.qmu.Lock()
		t.queue = make(chan *flavors.Instance, t.depth)
		t.done = make(chan struct{}, t.depth)
		t.qmu.Unlock()
		if 0 < len(t.actors) {
			for i := t.workers; 0 < i; i-- {
				go t.actorLoop(s, t.actors[i%len(t.actors)])
			}
		} else {
			for i := t.workers; 0 < i; i-- {
				go t.funcLoop(s)
			}
		}
	}
}

func (t *task) shutdown(s *slip.Scope) {
	t.qmu.Lock()
	defer t.qmu.Unlock()
	if t.queue != nil {
		for i := t.workers; 0 < i; i-- {
			t.queue <- nil
		}
		close(t.queue)
		for i := t.workers; 0 < i; i-- {
			<-t.done
		}
		close(t.done)
		t.queue = nil
		t.done = nil
	}
	for _, a := range t.actors {
		if a.HasMethod(":shutdown") {
			a.Receive(s, ":shutdown", slip.List{}, 0)
		}
	}
}

func (t *task) running() bool {
	t.qmu.Lock()
	defer t.qmu.Unlock()

	return t.queue != nil
}

func (t *task) actorLoop(s *slip.Scope, actor slip.Instance) {
	for {
		bi := <-t.queue
		if bi == nil {
			break
		}
		t.act(s, actor, bi)
	}
	t.done <- struct{}{}
}

func (t *task) funcLoop(s *slip.Scope) {
	for {
		bi := <-t.queue
		if bi == nil {
			break
		}
		t.call(s, bi)
	}
	t.done <- struct{}{}
}

func (t *task) receive(s *slip.Scope, bi *flavors.Instance) {
	var (
		bx       *box
		flowName string
	)
	t.received.Add(1)
	bi, bx = boxDup(bi)
	if t.flow != nil {
		flowName = t.flow.name
	}
	bx.track.Scan(flowName, t.name)
	switch {
	case t.queue != nil:
		t.queue <- bi
	case 0 < len(t.actors):
		t.act(s, t.actors[0], bi)
	case t.caller != nil:
		t.call(s, bi)
	}
}

func (t *task) act(s *slip.Scope, actor slip.Instance, bi *flavors.Instance) {
	defer t.handlePanic(s, bi)
	t.handleResult(s, actor.Receive(s, ":perform", slip.List{bi}, 0))
}

func (t *task) call(s *slip.Scope, bi *flavors.Instance) {
	defer t.handlePanic(s, bi)
	t.handleResult(s, t.caller.Call(s, slip.List{bi}, 0))
}

func (t *task) handleResult(s *slip.Scope, result slip.Object) {
	if list, _ := result.(slip.List); len(list) == 2 {
		var to *task
		switch tr := list[0].(type) {
		case nil:
			to = t.links[""]
		case slip.String:
			to = t.links[string(tr)]
		case slip.Symbol:
			to = t.links[string(tr)]
		}
		if bi, has := list[1].(*flavors.Instance); has && bi != nil && boxFlavor == bi.Flavor {
			t.processed.Add(1)
			tr := bi.Any.(*box).track
			ev := tr.history[len(tr.history)-1]
			t.duration.Add(uint64(time.Since(ev.when)))
			if to != nil {
				to.receive(s, bi)
			}
			return
		}
	}
	slip.NewPanic("Actor in task %s did not return a list of transition name and box instance.", t.name)
}

func (t *task) handlePanic(s *slip.Scope, bi *flavors.Instance) {
	if rec := recover(); rec != nil {
		t.errors.Add(1)
		tr := bi.Any.(*box).track
		ev := tr.history[len(tr.history)-1]
		t.duration.Add(uint64(time.Since(ev.when)))
		nb, bx := boxDup(bi)
		if rs, ok := rec.(fmt.Stringer); ok {
			bx.content = map[string]any{
				"content": bx.content,
				"error":   rs.String(),
			}
		} else {
			bx.content = map[string]any{
				"content": bx.content,
				"error":   fmt.Sprintf("%v", rs),
			}
		}
		if to, has := t.links["error"]; has {
			if to != nil {
				to.receive(s, nb)
			}
			return
		}
		if t.flow != nil && t.flow.errorTask != nil {
			t.flow.errorTask.receive(s, nb)
		}
	}
}

func (t *task) metrics() (alist slip.List) {
	alist = append(alist, slip.List{slip.Symbol("received"), slip.Tail{Value: slip.Fixnum(t.received.Load())}})
	ecnt := t.errors.Load()
	pcnt := t.processed.Load()
	dur := t.duration.Load()
	alist = append(alist, slip.List{slip.Symbol("processed"), slip.Tail{Value: slip.Fixnum(pcnt)}})
	alist = append(alist, slip.List{slip.Symbol("errors"), slip.Tail{Value: slip.Fixnum(ecnt)}})
	if 0 < pcnt || 0 < ecnt {
		alist = append(alist,
			slip.List{
				slip.Symbol("average"),
				slip.Tail{Value: slip.DoubleFloat(float64(dur) / float64(time.Second) / float64(ecnt+pcnt))},
			})
	}
	return
}

// MakeTask is only public for testing purposes.
func MakeTask(name slip.Object) (self *flavors.Instance, t *task) {
	self = taskFlavor.MakeInstance().(*flavors.Instance)
	t = &task{links: map[string]*task{}}
	switch tn := name.(type) {
	case slip.Symbol:
		t.name = string(tn)
	case slip.String:
		t.name = string(tn)
	}
	self.Any = t

	return
}

// task-flavor :init //////////////////////////////////////////////////////////

type taskInitCaller struct{}

func (caller taskInitCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		args = args[0].(slip.List)
	}
	tsk := task{self: obj, links: map[string]*task{}}
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
