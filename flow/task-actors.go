// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := TaskActors{Function: slip.Function{Name: "flow-task-actors", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-task-actors",
			Args: []*slip.DocArg{
				{
					Name: "task",
					Type: "flow-task",
					Text: "to return the actors of.",
				},
			},
			Return: "string",
			Text:   `__flow-task-actors__ returns the actors of a _flow-task_ instance.`,
			Examples: []string{
				`(setq task (make-instance 'flow-task :actors "tisk"))`,
				`(flow-task-actors task) => ("tisk")`,
			},
		}, &Pkg)
}

// TaskActors represents the flow-task-actors function.
type TaskActors struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *TaskActors) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	slip.ArgCountCheck(f, args, 1, 1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Type != taskFlavor {
		slip.TypePanic(s, depth, "task", args[0], "task")
	}
	return self.Any.(*task).actorList()
}

type taskActorsCaller struct{}

func (caller taskActorsCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	return obj.Any.(*task).actorList()
}

func (caller taskActorsCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":actors", "flow-task-actors", "flow-task", "task")
}
