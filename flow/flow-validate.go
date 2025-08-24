// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := FlowValidate{Function: slip.Function{Name: "flow-validate", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-validate",
			Args: []*slip.DocArg{
				{
					Name: "flow",
					Type: "flow",
					Text: "to validate.",
				},
			},
			Return: "list",
			Text:   `__flow-validate__ returns a list of validation failure messages if any.`,
			Examples: []string{
				`(setq flow (make-instance 'flow :name "flo"))`,
				`(flow-validate flow) => ("no entry task")`,
			},
		}, &Pkg)
}

// FlowValidate represents the flow-validate function.
type FlowValidate struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *FlowValidate) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.CheckArgCount(s, depth, f, args, 1, 1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Type != flowFlavor {
		slip.TypePanic(s, depth, "flow", args[0], "flow")
	}
	return self.Any.(*flow).validate(s)
}

type flowValidateCaller struct{}

func (caller flowValidateCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	return obj.Any.(*flow).validate(s)
}

func (caller flowValidateCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":validate", "flow-validate", "flow", "flow")
}
