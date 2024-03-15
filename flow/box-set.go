// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/ojg/alt"
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
					Text: "to set a value in.",
				},
				{
					Name: "value",
					Type: "object",
					Text: "to set in _box_ according to the path.",
				},
				{Name: "&optional"},
				{
					Name: "path",
					Type: "string|bag-path",
					Text: `to the location in the box to set the _value_.
The path must follow the JSONPath format.`,
				},
			},
			Return: "box",
			Text: `__flow-box-set__ sets a _value_ at the location described by _path_.
If no _path_ is provided the entire contents of the box is replaced.`,
			Examples: []string{
				`(setq box (make-instance 'flow-box-flavor :parse "{a:7}")) => #<flow-box-flavor 12345>`,
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
	self, ok := args[0].(*flavors.Instance)
	if !ok {
		slip.PanicType("box", args[0], "box")
	}
	_ = self.Receive(s, ":set", args[1:], depth)

	return self
}

type boxSetCaller struct{}

func (caller boxSetCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	switch len(args) {
	case 1:
		setBox(obj, args[0], nil)
	case 2:
		setBox(obj, args[0], args[1])
	default:
		flavors.PanicMethodArgChoice(obj, ":set", len(args), "1 or 2")
	}
	return obj
}

func (caller boxSetCaller) Docs() string {
	return methodDocFromFunc(":set", "flow-box-set", "flow-box-flavor", "box")
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
	bx := obj.Any.(*box)
	if bx.frozen {
		bx.content = alt.Dup(bx.content)
		bx.frozen = false
	}
	v := bag.ObjectToBag(value)
	if x == nil {
		bx.content = v
	} else {
		x.MustSet(bx.content, v)
	}
}
