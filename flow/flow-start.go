// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := FlowStart{Function: slip.Function{Name: "flow-start", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-start",
			Args: []*slip.DocArg{
				{
					Name: "flow",
					Type: "flow",
					Text: "to start.",
				},
			},
			Return: "nil",
			Text:   `__flow-start__ starts all the _tasks_ in the _flow_.`,
			Examples: []string{
				`(setq flow (make-instance 'flow-flavor :start "flo"))`,
				`(flow-add-task flow :name 'tisk :actor (lambda (b) (list 'ok b)))`,
				`(flow-start flow) => nil`,
			},
		}, &Pkg)
}

// FlowStart represents the flow-start function.
type FlowStart struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *FlowStart) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.ArgCountCheck(f, args, 1, 1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Flavor != flowFlavor {
		slip.PanicType("flow", args[0], "flow")
	}
	self.Any.(*flow).start(s)

	return nil
}

type flowStartCaller struct{}

func (caller flowStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)

	obj.Any.(*flow).start(s)

	return nil
}

func (caller flowStartCaller) Docs() string {
	return methodDocFromFunc(":start", "flow-start", "flow-flavor", "flow")
}
