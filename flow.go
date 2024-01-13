// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

import (
	"sort"
	"sync/atomic"
	"time"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/bag"
	"github.com/ohler55/slip/pkg/flavors"
	"github.com/ohler55/slip/pkg/gi"
)

var (
	flowFlavor *flavors.Flavor
)

func init() {
	flowFlavor = flavors.DefFlavor("flow-flavor",
		map[string]slip.Object{},
		[]string{
			"can-log-flavor",
		},
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`

Flows are named and are a container for tasks. There are some elements of the
flow that are shared across tasks. One shared element is the log level and an
error handler task. If set the flow error task is transitioned to on errors.


An optional entry task can be identified as the starting point for
processing. Data submitted to a flow is passed to the entry task. The entry
task is optional if a trigger task is included in the flow. A trigger taask is
one that generates a data _box_ when an event occurs such as receiving an HTTP
request or a timer triggers.


A flow is built by first adding tasks to the flow and then linking the tasks
together. Designating an entry task if needed as well. The "error" task name
is reserved for handling errors and panics.


Once a flow has been built data is submitted as an instance of the
_box-flavor_ which collects tracking information as it traverses the graph of
linked tasks. If provided the final tasks in a flow will place the _box_ with
tracking information on an exit channel.


See also: flow-task-flavor

`),
			},
			slip.List{
				slip.Symbol(":init-keywords"),
				slip.Symbol(":name"),
				slip.Symbol(":exit-channel"),
			},
		},
	)
	flowFlavor.DefMethod(":init", "", flowInitCaller{})
	flowFlavor.DefMethod(":name", "", flowNameCaller{})
	flowFlavor.DefMethod(":start", "", flowStartCaller{})
	flowFlavor.DefMethod(":shutdown", "", flowShutdownCaller{})
	flowFlavor.DefMethod(":running", "", flowRunningCaller{})
	flowFlavor.DefMethod(":add-task", "", flowAddTaskCaller{})
	flowFlavor.DefMethod(":tasks", "", flowTasksCaller{})
	flowFlavor.DefMethod(":remove-task", "", flowRemoveTaskCaller{})
	flowFlavor.DefMethod(":find-task", "", flowFindTaskCaller{})
	flowFlavor.DefMethod(":entry", "", flowEntryCaller{})
	flowFlavor.DefMethod(":set-entry", "", flowSetEntryCaller{})
	flowFlavor.DefMethod(":link", "", flowLinkCaller{})
	flowFlavor.DefMethod(":submit", "", flowSubmitCaller{})
	flowFlavor.DefMethod(":exit-channel", "", flowExitChannelCaller{})
	flowFlavor.DefMethod(":metrics", "", flowMetricsCaller{})
	flowFlavor.DefMethod(":reset-metrics", "", flowResetMetricsCaller{})
	flowFlavor.DefMethod(":set-level", ":after", flowSetLevelCaller{})
}

type flow struct {
	name     string
	self     *flavors.Instance
	group    *group
	tasks    map[string]*task
	entry    *task
	exitChan gi.Channel
	started  bool

	received  atomic.Uint64
	errors    atomic.Uint64
	processed atomic.Uint64
	duration  atomic.Uint64 // sum of processing times from entry to completed from box tracks
}

func (f *flow) start(s *slip.Scope) {
	f.received.Store(0)
	f.errors.Store(0)
	f.processed.Store(0)
	f.duration.Store(0)
	logger := f.self.Get("logger")
	if logger == nil {
		logger = slip.ReadString("(make-instance 'logger-flavor)").Eval(s, nil)
		f.self.Set("logger", logger)
	}
	for _, t := range f.tasks {
		t.self.Set("logger", logger)
		t.start(s)
	}
	f.started = true
}

func (f *flow) shutdown(s *slip.Scope) {
	for _, t := range f.tasks {
		t.shutdown(s)
	}
	f.started = false
}

func (f *flow) running() bool {
	return f.started
}

func (f *flow) addTask(args slip.List) *flavors.Instance {
	inst, tsk := MakeTask(args...)
	if _, has := f.tasks[tsk.name]; has {
		slip.NewPanic("Task %s already exists in flow %s.", tsk.name, f.name)
	}
	tsk.flow = f
	f.tasks[tsk.name] = tsk

	return inst
}

func (f *flow) removeTask(name slip.Object) {
	var key string
	switch tn := name.(type) {
	case slip.String:
		key = string(tn)
	case slip.Symbol:
		key = string(tn)
	default:
		slip.PanicType("flow :remove-task :task", tn, "string", "symbol")
	}
	delete(f.tasks, key)
}

func (f *flow) findTask(name slip.Object) (found slip.Object) {
	var key string
	switch tn := name.(type) {
	case slip.String:
		key = string(tn)
	case slip.Symbol:
		key = string(tn)
	default:
		slip.PanicType("flow :find-task :task-name", tn, "string", "symbol")
	}
	if t := f.tasks[key]; t != nil {
		found = t.self
	}
	return
}

func (f *flow) taskList() slip.List {
	tasks := make(slip.List, 0, len(f.tasks))
	for _, t := range f.tasks {
		tasks = append(tasks, t.self)
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].(*flavors.Instance).Any.(*task).name < tasks[j].(*flavors.Instance).Any.(*task).name
	})
	return tasks
}

