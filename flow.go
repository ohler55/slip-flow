// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main

import (
	"github.com/ohler55/slip"
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


Once built data is submitted as an instance of the _box-flavor_ which collects
tracking information as it traverses the graph of linked tasks. If provided
the final tasks in a flow will place the _box_ with tracking information on an
exit channel.


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
	// flowFlavor.DefMethod(":name", "", flowNameCaller{})
	// flowFlavor.DefMethod(":start", "", flowStartCaller{})
	// flowFlavor.DefMethod(":shutdown", "", flowShutdownCaller{})
	// flowFlavor.DefMethod(":running", "", flowRunningCaller{})
	// flowFlavor.DefMethod(":metrics", "", flowMetricsCaller{})
	// flowFlavor.DefMethod(":add-task", "", flowAddTaskCaller{})
	// flowFlavor.DefMethod(":link", "", flowLinkCaller{})
	// flowFlavor.DefMethod(":unlink", "", flowUnlinkCaller{})
	// flowFlavor.DefMethod(":remove-task", "", flowRemoveTaskCaller{})
	// flowFlavor.DefMethod(":tasks", "", flowTasksCaller{})
	// flowFlavor.DefMethod(":FindTask", "", flowFindTaskCaller{})
	// flowFlavor.DefMethod(":entry", "", flowEntryCaller{})
	// flowFlavor.DefMethod(":set-entry", "", flowSetEntryCaller{})
	// flowFlavor.DefMethod(":submit", "", flowSubmitCaller{})
	// TBD
}

type flow struct {
	name      string
	self      *flavors.Instance
	tasks     map[string]*task
	errorTask *task
	exitChan  gi.Channel

	// received  atomic.Uint64
	// errors    atomic.Uint64
	// processed atomic.Uint64
	// duration  atomic.Uint64 // sum of processing times from entry to completed from box tracks
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
