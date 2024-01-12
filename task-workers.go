// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := TaskWorkers{Function: slip.Function{Name: "flow-task-workers", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-task-workers",
			Args: []*slip.DocArg{
				{
					Name: "task",
					Type: "flow-task",
					Text: "to return the number of workers from.",
				},
			},
			Return: "fixnum",
			Text:   `__flow-task-workers__ returns the number of workers for a _flow-task-flavor_ instance.`,
			Examples: []string{
				`(setq task (make-instance 'flow-task-flavor :name "tisk" :workers 3))`,
				`(flow-task-workers task) => 3`,
			},
		}, &Pkg)
}

// TaskWorkers represents the flow-task-workers function.
type TaskWorkers struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *TaskWorkers) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	slip.ArgCountCheck(f, args, 1, 1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Flavor != taskFlavor {
		slip.PanicType("task", args[0], "task")
	}
	return slip.Fixnum(self.Any.(*task).workers)
}

type taskWorkersCaller struct{}

func (caller taskWorkersCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	return slip.Fixnum(obj.Any.(*task).workers)
}

func (caller taskWorkersCaller) Docs() string {
	return methodDocFromFunc(":workers", "flow-task-workers", "flow-task-flavor", "task")
}
