// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := TaskRunning{Function: slip.Function{Name: "flow-task-running", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-task-running",
			Args: []*slip.DocArg{
				{
					Name: "task",
					Type: "flow-task",
					Text: "to running.",
				},
			},
			Return: "nil",
			Text:   `__flow-task-running__ runnings the _task_ if workers is greater than zero.`,
			Examples: []string{
				`(setq task (make-instance 'flow-task :name "tisk" :running 3))`,
				`(flow-task-running task) => nil`,
			},
		}, &Pkg)
}

// TaskRunning represents the flow-task-running function.
type TaskRunning struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *TaskRunning) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.CheckArgCount(s, depth, f, args, 1, 1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Type != taskFlavor {
		slip.TypePanic(s, depth, "task", args[0], "task")
	}
	if self.Any.(*task).running() {
		return slip.True
	}
	return nil
}

type taskRunningCaller struct{}

func (caller taskRunningCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	if obj.Any.(*task).running() {
		return slip.True
	}
	return nil
}

func (caller taskRunningCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":running", "flow-task-running", "flow-task", "task")
}
