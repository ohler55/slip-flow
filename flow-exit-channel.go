// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
	"github.com/ohler55/slip/pkg/gi"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := FlowExitChannel{Function: slip.Function{Name: "flow-exit-channel", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-exit-channel",
			Args: []*slip.DocArg{
				{
					Name: "flow",
					Type: "flow",
					Text: "to return the exit-channel of.",
				},
			},
			Return: "string",
			Text:   `__flow-exit-channel__ returns the exit-channel of a _flow-flavor_ instance.`,
			Examples: []string{
				`(setq flow (make-instance 'flow-flavor :exit-channel "flo"))`,
				`(flow-exit-channel flow) => "flo"`,
			},
		}, &Pkg)
}

// FlowExitChannel represents the flow-flow-exit-channel function.
type FlowExitChannel struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *FlowExitChannel) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	slip.ArgCountCheck(f, args, 1, 1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Flavor != flowFlavor {
		slip.PanicType("flow", args[0], "flow")
	}
	return gi.Channel(self.Any.(*flow).exitChan)
}

type flowExitChannelCaller struct{}

func (caller flowExitChannelCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	return gi.Channel(obj.Any.(*flow).exitChan)
}

func (caller flowExitChannelCaller) Docs() string {
	return methodDocFromFunc(":exit-channel", "flow-exit-channel", "flow-flavor", "flow")
}
