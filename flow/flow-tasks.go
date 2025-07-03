// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := FlowTasks{Function: slip.Function{Name: "flow-tasks", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-tasks",
			Args: []*slip.DocArg{
				{
					Name: "flow",
					Type: "flow",
					Text: "to return a list of tasks from.",
				},
			},
			Return: "list",
			Text:   `__flow-tasks__ returns a list of all the _tasks_ in the _flow_.`,
			Examples: []string{
				`(setq flow (make-instance 'flow :tasks "flo"))`,
				`(flow-add-task flow :name 'tisk :actor (lambda (b) (list 'ok b)))`,
				`(flow-tasks flow) => (#<flow-task-flavor 12345>)`,
			},
		}, &Pkg)
}

// FlowTasks represents the flow-tasks function.
type FlowTasks struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *FlowTasks) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.ArgCountCheck(f, args, 1, 1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Type != flowFlavor {
		slip.PanicType("flow", args[0], "flow")
	}
	return self.Any.(*flow).taskList()
}

type flowTasksCaller struct{}

func (caller flowTasksCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)

	return obj.Any.(*flow).taskList()
}

func (caller flowTasksCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":tasks", "flow-tasks", "flow", "flow")
}
