// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := FlowAddTask{Function: slip.Function{Name: "flow-add-task", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-add-task",
			Args: []*slip.DocArg{
				{
					Name: "flow",
					Type: "flow",
					Text: "to add-task.",
				},
				{Name: "&key"},
				{
					Name: "name",
					Type: "string",
					Text: "of task to add.",
				},
				{
					Name: "actor",
					Type: "instance|function|list",
					Text: `if an instance that instance is used for processing and must have the
_perform_ method that expectes an instance of the _flow-box-flavor_. If the instance has a _start_ or _shutdown_
those will be called when starting or stoping a flow. If the value of _:actor_ is a function is must expect one
box argument just as the _:perform_ method does. If the actor is a list of instances those will be used as
workers.`,
				},
				{
					Name: "workers",
					Type: "fixnum",
					Text: "the number of workers for concurrent processing. Zero indicates no concurrent processing.",
				},
				{
					Name: "depth",
					Type: "fixnum",
					Text: "of the work queue.",
				},
			},
			Return: "nil",
			Text:   `__flow-add-task__ add-tasks all the _tasks_ in the _flow_.`,
			Examples: []string{
				`(setq flow (make-instance 'flow-flavor :add-task "flo"))`,
				`(flow-add-task flow :name 'tisk :actor (lambda (b) (list 'ok b)))`,
				`(flow-add-task flow) => nil`,
			},
		}, &Pkg)
}

// FlowAddTask represents the flow-add-task function.
type FlowAddTask struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *FlowAddTask) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.ArgCountCheck(f, args, 1, -1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Flavor != flowFlavor {
		slip.PanicType("flow", args[0], "flow")
	}
	return self.Any.(*flow).addTask(args[1:])
}

type flowAddTaskCaller struct{}

func (caller flowAddTaskCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)

	return obj.Any.(*flow).addTask(args)
}

func (caller flowAddTaskCaller) Docs() string {
	return methodDocFromFunc(":add-task", "flow-add-task", "flow-flavor", "flow")
}
