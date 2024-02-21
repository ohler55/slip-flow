// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := BoxUnwatch{Function: slip.Function{Name: "flow-box-unwatch", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-box-unwatch",
			Args: []*slip.DocArg{
				{
					Name: "box",
					Type: "flow-box",
					Text: "to unwatch a whatch channel to.",
				},
				{
					Name: "name",
					Type: "string",
					Text: `of the unwatcher to add.`,
				},
				{
					Name: "channel",
					Type: "channel",
					Text: `channel to push notifications to.`,
				},
			},
			Return: "nil",
			Text:   `__flow-box-unwatch__ adds a unwatcher _channel_ associated with the _name_.`,
			Examples: []string{
				`(setq box (make-instance 'flow-box-flavor :parse "{a:7}")) => #<flow-box-flavor 12345>`,
				`(setq chan (make-channel 5) => #<channel 12345>`,
				`(flow-box-unwatch box "done" chan) => nil`,
			},
		}, &Pkg)
}

// BoxUnwatch represents the flow-box-unwatch function.
type BoxUnwatch struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *BoxUnwatch) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	self, ok := args[0].(*flavors.Instance)
	if !ok {
		slip.PanicType("box", args[0], "box")
	}
	_ = self.Receive(s, ":unwatch", args[1:], depth)

	return nil
}

type boxUnwatchCaller struct{}

func (caller boxUnwatchCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	if len(args) == 1 {
		unwatchBox(obj, args[0])
	} else {
		flavors.PanicMethodArgChoice(obj, ":unwatch", len(args), "1")
	}
	return nil
}

func (caller boxUnwatchCaller) Docs() string {
	return methodDocFromFunc(":unwatch", "flow-box-unwatch", "flow-box-flavor", "box")
}

func unwatchBox(obj *flavors.Instance, name slip.Object) {
	var key string
	switch tn := name.(type) {
	case slip.String:
		key = string(tn)
	case slip.Symbol:
		key = string(tn)
	default:
		slip.PanicType("name", tn, "string", "symbol")
	}
	delete(obj.Any.(*box).watchers, key)
}
