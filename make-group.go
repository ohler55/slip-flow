// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := MakeGroup{Function: slip.Function{Name: "make-flow-group", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "make-flow-group",
			Args: []*slip.DocArg{
				{Name: "&key"},
				{
					Name: "logger",
					Type: "instance",
					Text: "Sets the logger for the group. The instance must have the _:log_ method",
				},
			},
			Return: "group",
			Text:   `__make-flow-group__ make a new instance of the _group-flavor_.`,
			Examples: []string{
				`(make-flow-group) => #<flow-group-flavor 12345>`,
			},
		}, &Pkg)
}

// MakeGroup represents the make-group function.
type MakeGroup struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *MakeGroup) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	self := groupFlavor.MakeInstance().(*flavors.Instance)
	self.Init(s, args, depth)

	return self
}
