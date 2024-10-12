// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := GroupShutdown{Function: slip.Function{Name: "flow-group-shutdown", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-group-shutdown",
			Args: []*slip.DocArg{
				{
					Name: "group",
					Type: "flow-group",
					Text: "to shutdown.",
				},
			},
			Return: "nil",
			Text:   `__flow-group-shutdown__ shutdown all the _flows_ in the _group_.`,
			Examples: []string{
				`(setq group (make-instance 'flow-group-flavor))`,
				`(flow-group-add group (make-instance 'flow-flavor)) => nil`,
				`(flow-group-shutdown group) => nil`,
			},
		}, &Pkg)
}

// GroupShutdown represents the flow-group-shutdown function.
type GroupShutdown struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *GroupShutdown) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.ArgCountCheck(f, args, 1, 1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Type != groupFlavor {
		slip.PanicType("group", args[0], "group")
	}
	self.Any.(*group).shutdown(s)

	return nil
}

type groupShutdownCaller struct{}

func (caller groupShutdownCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	obj.Any.(*group).shutdown(s)

	return nil
}

func (caller groupShutdownCaller) Docs() string {
	return methodDocFromFunc(":shutdown", "flow-group-shutdown", "flow-group-flavor", "group")
}
