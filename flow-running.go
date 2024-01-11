// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := FlowRunning{Function: slip.Function{Name: "flow-running", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-running",
			Args: []*slip.DocArg{
				{
					Name: "flow",
					Type: "flow",
					Text: "to running.",
				},
			},
			Return: "boolean",
			Text:   `__flow-running__ runnings all the _tasks_ in the _flow_.`,
			Examples: []string{
				`(setq flow (make-instance 'flow-flavor :running "flo"))`,
				`(flow-add-task flow :name 'tisk :actor (lambda (b) (list 'ok b)))`,
				`(flow-start flow) => nil`,
				`(flow-running flow) => nil`,
			},
		}, &Pkg)
}

// FlowRunning represents the flow-running function.
type FlowRunning struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *FlowRunning) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.ArgCountCheck(f, args, 1, 1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Flavor != flowFlavor {
		slip.PanicType("flow", args[0], "flow")
	}
	if self.Any.(*flow).running() {
		return slip.True
	}
	return nil
}

type flowRunningCaller struct{}

func (caller flowRunningCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)

	if obj.Any.(*flow).running() {
		return slip.True
	}
	return nil
}

func (caller flowRunningCaller) Docs() string {
	return methodDocFromFunc(":running", "flow-running", "flow-flavor", "flow")
}
