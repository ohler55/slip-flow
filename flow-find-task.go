// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := FlowFindTask{Function: slip.Function{Name: "flow-find-task", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-find-task",
			Args: []*slip.DocArg{
				{
					Name: "flow",
					Type: "flow",
					Text: "to find a task in.",
				},
				{
					Name: "task-name",
					Type: "string",
					Text: "name of the task to find.",
				},
			},
			Return: "instance",
			Text:   `__flow-find-task__ finds and returns the task with a name matching _task-name_.`,
			Examples: []string{
				`(setq flow (make-instance 'flow-flavor :find-task "flo"))`,
				`(flow-find-task flow :name 'tisk :actor (lambda (b) (list 'ok b)))`,
				`(flow-find-task flow "tisk") => #<flow-task-flavor 12345>`,
			},
		}, &Pkg)
}

// FlowFindTask represents the flow-find-task function.
type FlowFindTask struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *FlowFindTask) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.ArgCountCheck(f, args, 2, 2)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Flavor != flowFlavor {
		slip.PanicType("flow", args[0], "flow")
	}
	return self.Any.(*flow).findTask(args[1])
}

type flowFindTaskCaller struct{}

func (caller flowFindTaskCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)

	return obj.Any.(*flow).findTask(args[0])
}

func (caller flowFindTaskCaller) Docs() string {
	return methodDocFromFunc(":find-task", "flow-find-task", "flow-flavor", "flow")
}
