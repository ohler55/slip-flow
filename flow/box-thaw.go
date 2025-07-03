// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := BoxThaw{Function: slip.Function{Name: "flow-box-thaw", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-box-thaw",
			Args: []*slip.DocArg{
				{
					Name: "box",
					Type: "flow-box",
					Text: "to thaw.",
				},
			},
			Return: "",
			Text:   `__flow-box-thaw__ thaws an instance of the _flow-box_.`,
			Examples: []string{
				`(setq box (make-instance 'flow-box :parse "{a:7}"))`,
				`(flow-box-thaw box)`,
				`(send box :frozen) => nil`,
			},
		}, &Pkg)
}

// BoxThaw represents the flow-box-thaw function.
type BoxThaw struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *BoxThaw) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	self, ok := args[0].(*flavors.Instance)
	if !ok {
		slip.PanicType("box", args[0], "box")
	}
	return self.Receive(s, ":thaw", args[1:], depth)
}

type boxThawCaller struct{}

func (caller boxThawCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any.(*box).frozen = false

	return nil
}

func (caller boxThawCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":thaw", "flow-box-thaw", "flow-box", "box")
}
