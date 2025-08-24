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
			f := BoxRemove{Function: slip.Function{Name: "flow-box-remove", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-box-remove",
			Args: []*slip.DocArg{
				{
					Name: "box",
					Type: "flow-box",
					Text: "to remove a value from.",
				},
				{
					Name: "path",
					Type: "string|bag-path",
					Text: `to the location in the box to remove.
The path must follow the JSONPath format.`,
				},
			},
			Return: "box",
			Text:   `__flow-box-remove__ returns _box_ after removing all values that match the provided _path_.`,
			Examples: []string{
				`(setq box (make-instance 'flow-box :parse "{a:7}")) => #<flow-box 12345>`,
				`(flow-box-remove box "a") => #<flow-box 12345> ;; content is now {}`,
			},
		}, &Pkg)
}

// BoxRemove represents the flow-box-remove function.
type BoxRemove struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *BoxRemove) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	self, ok := args[0].(*flavors.Instance)
	if !ok {
		slip.TypePanic(s, depth, "box", args[0], "box")
	}
	_ = self.Receive(s, ":remove", args[1:], depth)

	return self
}

type boxRemoveCaller struct{}

func (caller boxRemoveCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	if len(args) == 1 {
		removeBox(s, obj, args[0], depth)
	} else {
		slip.MethodArgChoicePanic(s, depth, obj, ":remove", len(args), "1")
	}
	return obj
}

func (caller boxRemoveCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":remove", "flow-box-remove", "flow-box", "box")
}

func removeBox(s *slip.Scope, obj *flavors.Instance, path slip.Object, depth int) {
	var x jp.Expr
	switch p := path.(type) {
	case nil:
	case slip.String:
		x = jp.MustParseString(string(p))
	case bag.Path:
		x = jp.Expr(p)
	default:
		slip.TypePanic(s, depth, "path", p, "string")
	}
	bx := obj.Any.(*box)
	if x == nil {
		bx.content = nil
	} else {
		if bx.frozen {
			bx.content = alt.Dup(bx.content)
		}
		bx.content = x.MustRemove(bx.content)
	}
	bx.frozen = false
}
