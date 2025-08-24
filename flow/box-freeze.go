// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := BoxFreeze{Function: slip.Function{Name: "flow-box-freeze", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-box-freeze",
			Args: []*slip.DocArg{
				{
					Name: "box",
					Type: "flow-box",
					Text: "to freeze.",
				},
			},
			Return: "",
			Text:   `__flow-box-freeze__ freezes an instance of the _flow-box_.`,
			Examples: []string{
				`(setq box (make-instance 'flow-box :parse "{a:7}"))`,
				`(flow-box-freeze box)`,
				`(send box :frozen) => t`,
			},
		}, &Pkg)
}

// BoxFreeze represents the flow-box-freeze function.
type BoxFreeze struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *BoxFreeze) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	self, ok := args[0].(*flavors.Instance)
	if !ok {
		slip.TypePanic(s, depth, "box", args[0], "box")
	}
	return self.Receive(s, ":freeze", args[1:], depth)
}

type boxFreezeCaller struct{}

func (caller boxFreezeCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any.(*box).frozen = true

	return nil
}

func (caller boxFreezeCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":freeze", "flow-box-freeze", "flow-box", "box")
}
