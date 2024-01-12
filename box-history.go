// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := BoxHistory{Function: slip.Function{Name: "flow-box-history", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-box-history",
			Args: []*slip.DocArg{
				{
					Name: "box",
					Type: "flow-box",
					Text: "to get the history from.",
				},
			},
			Return: "list",
			Text: `__flow-box-history__ returns the history of the box track as a
list of triples where each triple is a list of the time, the task name, and the flow name.`,
			Examples: []string{
				`(setq box (make-instance 'flow-box-flavor :tracking-id 123))`,
				`(send box :scan "flo" "tisk")`,
				`(flow-box-history box) => ((@2023-12-15T19:23:17Z "tisk" "flo"))`,
			},
		}, &Pkg)
}

// BoxHistory represents the flow-box-history function.
type BoxHistory struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *BoxHistory) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	self, ok := args[0].(*flavors.Instance)
	if !ok {
		slip.PanicType("box", args[0], "box")
	}
	return self.Receive(s, ":history", args[1:], depth)
}

type boxHistoryCaller struct{}

func (caller boxHistoryCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)

	return obj.Any.(*box).track.historyList()
}

func (caller boxHistoryCaller) Docs() string {
	return methodDocFromFunc(":history", "flow-box-history", "flow-box-flavor", "box")
}
