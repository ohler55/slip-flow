// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := TaskUpdateLink{Function: slip.Function{Name: "flow-task-update-link", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-task-update-link",
			Args: []*slip.DocArg{
				{
					Name: "task",
					Type: "flow-task",
					Text: "to receive a box.",
				},
				{
					Name: "link-name",
					Type: "string",
					Text: "name of the link to update.",
				},
				{
					Name: "mid-point",
					Type: "list",
					Text: "new mid-points for a link.",
				},
			},
			Return: "nil",
			Text:   `__flow-task-update-link__ updates the _mid-points_ of a link.`,
			Examples: []string{
				`(setq task (make-instance 'flow-task-flavor :name "tisk"))`,
				`(flow-task-update-link task "ok" '((10 20))) => nil`,
			},
		}, &Pkg)
}

// TaskUpdateLink represents the flow-task-update-link function.
type TaskUpdateLink struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *TaskUpdateLink) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.ArgCountCheck(f, args, 3, 3)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Flavor != taskFlavor {
		slip.PanicType("task", args[0], "task")
	}
	self.Any.(*task).updateLink(args[1:])

	return nil
}

type taskUpdateLinkCaller struct{}

func (caller taskUpdateLinkCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any.(*task).updateLink(args)

	return nil
}

func (caller taskUpdateLinkCaller) Docs() string {
	return methodDocFromFunc(":update-link", "flow-task-update-link", "flow-task-flavor", "task")
}
