// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"fmt"
	"sort"
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
	Pkg.Initialize(nil)
	taskFlavor = flavors.DefFlavor("flow-task",
		map[string]slip.Object{
			"x":   nil,
			"y":   nil,
			"svg": nil,
		},
		[]string{
			"can-log",
		},
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`Tasks are the shell around an actor that performs specific actions on a _box_
that is passed from one _task_ to another in a _flow_.


Each _task_ is named and can process a _box_ either synchronously by setting
the _:workers_ to zero or asynchronous if _:workers_ is set to one or
more. After processing the _box_ is sent through a link to the destination
_task_ at the end of the link.


See also: flow

`),
			},
			slip.List{
				slip.Symbol(":init-keywords"),
				slip.Symbol(":name"),
				slip.Symbol(":actor"),
				slip.Symbol(":workers"),
				slip.Symbol(":depth"),
			},
			slip.Symbol(":gettable-instance-variables"),
			slip.Symbol(":settable-instance-variables"),
			slip.Symbol(":inittable-instance-variables"),
		},
		&Pkg,
	)
	taskFlavor.DefMethod(":init", "", taskInitCaller{})
	taskFlavor.DefMethod(":name", "", taskNameCaller{})
	taskFlavor.DefMethod(":actors", "", taskActorsCaller{})
	taskFlavor.DefMethod(":workers", "", taskWorkersCaller{})
	taskFlavor.DefMethod(":depth", "", taskDepthCaller{})
	taskFlavor.DefMethod(":start", "", taskStartCaller{})
	taskFlavor.DefMethod(":shutdown", "", taskShutdownCaller{})
	taskFlavor.DefMethod(":running", "", taskRunningCaller{})
	taskFlavor.DefMethod(":receive", "", taskReceiveCaller{})
	taskFlavor.DefMethod(":metrics", "", taskMetricsCaller{})
	taskFlavor.DefMethod(":reset-metrics", "", taskResetMetricsCaller{})
	taskFlavor.DefMethod(":links", "", taskLinksCaller{})
	taskFlavor.DefMethod(":unlink", "", taskUnlinkCaller{})
	taskFlavor.DefMethod(":transition", "", taskTransitionCaller{})
	taskFlavor.DefMethod(":update-link", "", taskUpdateLinkCaller{})
}

type link struct {
	task *task
	mids slip.List
}

type task struct {
	name      string
	self      *flavors.Instance
	flow      *flow
	links     map[string]*link
	actors    []slip.Instance
	caller    slip.Caller
	funcName  string
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
		if t.depth <= 0 {
			t.depth = t.workers * 2
		}
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
	if t.flow != nil {
		flowName = t.flow.name
	}
	if levelInfo <= int(t.self.Get("log-level").(slip.Fixnum)) {
		msg := fmt.Sprintf("%s:%s received box %s", flowName, t.name, bi.Any.(*box).track.id)
		t.self.Receive(s, ":info", slip.List{slip.String(msg)}, 0)
	}
	bi, bx = boxDup(bi)
	if len(bx.track.history) == 0 && t.flow != nil {
		t.flow.received.Add(1)
	}
	t.received.Add(1)
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
		var linkName string
		switch tr := list[0].(type) {
		case nil:
			if list[1] == nil && (list[0] == nil || len(t.links) == 0) {
				// No links and both the transition and box are nil for this
				// is the end of the line.
				return
			}
		case slip.String:
			linkName = string(tr)
		case slip.Symbol:
			linkName = string(tr)
		}
		if bi, has := list[1].(*flavors.Instance); has && bi != nil && boxFlavor == bi.Flavor {
			t.transition(s, linkName, bi)
			return
		}
	}
	slip.NewPanic("Actor did not return a list of link name and box instance.")
}

func (t *task) transition(s *slip.Scope, linkName string, bi *flavors.Instance) {
	var to *task
	if lnk := t.links[linkName]; lnk != nil {
		to = lnk.task
	}
	t.processed.Add(1)
	tr := bi.Any.(*box).track
	ev := tr.history[len(tr.history)-1]
	t.duration.Add(uint64(time.Since(ev.when)))
	if levelInfo <= int(t.self.Get("log-level").(slip.Fixnum)) {
		var flowName string
		if t.flow != nil {
			flowName = t.flow.name
		}
		msg := fmt.Sprintf("%s:%s following %s with box %s",
			flowName, t.name, linkName, bi.Any.(*box).track.id)
		t.self.Receive(s, ":info", slip.List{slip.String(msg)}, 0)
	}
	if to != nil {
		to.receive(s, bi)
	}
}

func (t *task) handlePanic(s *slip.Scope, bi *flavors.Instance) {
	if rec := recover(); rec != nil {
		t.handleError(s, bi, rec)
	}
}

