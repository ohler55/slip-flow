// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main

import (
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/bag"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := BoxSet{Function: slip.Function{Name: "flow-box-set", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-box-set",
			Args: []*slip.DocArg{
				{
					Name: "box",
					Type: "flow-box",
					Text: "The _box_ to set a value in.",
				},
				{
					Name: "value",
					Type: "object",
					Text: "The _value_ to set in _box_ according to the path.",
				},
				{Name: "&optional"},
				{
					Name: "path",
					Type: "string|bag-path",
					Text: `The path to the location in the box to set the _value_.
The path must follow the JSONPath format.`,
				},
			},
			Return: "box",
			Text: `__flow-box-set__ sets a _value_ at the location described by _path_.
If no _path_ is provided the entire contents of the box is replaced.

This is the same as the _:set_ method of the _flow-box-flavor_ except none of the method's
daemons are invoked hence it has a slight performance advantage.`,
			Examples: []string{
				`(setq box (make-instance 'flow-box-flavor :parse "{a:7}"))`,
				`(flow-box-set box 3 "a") => #<flow-box-flavor 12345> ;; content is now {a:3}`,
			},
		}, &Pkg)
}

// BoxSet represents the flow-box-set function.
type BoxSet struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *BoxSet) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	if len(args) < 2 || 3 < len(args) {
		slip.PanicArgCount(f, 2, 3)
	}
	obj, ok := args[0].(*flavors.Instance)
	if !ok || obj.Flavor != boxFlavor {
		slip.PanicType("box", args[0], "flow-box")
	}
	if 2 < len(args) {
		setBox(obj, args[1], args[2])
	} else {
		setBox(obj, args[1], nil)
	}
	return obj
}

func setBox(obj *flavors.Instance, value, path slip.Object) {
	var x jp.Expr
	switch p := path.(type) {
	case nil:
	case slip.String:
		x = jp.MustParseString(string(p))
	case bag.Path:
		x = jp.Expr(p)
	default:
		slip.PanicType("path", p, "string")
	}
	v := bag.ObjectToBag(value)
	if x == nil {
		obj.Any.(*box).content = v
	} else {
		x.MustSet(obj.Any.(*box).content, v)
	}
}
