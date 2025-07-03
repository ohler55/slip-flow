// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
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
			f := BoxParse{Function: slip.Function{Name: "flow-box-parse", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-box-parse",
			Args: []*slip.DocArg{
				{
					Name: "box",
					Type: "flow-box",
					Text: "to parse a value in.",
				},
				{
					Name: "string",
					Type: "string",
					Text: "to parse and set in the instance according to the _path_.",
				},
				{Name: "&optional"},
				{
					Name: "path",
					Type: "string|bag-path",
					Text: `to the location in the box to set the parsed value.
The path must follow the JSONPath format.`,
				},
			},
			Return: "box",
			Text: `__flow-box-parse__ parses the _string_ and sets the result at the location
described by _path_. If no _path_ is provided the entire contents of the box is replaced.`,
			Examples: []string{
				`(setq box (make-instance 'flow-box :parse "{a:7}")) => #<flow-box 12345>`,
				`(flow-box-parse box 3 "[a]") => #<flow-box 12345> ;; content is now {a:[3]}`,
			},
		}, &Pkg)
}

// BoxParse represents the flow-box-parse function.
type BoxParse struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *BoxParse) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	self, ok := args[0].(*flavors.Instance)
	if !ok {
		slip.PanicType("box", args[0], "box")
	}
	_ = self.Receive(s, ":parse", args[1:], depth)

	return self
}

type boxParseCaller struct{}

func (caller boxParseCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	switch len(args) {
	case 1:
		parseBox(obj, args[0], nil)
	case 2:
		parseBox(obj, args[0], args[1])
	default:
		slip.PanicMethodArgChoice(obj, ":parse", len(args), "1 or 2")
	}
	return obj
}

func (caller boxParseCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":parse", "flow-box-parse", "flow-box", "box")
}

func parseBox(obj *flavors.Instance, value, path slip.Object) {
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
	ss, ok := value.(slip.String)
	if !ok {
		slip.PanicType("string", value, "string")
	}
	v := sen.MustParse([]byte(ss))
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
