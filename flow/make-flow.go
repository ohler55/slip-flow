// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := MakeFlow{Function: slip.Function{Name: "make-flow", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "make-flow",
			Args: []*slip.DocArg{
				{Name: "&key"},
				{
					Name: "name",
					Type: "string",
					Text: `Sets name of the flow`,
				},
				{
					Name: "exit-channel",
					Type: "channel",
					Text: "Sets the exit-channel of the flow.",
				},
				{
					Name: "logger",
					Type: "instance",
					Text: "Sets the logger for the flow. The instance must have the _:log_ method",
				},
			},
			Return: "flow",
			Text:   `__make-flow__ make a new instance of the _flow-flavor_.`,
			Examples: []string{
				`(make-flow :name "flo") => #<flow-flavor 12345>`,
			},
		}, &Pkg)
}

// MakeFlow represents the make-flow function.
type MakeFlow struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *MakeFlow) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	self := flowFlavor.MakeInstance().(*flavors.Instance)
	self.Init(s, args, depth)

	return self
}
