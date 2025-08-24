// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := TaskStart{Function: slip.Function{Name: "flow-task-start", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-task-start",
			Args: []*slip.DocArg{
				{
					Name: "task",
					Type: "flow-task",
					Text: "to start.",
				},
			},
			Return: "nil",
			Text:   `__flow-task-start__ starts the _task_ if workers is greater than zero.`,
			Examples: []string{
				`(setq task (make-instance 'flow-task :name "tisk" :start 3))`,
				`(flow-task-start task) => nil`,
			},
		}, &Pkg)
}

// TaskStart represents the flow-task-start function.
type TaskStart struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *TaskStart) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.CheckArgCount(s, depth, f, args, 1, 1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Type != taskFlavor {
		slip.TypePanic(s, depth, "task", args[0], "task")
	}
	self.Any.(*task).start(s)

	return nil
}

type taskStartCaller struct{}

func (caller taskStartCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	obj.Any.(*task).start(s)

	return nil
}

func (caller taskStartCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":start", "flow-task-start", "flow-task", "task")
}
