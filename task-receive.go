// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := TaskReceive{Function: slip.Function{Name: "flow-task-receive", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-task-receive",
			Args: []*slip.DocArg{
				{
					Name: "task",
					Type: "flow-task",
					Text: "to receive a box.",
				},
				{
					Name: "box",
					Type: "instance",
					Text: "to receive.",
				},
			},
			Return: "nil",
			Text:   `__flow-task-receive__ receives a _box_ to be processed by a _task_.`,
			Examples: []string{
				`(setq task (make-instance 'flow-task-flavor :name "tisk"))`,
				`(flow-task-receive task (make-flow-box :parse "[1 2]") => nil`,
			},
		}, &Pkg)
}

// TaskReceive represents the flow-task-receive function.
type TaskReceive struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *TaskReceive) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.ArgCountCheck(f, args, 2, 2)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Flavor != taskFlavor {
		slip.PanicType("task", args[0], "task")
	}
	var bi *flavors.Instance
	if bi, ok = args[1].(*flavors.Instance); !ok || boxFlavor != bi.Flavor {
		slip.PanicType("box", args[1], "box")
	}
	self.Any.(*task).receive(s, bi)

	return nil
}

type taskReceiveCaller struct{}

func (caller taskReceiveCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	bi, ok := args[0].(*flavors.Instance)
	if !ok || boxFlavor != bi.Flavor {
		slip.PanicType("box", args[0], "box")
	}
	obj.Any.(*task).receive(s, bi)

	return nil
}

func (caller taskReceiveCaller) Docs() string {
	return methodDocFromFunc(":receive", "flow-task-receive", "flow-task-flavor", "task")
}
