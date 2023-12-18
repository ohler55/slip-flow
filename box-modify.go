// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main

import (
	"strings"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/bag"
	"github.com/ohler55/slip/pkg/cl"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := BoxModify{Function: slip.Function{Name: "flow-box-modify", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-box-modify",
			Args: []*slip.DocArg{
				{
					Name: "box",
					Type: "flow-box",
					Text: "to modify a value in.",
				},
				{
					Name: "function",
					Type: "function",
					Text: "to modify the value at _path_. It must expect a single argument.",
				},
				{Name: "&optional"},
				{
					Name: "path",
					Type: "string|bag-path",
					Text: `to the location in the box to modify the _value_.
The path must follow the JSONPath format.`,
				},
				{Name: "&key"},
				{
					Name: "as-bag",
					Type: "boolean",
					Text: `if true the _function_ expects a _bag_ otherwise it expects a lisp object.`,
				},
			},
			Return: "box",
			Text: `__flow-box-modify__ modifies a _value_ at the location described by _path_
using the _function_ specified. If no _path_ is provided the entire contents is passed to the _function_.`,
			Examples: []string{
				`(setq box (make-instance 'flow-box-flavor :parse "{a:[1 2 3]}")) => #<flow-box-flavor 12345>`,
				`(flow-box-modify box 'reverse "a") => #<flow-box-flavor 12345> ;; content is now {a:[3 2 1]}`,
			},
		}, &Pkg)
}

// BoxModify represents the flow-box-modify function.
type BoxModify struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *BoxModify) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	self, ok := args[0].(*flavors.Instance)
	if !ok {
		slip.PanicType("box", args[0], "box")
	}
	_ = self.Receive(s, ":modify", args[1:], depth)

	return self
}

type boxModifyCaller struct{}

func (caller boxModifyCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	modifyBox(s, obj, args, depth+1)
	return obj
}

func (caller boxModifyCaller) Docs() string {
	return methodDocFromFunc(":modify", "flow-box-modify", "flow-box-flavor", "box")
}

func modifyBox(s *slip.Scope, obj *flavors.Instance, args slip.List, depth int) {
	caller := cl.ResolveToCaller(s, args[0], depth)
	var (
		x     jp.Expr
		asBag bool
	)
	if 1 < len(args) {
		switch p := args[1].(type) {
		case nil:
		case slip.String:
			x = jp.MustParseString(string(p))
		case bag.Path:
			x = jp.Expr(p)
		default:
			slip.PanicType("path", p, "string")
		}
		if 2 < len(args) {
			for pos := 2; pos < len(args); pos += 2 {
				sym, ok := args[pos].(slip.Symbol)
				if !ok {
					slip.PanicType("keyword", args[pos], "keyword")
				}
				if len(args)-1 <= pos {
					slip.NewPanic("keyword %s is missing a value", sym)
				}
				if strings.EqualFold(string(sym), ":as-bag") {
					asBag = args[pos+1] != nil
				} else {
					slip.PanicType("keyword", sym, ":as-bag")
				}
			}
		}
	}
	bx := obj.Any.(*box)
	if bx.frozen {
		bx.content = alt.Dup(bx.content)
		bx.frozen = false
	}
	if x == nil {
		bx.content = modifyValue(s, bx.content, caller, asBag, depth)
	} else {
		bx.content = x.MustModify(bx.content, func(element any) (altered any, changed bool) {
			return modifyValue(s, element, caller, asBag, depth), true
		})
	}
}

func modifyValue(s *slip.Scope, value any, caller slip.Caller, asBag bool, depth int) any {
	bagFlavor := bag.Flavor()
	if asBag {
		bg := bagFlavor.MakeInstance().(*flavors.Instance)
		bg.Any = value
		obj := caller.Call(s, slip.List{bg}, depth)
		if bg, _ := obj.(*flavors.Instance); bg != nil && bg.Flavor == bagFlavor {
			return bg.Any
		}
		return slip.Simplify(obj)
	}
	obj := slip.SimpleObject(value)
	obj = caller.Call(s, slip.List{obj}, depth)
	if bg, _ := obj.(*flavors.Instance); bg != nil && bg.Flavor == bagFlavor {
		return bg.Any
	}
	return slip.Simplify(obj)
}
