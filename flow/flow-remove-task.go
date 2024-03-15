// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := FlowRemoveTask{Function: slip.Function{Name: "flow-remove-task", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-remove-task",
			Args: []*slip.DocArg{
				{
					Name: "flow",
					Type: "flow",
					Text: "to remove a task from.",
				},
				{
					Name: "task-name",
					Type: "string",
					Text: "name of the task to remove.",
				},
			},
			Return: "nil",
			Text:   `__flow-remove-task__ remove a task from the tasks in the _flow_.`,
			Examples: []string{
				`(setq flow (make-instance 'flow-flavor :remove-task "flo"))`,
				`(flow-add-task flow :name 'tisk :actor (lambda (b) (list 'ok b)))`,
				`(flow-remove-task flow "tisk") => nil`,
				`(flow-tasks flow) => ()`,
			},
		}, &Pkg)
}

// FlowRemoveTask represents the flow-remove-task function.
type FlowRemoveTask struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *FlowRemoveTask) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.ArgCountCheck(f, args, 2, 2)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Flavor != flowFlavor {
		slip.PanicType("flow", args[0], "flow")
	}
	self.Any.(*flow).removeTask(args[1])

	return nil
}

type flowRemoveTaskCaller struct{}

func (caller flowRemoveTaskCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)

	obj.Any.(*flow).removeTask(args[0])

	return nil
}

func (caller flowRemoveTaskCaller) Docs() string {
	return methodDocFromFunc(":remove-task", "flow-remove-task", "flow-flavor", "flow")
}
