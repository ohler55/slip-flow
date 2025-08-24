// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := TaskDepth{Function: slip.Function{Name: "flow-task-depth", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-task-depth",
			Args: []*slip.DocArg{
				{
					Name: "task",
					Type: "flow-task",
					Text: "to return the number of depth from.",
				},
			},
			Return: "fixnum",
			Text:   `__flow-task-depth__ returns the number of depth for a _flow-task_ instance.`,
			Examples: []string{
				`(setq task (make-instance 'flow-task :name "tisk" :depth 3))`,
				`(flow-task-depth task) => 3`,
			},
		}, &Pkg)
}

// TaskDepth represents the flow-task-depth function.
type TaskDepth struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *TaskDepth) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	slip.CheckArgCount(s, depth, f, args, 1, 1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Type != taskFlavor {
		slip.TypePanic(s, depth, "task", args[0], "task")
	}
	return slip.Fixnum(self.Any.(*task).depth)
}

type taskDepthCaller struct{}

func (caller taskDepthCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	return slip.Fixnum(obj.Any.(*task).depth)
}

func (caller taskDepthCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":depth", "flow-task-depth", "flow-task", "task")
}
