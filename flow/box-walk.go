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
			f := BoxWalk{Function: slip.Function{Name: "flow-box-walk", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-box-walk",
			Args: []*slip.DocArg{
				{
					Name: "box",
					Type: "flow-box",
					Text: "to walk.",
				},
				{
					Name: "function",
					Type: "function",
					Text: "to apply to each node in the instance matching the _path_.",
				},
				{Name: "&optional"},
				{
					Name: "path",
					Type: "string|bag-path",
					Text: `to the location in the box to walk. Default: "..".
The path must follow the JSONPath format.`,
				},
				{Name: "&key"},
				{
					Name: "as-lisp",
					Type: "boolean",
					Text: `if not nil then the value to the _function_ is a LISP value otherwise a new _bag_.`,
				},
			},
			Return: "box",
			Text:   `__flow-box-walk__ walks the values at the location described by _path_.`,
			Examples: []string{
				`(setq box (make-instance 'flow-box-flavor :parse "{a:[1 2 3]}")) => #<flow-box-flavor 12345>`,
				`(flow-box-walk box 'reverse "a") => #<flow-box-flavor 12345> ;; content is now {a:[3 2 1]}`,
			},
		}, &Pkg)
}

// BoxWalk represents the flow-box-walk function.
type BoxWalk struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *BoxWalk) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	self, ok := args[0].(*flavors.Instance)
	if !ok {
		slip.PanicType("box", args[0], "box")
	}
	_ = self.Receive(s, ":walk", args[1:], depth)

	return self
}

type boxWalkCaller struct{}

func (caller boxWalkCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	walkBox(s, obj, args, depth)
	return nil
}

func (caller boxWalkCaller) Docs() string {
	return methodDocFromFunc(":walk", "flow-box-walk", "flow-box-flavor", "box")
}

func walkBox(s *slip.Scope, obj *flavors.Instance, args slip.List, depth int) {
	fn := args[0]
	path := jp.D()
	var asBag bool
	if 1 < len(args) {
		switch p := args[1].(type) {
		case nil:
		case slip.String:
			path = jp.MustParseString(string(p))
		case bag.Path:
			path = jp.Expr(p)
		default:
			slip.PanicType("path", p, "string", "bag-path")
		}
		asBag = 2 < len(args) && args[2] != nil
	}
	d2 := depth + 1
CallFunc:
	switch tf := fn.(type) {
	case *slip.Lambda:
		if asBag {
			for _, v := range path.Get(obj.Any.(*box).content) {
				arg := bag.Flavor().MakeInstance().(*flavors.Instance)
				arg.Any = v
				_ = tf.Call(s, slip.List{arg}, d2)
			}
		} else {
			for _, v := range path.Get(obj.Any.(*box).content) {
				arg := slip.SimpleObject(v)
				_ = tf.Call(s, slip.List{arg}, d2)
			}
		}
	case *slip.FuncInfo:
		if asBag {
			for _, v := range path.Get(obj.Any.(*box).content) {
				arg := bag.Flavor().MakeInstance().(*flavors.Instance)
				arg.Any = v
				_ = tf.Apply(s, slip.List{arg}, d2)
			}
		} else {
			for _, v := range path.Get(obj.Any.(*box).content) {
				arg := slip.SimpleObject(v)
				tf.Apply(s, slip.List{arg}, d2)
			}
		}
	case slip.Symbol:
		fn = slip.FindFunc(string(tf))
		goto CallFunc
	case slip.List:
		fn = s.Eval(tf, d2)
		goto CallFunc
	default:
		slip.PanicType("function", tf, "function")
	}
}
