// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := GroupStart{Function: slip.Function{Name: "flow-group-start", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-group-start",
			Args: []*slip.DocArg{
				{
					Name: "group",
					Type: "flow-group",
					Text: "to start.",
				},
			},
			Return: "nil",
			Text:   `__flow-group-start__ starts all the _flows_ in the _group_.`,
			Examples: []string{
				`(setq group (make-instance 'flow-group-flavor))`,
				`(flow-group-add group (make-instance 'flow-flavor)) => nil`,
				`(flow-group-start group) => nil`,
			},
		}, &Pkg)
}

// GroupStart represents the flow-group-start function.
type GroupStart struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *GroupStart) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.ArgCountCheck(f, args, 1, 1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Flavor != groupFlavor {
		slip.PanicType("group", args[0], "group")
	}
	self.Any.(*group).start(s)

	return nil
}

type groupStartCaller struct{}

func (caller groupStartCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	obj.Any.(*group).start(s)

	return nil
}

func (caller groupStartCaller) Docs() string {
	return methodDocFromFunc(":start", "flow-group-start", "flow-group-flavor", "group")
}
