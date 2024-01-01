// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := TaskName{Function: slip.Function{Name: "flow-task-name", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-task-name",
			Args: []*slip.DocArg{
				{
					Name: "task",
					Type: "flow-task",
					Text: "to return the name of.",
				},
			},
			Return: "string",
			Text:   `__flow-task-name__ returns the name of a _flow-task-flavor_ instance.`,
			Examples: []string{
				`(setq task (make-instance 'flow-task-flavor :name "tisk"))`,
				`(flow-task-name task) => "tisk"`,
			},
		}, &Pkg)
}

// TaskName represents the flow-task-name function.
type TaskName struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *TaskName) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	slip.ArgCountCheck(f, args, 1, 1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Flavor != taskFlavor {
		slip.PanicType("task", args[0], "task")
	}
	return slip.String(self.Any.(*task).name)
}

type taskNameCaller struct{}

func (caller taskNameCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	return slip.String(obj.Any.(*task).name)
}

func (caller taskNameCaller) Docs() string {
	return methodDocFromFunc(":name", "flow-task-name", "flow-task-flavor", "task")
}
