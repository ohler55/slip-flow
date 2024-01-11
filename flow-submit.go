// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := FlowSubmit{Function: slip.Function{Name: "flow-submit", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-submit",
			Args: []*slip.DocArg{
				{
					Name: "flow",
					Type: "flow",
					Text: "to submit a box to.",
				},
				{
					Name: "box",
					Type: "instance",
					Text: "instance of the flow-box-flavor to process by the flow.",
				},
			},
			Return: "nil",
			Text:   `__flow-submit__ submits an instance of the _flow-box-flavor_ for processing by the _flow_.`,
			Examples: []string{
				`(setq flow (make-instance 'flow-flavor :submit "flo"))`,
				`(flow-add-task flow :name 'tisk :actor (lambda (b) (list 'ok b)))`,
				`(flow-submit flow "tisk") => nil`,
			},
		}, &Pkg)
}

// FlowSubmit represents the flow-submit function.
type FlowSubmit struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *FlowSubmit) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.ArgCountCheck(f, args, 2, 2)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Flavor != flowFlavor {
		slip.PanicType("flow", args[0], "flow")
	}
	self.Any.(*flow).submit(s, args[1])

	return nil
}

type flowSubmitCaller struct{}

func (caller flowSubmitCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)

	obj.Any.(*flow).submit(s, args[0])

	return nil
}

func (caller flowSubmitCaller) Docs() string {
	return methodDocFromFunc(":submit", "flow-submit", "flow-flavor", "flow")
}
