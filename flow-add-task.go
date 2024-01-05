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

// FlowAddTask represents the flow-flow-add-task function.
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
	self.Any.(*flow).addTask(args[1:])

	return nil
}

type flowAddTaskCaller struct{}

func (caller flowAddTaskCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)

	obj.Any.(*flow).addTask(args)

	return nil
}

func (caller flowAddTaskCaller) Docs() string {
	return methodDocFromFunc(":add-task", "flow-add-task", "flow-flavor", "flow")
}
