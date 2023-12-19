// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := BoxFrozen{Function: slip.Function{Name: "flow-box-frozen", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-box-frozen",
			Args: []*slip.DocArg{
				{
					Name: "box",
					Type: "flow-box",
					Text: "to return then frozen state from.",
				},
			},
			Return: "boolean",
			Text:   `__flow-box-frozen__ returns the frozen state of an instance of the _flow-box-flavor_.`,
			Examples: []string{
				`(setq box (make-instance 'flow-box-flavor :parse "{a:7}"))`,
				`(flow-box-freeze box)`,
				`(flow-box-frozen box) => t`,
			},
		}, &Pkg)
}

// BoxFrozen represents the flow-box-frozen function.
type BoxFrozen struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *BoxFrozen) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	self, ok := args[0].(*flavors.Instance)
	if !ok {
		slip.PanicType("box", args[0], "box")
	}
	return self.Receive(s, ":frozen", args[1:], depth)
}

type boxFrozenCaller struct{}

func (caller boxFrozenCaller) Call(s *slip.Scope, args slip.List, depth int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)
	if obj.Any.(*box).frozen {
		value = slip.True
	}
	return
}

func (caller boxFrozenCaller) Docs() string {
	return methodDocFromFunc(":frozen", "flow-box-frozen", "flow-box-flavor", "box")
}