func (f *flow) setEntry(name slip.Object) (found slip.Object) {
	var key string
	switch tn := name.(type) {
	case nil:
		f.entry = nil
		return
	case slip.String:
		key = string(tn)
	case slip.Symbol:
		key = string(tn)
	default:
		slip.PanicType("flow :set-entry :task-name", tn, "string", "symbol")
	}
	if t := f.tasks[key]; t != nil {
		found = t.self
		f.entry = t
	} else {
		slip.NewPanic("task %s not found", key)
	}
	return
}

func (f *flow) link(args slip.List) {
	// Argument count already checked.
	var (
		from *task
		to   *task
	)
	name := strFromArg(args[0], "flow :link :link-name")
	if from = f.tasks[strFromArg(args[1], "flow :link :from")]; from == nil {
		slip.NewPanic("task %s not found", args[1])
	}
	if to = f.tasks[strFromArg(args[2], "flow :link :to")]; to == nil {
		slip.NewPanic("task %s not found", args[2])
	}
	from.links[name] = to
}

func (f *flow) exit(bi slip.Object) {
	history := bi.(*flavors.Instance).Any.(*box).track.history
	if 0 < len(history) {
		first := history[0]
		last := history[len(history)-1]
		if last.task == "error" {
			f.errors.Add(1)
		} else {
			f.processed.Add(1)
			f.duration.Add(uint64(last.when.Sub(first.when)))
		}
	}
	if f.exitChan != nil {
		f.exitChan <- bi
	}
}

func (f *flow) submit(s *slip.Scope, data slip.Object) {
	if f.entry == nil {
		slip.NewPanic("no entry task has been set for the %s flow", f.name)
	}
	if !f.started {
		f.start(s)
	}
	var bi *flavors.Instance // box-flavor
	if inst, _ := data.(*flavors.Instance); inst != nil {
		switch {
		case inst.Flavor == boxFlavor:
			bi = inst
		case inst.Flavor == bag.Flavor():
			var bx *box
			bi, bx = MakeBox(gi.NewUUID())
			bx.content = inst.Any
			bx.frozen = true
		default:
			slip.PanicType("box", data, "flow-box-flavor", "bag-flavor")
		}
	} else {
		var bx *box
		bi, bx = MakeBox(gi.NewUUID())
		bx.content = bag.ObjectToBag(data)
	}
	f.entry.receive(s, bi)
}

func (f *flow) metrics() (alist slip.List) {
	alist = append(alist, slip.List{slip.Symbol("received"), slip.Tail{Value: slip.Fixnum(f.received.Load())}})
	ecnt := f.errors.Load()
	pcnt := f.processed.Load()
	dur := f.duration.Load()
	alist = append(alist, slip.List{slip.Symbol("processed"), slip.Tail{Value: slip.Fixnum(pcnt)}})
	alist = append(alist, slip.List{slip.Symbol("errors"), slip.Tail{Value: slip.Fixnum(ecnt)}})
	if 0 < pcnt {
		alist = append(alist,
			slip.List{
				slip.Symbol("average"),
				slip.Tail{Value: slip.DoubleFloat(float64(dur) / float64(time.Second) / float64(pcnt))},
			})
	}
	return
}

func (f *flow) resetMetrics() {
	f.received.Store(0)
	f.errors.Store(0)
	f.processed.Store(0)
	f.duration.Store(0)
	for _, t := range f.tasks {
		t.resetMetrics()
	}
}

func strFromArg(arg slip.Object, argName string) (str string) {
	switch ta := arg.(type) {
	case nil:
		str = ""
	case slip.String:
		str = string(ta)
	case slip.Symbol:
		str = string(ta)
	default:
		slip.PanicType(argName, ta, "string", "symbol")
	}
	return
}

type flowInitCaller struct{}

func (caller flowInitCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		args = args[0].(slip.List)
	}
	flo := flow{self: obj, tasks: map[string]*task{}}
	for i := 0; i < len(args)-1; i += 2 {
		switch args[i] {
		case slip.Symbol(":name"):
			switch tv := args[i+1].(type) {
			case slip.String:
				flo.name = string(tv)
			case slip.Symbol:
				flo.name = string(tv)
			default:
				slip.PanicType("flow :init :name", args[i+1], "string", "symbol")
			}
		case slip.Symbol(":exit-channel"):
			if ch, ok := args[i+1].(gi.Channel); ok {
				flo.exitChan = ch
			} else {
				slip.PanicType("flow :init :exit-channel", args[i+1], "gi:chanel")
			}
		}
	}
	obj.Any = &flo

	return nil
}

func (caller flowInitCaller) Docs() string {
	return `__:init__ &key _name_ _exit-channel_ _logger_
   _:name_ [string] sets the name of the flow.
   _:exit-channel_ [gi:channel] if provided the exit tasks of a flow place the _box_ being processed on this channel.
   _:logger_ [instance] an instance that has the _:log_ method.


Sets the initial value when _make-instance_ is called.
`
}

type flowSetLevelCaller struct{}

func (caller flowSetLevelCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	f := obj.Any.(*flow)
	level := f.self.Get("log-level")
	for _, t := range f.tasks {
		_ = t.self.Receive(s, ":set-level", slip.List{level}, 0)
	}
	return nil
}

func (caller flowSetLevelCaller) Docs() string {
	return `__:after :setLevel__


Sets the _log-level_ of all the tasks in the flow.
`
}
