// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/bag"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := BoxHas{Function: slip.Function{Name: "flow-box-has", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-box-has",
			Args: []*slip.DocArg{
				{
					Name: "box",
					Type: "flow-box",
					Text: "to has a value from.",
				},
				{
					Name: "path",
					Type: "string|bag-path",
					Text: `to the location in the box to check the value of.
The path must follow the JSONPath format.`,
				},
			},
			Return: "boolean",
			Text:   `__flow-box-has__ returns true if a value at the location described by _path_ exists.`,
			Examples: []string{
				`(setq box (make-instance 'flow-box-flavor :has "{a:7}"))`,
				`(flow-box-has box "a") => t`,
				`(flow-box-has box "b") => nil`,
			},
		}, &Pkg)
}

// BoxHas represents the flow-box-has function.
type BoxHas struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *BoxHas) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	self, ok := args[0].(*flavors.Instance)
	if !ok {
		slip.PanicType("box", args[0], "box")
	}
	return self.Receive(s, ":has", args[1:], depth)
}

type boxHasCaller struct{}

func (caller boxHasCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)
	if len(args) == 1 {
		value = hasBox(obj, args[0])
	} else {
		slip.PanicMethodArgChoice(obj, ":has", len(args), "1")
	}
	return
}

func (caller boxHasCaller) Docs() string {
	return methodDocFromFunc(":has", "flow-box-has", "flow-box-flavor", "box")
}

func hasBox(obj *flavors.Instance, path slip.Object) slip.Object {
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
	if x == nil || x.Has(obj.Any.(*box).content) {
		return slip.True
	}
	return nil
}
