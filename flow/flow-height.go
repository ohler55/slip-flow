// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := FlowHeight{Function: slip.Function{Name: "flow-height", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-height",
			Args: []*slip.DocArg{
				{
					Name: "flow",
					Type: "flow",
					Text: "to return the height of.",
				},
			},
			Return: "string",
			Text:   `__flow-height__ returns the height of a _flow_ instance.`,
			Examples: []string{
				`(setq flow (make-instance 'flow :name "flo" :height 300))`,
				`(flow-height flow) => 300`,
			},
		}, &Pkg)
}

// FlowHeight represents the flow-height function.
type FlowHeight struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *FlowHeight) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	slip.ArgCountCheck(f, args, 1, 1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Type != flowFlavor {
		slip.PanicType("flow", args[0], "flow")
	}
	return self.Any.(*flow).height()
}

type flowHeightCaller struct{}

func (caller flowHeightCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	return obj.Any.(*flow).height()
}

func (caller flowHeightCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":height", "flow-height", "flow", "flow")
}