func (t *task) handleError(s *slip.Scope, bi *flavors.Instance, err any) {
	t.errors.Add(1)
	tr := bi.Any.(*box).track
	ev := tr.history[len(tr.history)-1]
	t.duration.Add(uint64(time.Since(ev.when)))
	nb, bx := boxDup(bi)
	msg := fmt.Sprintf("%v", err)
	if se, _ := err.(slip.Error); se != nil {
		msg = se.Error()
	}
	bx.content = map[string]any{
		"content": bx.content,
		"error":   msg,
	}
	if lnk := t.links["error"]; lnk != nil {
		lnk.task.receive(s, nb)
		return
	}
	if t.flow != nil {
		if et := t.flow.tasks["error"]; et != nil {
			et.receive(s, nb)
			return
		}
		msg := fmt.Sprintf("%s:%s box %s: %s",
			t.flow.name, t.name, bi.Any.(*box).track.id, bx.content.(map[string]any)["error"])
		t.self.Receive(s, ":error", slip.List{slip.String(msg)}, 0)
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

func (t *task) resetMetrics() {
	t.received.Store(0)
	t.errors.Store(0)
	t.processed.Store(0)
	t.duration.Store(0)
}

func (t *task) linkList() (la slip.List) {
	if 0 < len(t.links) {
		keys := make([]string, 0, len(t.links))
		for k := range t.links {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		la = make(slip.List, len(keys))
		for i, k := range keys {
			var ti *flavors.Instance
			lnk := t.links[k]
			if lnk != nil {
				ti = lnk.task.self
			}
			la[i] = append(slip.List{slip.String(k), ti}, lnk.mids...)
		}
	}
	return
}

func (t *task) actorList() (al slip.List) {
	switch {
	case 0 < len(t.actors):
		al = make(slip.List, len(t.actors))
		for i, v := range t.actors {
			al[i] = v
		}
	case 0 < len(t.funcName):
		al = append(al, slip.Symbol(t.funcName))
	case t.caller != nil:
		if obj, ok := t.caller.(*slip.Lambda); ok {
			al = append(al, obj)
		}
	}
	return
}

func (t *task) unlink(args slip.List) {
	var name string
	switch ta := args[0].(type) {
	case nil:
		// leave name as ""
	case slip.String:
		name = string(ta)
	case slip.Symbol:
		name = string(ta)
	default:
		slip.PanicType("link", ta, "string", "symbol")
	}
	t.qmu.Lock()
	delete(t.links, name)
	t.qmu.Unlock()
}

func (t *task) updateLink(args slip.List) {
	var (
		name string
		mids slip.List
	)
	switch ta := args[0].(type) {
	case nil:
		// leave name as ""
	case slip.String:
		name = string(ta)
	case slip.Symbol:
		name = string(ta)
	default:
		slip.PanicType("link-name", ta, "string", "symbol")
	}
	switch ta := args[1].(type) {
	case nil:
		// leave empty or nil
	case slip.List:
		mids = checkMidPoints(ta)
	default:
		slip.PanicType("mid-points", ta, "list")
	}
	t.qmu.Lock()
	defer t.qmu.Unlock()
	if lnk := t.links[name]; lnk != nil {
		lnk.mids = mids
	} else {
		slip.NewPanic("task %s has no %s link", t.name, name)
	}
}

func (t *task) validate(s *slip.Scope) (fails slip.List) {
	allowed := map[string]bool{}
	for _, a := range t.actors {
		if a.HasMethod(":links") {
			al, _ := a.Receive(s, ":links", slip.List{}, 0).(slip.List)
			for _, v := range al {
				if name, ok := v.(slip.String); ok {
					allowed[string(name)] = true
				}
			}
		}
	}
	if 0 < len(allowed) {
		for name := range t.links {
			if !allowed[name] {
				fails = append(fails, slip.String(fmt.Sprintf("task %s can not return a %s link", t.name, name)))
			}
		}
	}
	return
}

// MakeTask is only public for testing purposes.
func MakeTask(args ...slip.Object) (self *flavors.Instance, t *task) {
	self = taskFlavor.MakeInstance().(*flavors.Instance)
	t = makeTaskStruct(self, args)
	return
}

func makeTaskStruct(self *flavors.Instance, args slip.List) (tsk *task) {
	tsk = &task{self: self, links: map[string]*link{}}
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
			case *slip.Lambda:
				tsk.caller = tv
			case *slip.FuncInfo:
				tsk.caller = tv.Create(nil).(slip.Funky).Caller()
				tsk.funcName = tv.Name
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
		case slip.Symbol(":x"):
			if args[i+1] != nil {
				if num, ok := args[i+1].(slip.Fixnum); ok {
					self.Let("x", num)
				} else {
					slip.PanicType("task :init :x", args[i+1], "fixnum")
				}
			}
		case slip.Symbol(":y"):
			if args[i+1] != nil {
				if num, ok := args[i+1].(slip.Fixnum); ok {
					self.Let("y", num)
				} else {
					slip.PanicType("task :init :x", args[i+1], "fixnum")
				}
			}
		case slip.Symbol(":svg"):
			if args[i+1] != nil {
				if str, ok := args[i+1].(slip.String); ok {
					self.Let("svg", str)
				} else {
					slip.PanicType("task :init :svg", args[i+1], "string")
				}
			}
		}
	}
	self.Any = tsk

	return
}

type taskInitCaller struct{}

func (caller taskInitCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	self := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		args = args[0].(slip.List)
	}
	_ = makeTaskStruct(self, args)

	return nil
}

func (caller taskInitCaller) Docs() string {
	return `__:init__ &key _name_ _workers_ _actor_ _logger_
   _:name_ [string] sets the name of the task.
   _:workers_ [fixnum] the number of workers for concurrent processing. Zero indicates no concurrent processing.
   _:depth_ [fixnum] of the work queue.
   _:logger_ [instance] an instance that has the _:log_ method.
   _:actor_ [instance|function|list] if an instance that instance is used for processing and must have the
_perform_ method that expectes an instance of the _flow-box_. If the instance has a _start_ or _shutdown_
those will be called when starting or stoping a flow. If the value of _:actor_ is a function is must expect one
box argument just as the _:perform_ method does. If the actor is a list of instances those will be used as
workers.


Sets the initial value when _make-instance_ is called.
`
}
