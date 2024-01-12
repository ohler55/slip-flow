// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := BoxNative{Function: slip.Function{Name: "flow-box-native", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-box-native",
			Args: []*slip.DocArg{
				{
					Name: "box",
					Type: "flow-box",
					Text: "to convert to a native LISP s-expression.",
				},
			},
			Return: "object",
			Text: `__flow-box-native__ converts the content of an instance of the
_flow-box-flavor_ to a native LISP s-expression.`,
			Examples: []string{
				`(setq box (make-instance 'flow-box-flavor :parse "{a:7}"))`,
				`(flow-box-native box) => (("a" . 7))`,
			},
		}, &Pkg)
}

// BoxNative represents the flow-box-native function.
type BoxNative struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *BoxNative) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	self, ok := args[0].(*flavors.Instance)
	if !ok {
		slip.PanicType("box", args[0], "box")
	}
	return self.Receive(s, ":native", args[1:], depth)
}

type boxNativeCaller struct{}

func (caller boxNativeCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		flavors.PanicMethodArgChoice(obj, ":native", len(args), "0")
	}
	// fmt.Printf("*** any: %T\n", obj.Any)
	// fmt.Printf("*** content: %T\n", obj.Any.(*box).content)
	return slip.SimpleObject(obj.Any.(*box).content)
}

func (caller boxNativeCaller) Docs() string {
	return methodDocFromFunc(":native", "flow-box-native", "flow-box-flavor", "box")
}
