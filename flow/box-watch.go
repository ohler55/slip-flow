// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
	"github.com/ohler55/slip/pkg/gi"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := BoxWatch{Function: slip.Function{Name: "flow-box-watch", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-box-watch",
			Args: []*slip.DocArg{
				{
					Name: "box",
					Type: "flow-box",
					Text: "to watch a whatch channel to.",
				},
				{
					Name: "name",
					Type: "string",
					Text: `of the watcher to add.`,
				},
				{
					Name: "channel",
					Type: "channel",
					Text: `channel to push notifications to.`,
				},
			},
			Return: "nil",
			Text:   `__flow-box-watch__ adds a watcher _channel_ associated with the _name_.`,
			Examples: []string{
				`(setq box (make-instance 'flow-box :parse "{a:7}")) => #<flow-box 12345>`,
				`(setq chan (make-channel 5) => #<channel 12345>`,
				`(flow-box-watch box "done" chan) => nil`,
			},
		}, &Pkg)
}

// BoxWatch represents the flow-box-watch function.
type BoxWatch struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *BoxWatch) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	self, ok := args[0].(*flavors.Instance)
	if !ok {
		slip.TypePanic(s, depth, "box", args[0], "box")
	}
	_ = self.Receive(s, ":watch", args[1:], depth)

	return nil
}

type boxWatchCaller struct{}

func (caller boxWatchCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	if len(args) == 2 {
		watchBox(s, obj, args[0], args[1], depth)
	} else {
		slip.MethodArgChoicePanic(s, depth, obj, ":watch", len(args), "2")
	}
	return nil
}

func (caller boxWatchCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":watch", "flow-box-watch", "flow-box", "box")
}

func watchBox(s *slip.Scope, obj *flavors.Instance, name, channel slip.Object, depth int) {
	var key string

	switch tn := name.(type) {
	case slip.String:
		key = string(tn)
	case slip.Symbol:
		key = string(tn)
	default:
		slip.TypePanic(s, depth, "name", tn, "string", "symbol")
	}
	if c, ok := channel.(gi.Channel); ok {
		obj.Any.(*box).watchers[key] = c
	} else {
		slip.TypePanic(s, depth, "channel", channel, "channel")
	}
}
