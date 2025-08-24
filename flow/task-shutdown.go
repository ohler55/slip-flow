// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := TaskShutdown{Function: slip.Function{Name: "flow-task-shutdown", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-task-shutdown",
			Args: []*slip.DocArg{
				{
					Name: "task",
					Type: "flow-task",
					Text: "to shutdown.",
				},
			},
			Return: "nil",
			Text:   `__flow-task-shutdown__ shutdowns the _task_ if workers is greater than zero.`,
			Examples: []string{
				`(setq task (make-instance 'flow-task :name "tisk" :shutdown 3))`,
				`(flow-task-start task) => nil`,
				`(flow-task-shutdown task) => nil`,
			},
		}, &Pkg)
}

// TaskShutdown represents the flow-task-shutdown function.
type TaskShutdown struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *TaskShutdown) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.ArgCountCheck(f, args, 1, 1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Type != taskFlavor {
		slip.TypePanic(s, depth, "task", args[0], "task")
	}
	self.Any.(*task).shutdown(s)

	return nil
}

type taskShutdownCaller struct{}

func (caller taskShutdownCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	obj.Any.(*task).shutdown(s)

	return nil
}

func (caller taskShutdownCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":shutdown", "flow-task-shutdown", "flow-task", "task")
}
