// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := GroupAdd{Function: slip.Function{Name: "flow-group-add", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-group-add",
			Args: []*slip.DocArg{
				{
					Name: "group",
					Type: "flow-group",
					Text: "to add.",
				},
				{
					Name: "flow",
					Type: "flow",
					Text: "to add.",
				},
			},
			Return: "nil",
			Text:   `__flow-group-add__ adds a _flow_ to the _group_.`,
			Examples: []string{
				`(setq group (make-instance 'flow-group-flavor))`,
				`(flow-group-add group (make-instance 'flow-flavor :name "flo")) => nil`,
				`(flow-group-add group) => nil`,
			},
		}, &Pkg)
}

// GroupAdd represents the flow-group-add function.
type GroupAdd struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *GroupAdd) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.ArgCountCheck(f, args, 2, 2)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Flavor != groupFlavor {
		slip.PanicType("group", args[0], "group")
	}
	self.Any.(*group).add(args[1])

	return nil
}

type groupAddCaller struct{}

func (caller groupAddCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	obj.Any.(*group).add(args[0])

	return nil
}

func (caller groupAddCaller) Docs() string {
	return methodDocFromFunc(":add", "flow-group-add", "flow-group-flavor", "group")
}
