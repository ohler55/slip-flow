// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := MakeFlowBox{Function: slip.Function{Name: "make-flow-box", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "make-flow-box",
			Args: []*slip.DocArg{
				{Name: "&key"},
				{
					Name: "tracking-id",
					Type: "string|fixnum|uuid",
					Text: `Sets the tracking id of the box to the provided value which
can be a string, fixnum, or gi:uuid.`,
				},
				{
					Name: "track",
					Type: "flow-track instance",
					Text: "Sets the tracking id and events of the box to the provided values.",
				},
				{
					Name: "set",
					Type: "bag-flavor instance|JSON compatible LISP",
					Text: "Sets the contents with the LISP or _bag-flavor_ instance.",
				},
				{
					Name: "parse",
					Type: "string",
					Text: "Sets the contents with the LISP or _bag-flavor_ instance.",
				},
				{
					Name: "parse",
					Type: "bag-flavor instance|JSON compatible LISP",
					Text: "A JSON or SEN string to form the content of the box.",
				},
				{
					Name: "read",
					Type: "bag-flavor instance|JSON compatible LISP",
					Text: "Read from an _input-stream_ and parses read JSON or SEN to form the content.",
				},
				{
					Name: "watch",
					Type: "symbol bound to a gi:channel",
					Text: "Add a watcher to the box with the name of the symbol which must be bound to a gi:channel",
				},
			},
			Return: "box",
			Text: `__make-flow-box__ make a new instance of the _flow-box-flavor_.

If no _path_ is provided the entire contents of the box is replaced.

This is the same as the calling _(make-instance 'flow-box-flavor)_ with the same
arguments as initializers.`,
			Examples: []string{
				`(make-flow-box :parse "{a:7}") => #<flow-box-flavor 12345> ;; content is {a:7}`,
			},
		}, &Pkg)
}

// MakeFlowBox represents the make-flow-box function.
type MakeFlowBox struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *MakeFlowBox) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	self := boxFlavor.MakeInstance().(*flavors.Instance)
	_ = self.Receive(s, ":init", slip.List{args}, depth)

	return self
}
