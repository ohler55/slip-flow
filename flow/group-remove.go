// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := GroupRemove{Function: slip.Function{Name: "flow-group-remove", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-group-remove",
			Args: []*slip.DocArg{
				{
					Name: "group",
					Type: "flow-group",
					Text: "to remove.",
				},
				{
					Name: "flow",
					Type: "string",
					Text: "name of flow to remove.",
				},
			},
			Return: "nil",
			Text:   `__flow-group-remove__ removes a _flow_ from the _group_.`,
			Examples: []string{
				`(setq group (make-instance 'flow-group))`,
				`(flow-group-add group (make-instance 'flow-flavor :name "flo")) => nil`,
				`(flow-group-remove group 'flo) => nil`,
			},
		}, &Pkg)
}

// GroupRemove represents the flow-group-remove function.
type GroupRemove struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *GroupRemove) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.CheckArgCount(s, depth, f, args, 2, 2)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Type != groupFlavor {
		slip.TypePanic(s, depth, "group", args[0], "group")
	}
	self.Any.(*group).remove(s, args[1], depth)

	return nil
}

type groupRemoveCaller struct{}

func (caller groupRemoveCaller) Call(s *slip.Scope, args slip.List, depth int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	obj.Any.(*group).remove(s, args[0], depth)

	return nil
}

func (caller groupRemoveCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":remove", "flow-group-remove", "flow-group", "group")
}
