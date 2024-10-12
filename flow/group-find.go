// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := GroupFind{Function: slip.Function{Name: "flow-group-find", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-group-find",
			Args: []*slip.DocArg{
				{
					Name: "group",
					Type: "flow-group",
					Text: "to find.",
				},
				{
					Name: "flow",
					Type: "string",
					Text: "name of flow to find.",
				},
			},
			Return: "instance",
			Text:   `__flow-group-find__ finds a _flow_ in the _group_.`,
			Examples: []string{
				`(setq group (make-instance 'flow-group-flavor))`,
				`(flow-group-add group (make-instance 'flow-flavor :name "flo")) => nil`,
				`(flow-group-find group 'flo) => #<flow-flavor 12345>`,
			},
		}, &Pkg)
}

// GroupFind represents the flow-group-find function.
type GroupFind struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *GroupFind) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.ArgCountCheck(f, args, 2, 2)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Type != groupFlavor {
		slip.PanicType("group", args[0], "group")
	}
	return self.Any.(*group).find(args[1])
}

type groupFindCaller struct{}

func (caller groupFindCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	return obj.Any.(*group).find(args[0])
}

func (caller groupFindCaller) Docs() string {
	return methodDocFromFunc(":find", "flow-group-find", "flow-group-flavor", "group")
}
