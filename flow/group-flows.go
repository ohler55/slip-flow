// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := GroupFlows{Function: slip.Function{Name: "flow-group-flows", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-group-flows",
			Args: []*slip.DocArg{
				{
					Name: "group",
					Type: "flow-group",
					Text: "to flows.",
				},
			},
			Return: "list",
			Text:   `__flow-group-flows__ returns a list of all _flows_ in the _group_.`,
			Examples: []string{
				`(setq group (make-instance 'flow-group-flavor))`,
				`(flow-group-add group (make-instance 'flow-flavor :name "flo")) => nil`,
				`(flow-group-flows group) => (#<flow-flavor 12345>)`,
			},
		}, &Pkg)
}

// GroupFlows represents the flow-group-flows function.
type GroupFlows struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *GroupFlows) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.ArgCountCheck(f, args, 1, 1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Flavor != groupFlavor {
		slip.PanicType("group", args[0], "group")
	}
	return self.Any.(*group).flowList()
}

type groupFlowsCaller struct{}

func (caller groupFlowsCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	return obj.Any.(*group).flowList()
}

func (caller groupFlowsCaller) Docs() string {
	return methodDocFromFunc(":flows", "flow-group-flows", "flow-group-flavor", "group")
}
