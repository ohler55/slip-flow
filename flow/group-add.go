// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

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
				`(setq group (make-instance 'flow-group))`,
				`(flow-group-add group (make-instance 'flow :name "flo")) => nil`,
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
	slip.CheckArgCount(s, depth, f, args, 2, 2)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Type != groupFlavor {
		slip.TypePanic(s, depth, "group", args[0], "group")
	}
	self.Any.(*group).add(s, args[1], depth)

	return nil
}

type groupAddCaller struct{}

func (caller groupAddCaller) Call(s *slip.Scope, args slip.List, depth int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	obj.Any.(*group).add(s, args[0], depth)

	return nil
}

func (caller groupAddCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":add", "flow-group-add", "flow-group", "group")
}
