// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
	"github.com/ohler55/slip/pkg/gi"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := FlowSetExitChannel{Function: slip.Function{Name: "flow-set-exit-channel", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-set-exit-channel",
			Args: []*slip.DocArg{
				{
					Name: "flow",
					Type: "flow",
					Text: "to return the exit-channel of.",
				},
				{
					Name: "channel",
					Type: "channel",
					Text: "to set the exit-channel to.",
				},
			},
			Return: "channel",
			Text: `__flow-set-exit-channel__ returns the exit-channel of a _flow-flavor_ instance
after setting exit-channel.`,
			Examples: []string{
				`(setq flow (make-instance 'flow-flavor))`,
				`(flow-set-exit-channel flow (make-channel 2)) => #<channel 12345>`,
				`(flow-exit-channel flow) => #<channel 12345>`,
			},
		}, &Pkg)
}

// FlowSetExitChannel represents the flow-set-exit-channel function.
type FlowSetExitChannel struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *FlowSetExitChannel) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	slip.ArgCountCheck(f, args, 2, 2)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Flavor != flowFlavor {
		slip.PanicType("flow", args[0], "flow")
	}
	switch ta := args[1].(type) {
	case nil:
		self.Any.(*flow).exitChan = nil
	case gi.Channel:
		self.Any.(*flow).exitChan = ta
	default:
		slip.PanicType("channel", ta, "channel")
	}
	return self.Any.(*flow).exitChan
}

type flowSetExitChannelCaller struct{}

func (caller flowSetExitChannelCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)
	switch ta := args[0].(type) {
	case nil:
		obj.Any.(*flow).exitChan = nil
	case gi.Channel:
		obj.Any.(*flow).exitChan = ta
	default:
		slip.PanicType("channel", ta, "channel")
	}
	return obj.Any.(*flow).exitChan
}

func (caller flowSetExitChannelCaller) Docs() string {
	return methodDocFromFunc(":set-exit-channel", "flow-set-exit-channel", "flow-flavor", "flow")
}
