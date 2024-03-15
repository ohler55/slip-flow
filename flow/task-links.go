// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := TaskLinks{Function: slip.Function{Name: "flow-task-links", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-task-links",
			Args: []*slip.DocArg{
				{
					Name: "task",
					Type: "flow-task",
					Text: "to return the links for.",
				},
			},
			Return: "nil",
			Text:   `__flow-task-links__ returns the links the _task_.`,
			Examples: []string{
				`(setq flow (make-instance 'flow-flavor :link "flo"))`,
				`(setq tick (flow-add-task flow :name 'tick :actor (lambda (b) (list 'ok b))))`,
				`(setq tock (flow-add-task flow :name 'tock :actor (lambda (b) (list 'ok b))))`,
				`(flow-link flow "ok" 'tick 'tock) => nil`,
				`(flow-task-links tick) => (("tock" . #<flow-task-flavor 12345>))`,
			},
		}, &Pkg)
}

// TaskLinks represents the flow-task-links function.
type TaskLinks struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *TaskLinks) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.ArgCountCheck(f, args, 1, 1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Flavor != taskFlavor {
		slip.PanicType("task", args[0], "task")
	}
	return self.Any.(*task).linkList()
}

type taskLinksCaller struct{}

func (caller taskLinksCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	return obj.Any.(*task).linkList()
}

func (caller taskLinksCaller) Docs() string {
	return methodDocFromFunc(":links", "flow-task-links", "flow-task-flavor", "task")
}
