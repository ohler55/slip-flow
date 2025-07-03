// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := FlowWidth{Function: slip.Function{Name: "flow-width", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-width",
			Args: []*slip.DocArg{
				{
					Name: "flow",
					Type: "flow",
					Text: "to return the width of.",
				},
			},
			Return: "string",
			Text:   `__flow-width__ returns the width of a _flow_ instance.`,
			Examples: []string{
				`(setq flow (make-instance 'flow :name "flo" :width 300))`,
				`(flow-width flow) => 300`,
			},
		}, &Pkg)
}

// FlowWidth represents the flow-width function.
type FlowWidth struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *FlowWidth) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	slip.ArgCountCheck(f, args, 1, 1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Type != flowFlavor {
		slip.PanicType("flow", args[0], "flow")
	}
	return self.Any.(*flow).width()
}

type flowWidthCaller struct{}

func (caller flowWidthCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	return obj.Any.(*flow).width()
}

func (caller flowWidthCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":width", "flow-width", "flow", "flow")
}
