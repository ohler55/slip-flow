// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := FlowWrite{Function: slip.Function{Name: "flow-write", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-write",
			Args: []*slip.DocArg{
				{
					Name: "flow",
					Type: "flow",
					Text: "to return the write for.",
				},
				{Name: "&optional"},
				{
					Name: "destination",
					Type: "output-stream|string|t|nil",
					Text: `The destination to write to. If _t_ then write to _*standard-output*.
If _nil_ then return a string.`,
				},
				{
					Name: "clos",
					Type: "boolean",
					Text: `Use CLOS methods or functions to build the _flow_ in the code written.`,
				},
			},
			Return: "string|nil",
			Text: `__flow-write__ writes the _flow_ to the _destination_.
If _destination_ is _nil_ then a string is returned otherwise _nil_ is returned.`,
		}, &Pkg)
}

// FlowWrite represents the flow-write function.
type FlowWrite struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *FlowWrite) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.ArgCountCheck(f, args, 1, 3)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Type != flowFlavor {
		slip.TypePanic(s, depth, "flow", args[0], "flow")
	}
	return self.Any.(*flow).write(s, args[1:], depth)
}

type flowWriteCaller struct{}

func (caller flowWriteCaller) Call(s *slip.Scope, args slip.List, depth int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	return obj.Any.(*flow).write(s, args, depth)
}

func (caller flowWriteCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":write", "flow-write", "flow", "flow")
}
