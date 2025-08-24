// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := FlowName{Function: slip.Function{Name: "flow-name", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-name",
			Args: []*slip.DocArg{
				{
					Name: "flow",
					Type: "flow",
					Text: "to return the name of.",
				},
			},
			Return: "string",
			Text:   `__flow-name__ returns the name of a _flow_ instance.`,
			Examples: []string{
				`(setq flow (make-instance 'flow :name "flo"))`,
				`(flow-name flow) => "flo"`,
			},
		}, &Pkg)
}

// FlowName represents the flow-name function.
type FlowName struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *FlowName) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	slip.CheckArgCount(s, depth, f, args, 1, 1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Type != flowFlavor {
		slip.TypePanic(s, depth, "flow", args[0], "flow")
	}
	return slip.String(self.Any.(*flow).name)
}

type flowNameCaller struct{}

func (caller flowNameCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	return slip.String(obj.Any.(*flow).name)
}

func (caller flowNameCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":name", "flow-name", "flow", "flow")
}
