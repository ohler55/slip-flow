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
			f := BoxGet{Function: slip.Function{Name: "flow-box-get", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-box-get",
			Args: []*slip.DocArg{
				{
					Name: "box",
					Type: "flow-box",
					Text: "to get a value from.",
				},
				{Name: "&optional"},
				{
					Name: "path",
					Type: "string|bag-path",
					Text: `to the location in the box to get the value.
The path must follow the JSONPath format.`,
				},
				{
					Name: "as-bag",
					Type: "boolean",
					Text: `if not nil then the returned value is a _bag_ otherwise a new LISP value is returned.`,
				},
			},
			Return: "object",
			Text: `__flow-box-get__ gets the value at the location
described by _path_. If no _path_ is provided the entire contents of the box is returned.`,
			Examples: []string{
				`(setq box (make-instance 'flow-box-flavor :get "{a:7}"))`,
				`(flow-box-get box "a") => 7`,
			},
		}, &Pkg)
}

// BoxGet represents the flow-box-get function.
type BoxGet struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *BoxGet) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	self, ok := args[0].(*flavors.Instance)
	if !ok {
		slip.PanicType("box", args[0], "box")
	}
	return self.Receive(s, ":get", args[1:], depth)
}

type boxGetCaller struct{}

func (caller boxGetCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)
	switch len(args) {
	case 0:
		value = getBox(obj, nil, false)
	case 1:
		value = getBox(obj, args[0], false)
	case 2:
		value = getBox(obj, args[0], args[1] != nil)
	default:
		flavors.PanicMethodArgCount(obj, ":get", len(args), 0, 2)
	}
	return
}

func (caller boxGetCaller) Docs() string {
	return methodDocFromFunc(":get", "flow-box-get", "flow-box-flavor", "box")
}

func getBox(obj *flavors.Instance, path slip.Object, asBag bool) slip.Object {
	var x jp.Expr
	switch p := path.(type) {
	case nil:
	case slip.String:
		x = jp.MustParseString(string(p))
	case bag.Path:
		x = jp.Expr(p)
	default:
		slip.PanicType("path", p, "string", "bag-path")
	}
	bx := obj.Any.(*box)
	var value any
	if x == nil {
		value = bx.content
	} else {
		value = x.First(bx.content)
	}
	if value == nil {
		return nil
	}
	if asBag {
		obj = bag.Flavor().MakeInstance().(*flavors.Instance)
		if bx.frozen {
			value = alt.Dup(value)
		}
		obj.Any = value

		return obj
	}
	return slip.SimpleObject(value)
}
