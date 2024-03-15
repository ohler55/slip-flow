// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := TaskUnlink{Function: slip.Function{Name: "flow-task-unlink", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-task-unlink",
			Args: []*slip.DocArg{
				{
					Name: "task",
					Type: "flow-task",
					Text: "to unlink a box.",
				},
				{
					Name: "link",
					Type: "string",
					Text: "to unlink.",
				},
			},
			Return: "nil",
			Text:   `__flow-task-unlink__ removes a link from the _task_.`,
			Examples: []string{
				`(setq flow (make-instance 'flow-flavor :submit "flo"))`,
				`(setq tick (flow-add-task flow :name 'tick :actor (lambda (b) (list 'ok b))))`,
				`(flow-add-task flow :name 'tock :actor (lambda (b) (list 'ok b)))`,
				`(flow-link flow 'ok 'tick 'tock) => nil`,
				`(flow-task-unlink tick 'ok) => nil`,
			},
		}, &Pkg)
}

// TaskUnlink represents the flow-task-unlink function.
type TaskUnlink struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *TaskUnlink) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.ArgCountCheck(f, args, 2, 2)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Flavor != taskFlavor {
		slip.PanicType("task", args[0], "task")
	}
	self.Any.(*task).unlink(args[1:])

	return nil
}

type taskUnlinkCaller struct{}

func (caller taskUnlinkCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	obj.Any.(*task).unlink(args)

	return nil
}

func (caller taskUnlinkCaller) Docs() string {
	return methodDocFromFunc(":unlink", "flow-task-unlink", "flow-task-flavor", "task")
}
