// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"io"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/bag"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := BoxRead{Function: slip.Function{Name: "flow-box-read", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-box-read",
			Args: []*slip.DocArg{
				{
					Name: "box",
					Type: "flow-box",
					Text: "to read a value in.",
				},
				{
					Name: "stream",
					Type: "input-stream",
					Text: "to read from and set in the instance according to the _path_.",
				},
				{Name: "&optional"},
				{
					Name: "path",
					Type: "string|bag-path",
					Text: `to the location in the box to set the readd value.
The path must follow the JSONPath format.`,
				},
			},
			Return: "box",
			Text: `__flow-box-read__ reads from the _stream_ and sets the result at the location
described by _path_. If no _path_ is provided the entire contents of the box is replaced.`,
			Examples: []string{
				`(setq box (make-instance 'flow-box-flavor :read "{a:7}"))`,
				`(flow-box-read box (make-string-input-steam "[3]") "a") => #<flow-box-flavor 12345>`,
				` ;; content is now {a:[3]}`,
			},
		}, &Pkg)
}

// BoxRead represents the flow-box-read function.
type BoxRead struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *BoxRead) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	self, ok := args[0].(*flavors.Instance)
	if !ok {
		slip.PanicType("box", args[0], "box")
	}
	_ = self.Receive(s, ":read", args[1:], depth)

	return self
}

type boxReadCaller struct{}

func (caller boxReadCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	switch len(args) {
	case 1:
		readBox(obj, args[0], nil)
	case 2:
		readBox(obj, args[0], args[1])
	default:
		flavors.PanicMethodArgChoice(obj, ":read", len(args), "1 or 2")
	}
	return obj
}

func (caller boxReadCaller) Docs() string {
	return methodDocFromFunc(":parse", "flow-box-parse", "flow-box-flavor", "box")
}

func readBox(obj *flavors.Instance, value, path slip.Object) {
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
	r, ok := value.(io.Reader)
	if !ok {
		slip.PanicType("stream", value, "input-stream")
	}
	v := sen.MustParseReader(r)
	if options.Converter != nil {
		v = options.Converter.Convert(v)
	}
	bx := obj.Any.(*box)
	if bx.frozen {
		bx.content = alt.Dup(bx.content)
		bx.frozen = false
	}
	if x == nil {
		bx.content = v
	} else {
		x.MustSet(bx.content, v)
	}
}
