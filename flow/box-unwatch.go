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
					Text: "to unwatch a named channel of.",
				},
				{Name: "&optional"},
				{
					Name: "name",
					Type: "string",
					Text: `of the watcher to remove.`,
				},
			},
			Return: "nil",
			Text: `__flow-box-unwatch__ removed a watcher _channel_ associated with the _name_
of if not _name_ is provided all watchers are removed.`,
			Examples: []string{
				`(setq box (make-instance 'flow-box :parse "{a:7}")) => #<flow-box 12345>`,
				`(setq chan (make-channel 5) => #<channel 12345>`,
				`(flow-box-watch box "done" chan) => nil`,
				`(flow-box-unwatch box "done") => nil`,
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
		slip.TypePanic(s, depth, "box", args[0], "box")
	}
	_ = self.Receive(s, ":unwatch", args[1:], depth)

	return nil
}

type boxUnwatchCaller struct{}

func (caller boxUnwatchCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	switch len(args) {
	case 0:
		unwatchBox(s, obj, nil, depth)
	case 1:
		unwatchBox(s, obj, args[0], depth)
	default:
		slip.PanicMethodArgChoice(obj, ":unwatch", len(args), "0 or 1")
	}
	return nil
}

func (caller boxUnwatchCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":unwatch", "flow-box-unwatch", "flow-box", "box")
}

func unwatchBox(s *slip.Scope, obj *flavors.Instance, name slip.Object, depth int) {
	if name == nil {
		obj.Any.(*box).watchers = map[string]gi.Channel{}
	} else {
		var key string
		switch tn := name.(type) {
		case slip.String:
			key = string(tn)
		case slip.Symbol:
			key = string(tn)
		default:
			slip.TypePanic(s, depth, "name", tn, "string", "symbol")
		}
		delete(obj.Any.(*box).watchers, key)
	}
}
