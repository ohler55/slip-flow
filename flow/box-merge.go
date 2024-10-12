// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := BoxMerge{Function: slip.Function{Name: "flow-box-merge", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-box-merge",
			Args: []*slip.DocArg{
				{
					Name: "box",
					Type: "flow-box",
					Text: "to merge a value in.",
				},
				{
					Name: "other",
					Type: "flow-box",
					Text: "the box to merge with _box_.",
				},
			},
			Return: "box",
			Text:   `__flow-box-merge__ merges an _other_ with _box_. Arrays in the _box_ are not expanded.`,
			Examples: []string{
				`(setq box (make-instance 'flow-box-flavor :parse "{a:7}")) => #<flow-box-flavor 12345>`,
				`(setq other (make-instance 'flow-box-flavor :parse "{b:8}")) => #<flow-box-flavor 12346>`,
				`(flow-box-merge box other) => #<flow-box-flavor 12345> ;; content is now {a:7 b:8}`,
			},
		}, &Pkg)
}

// BoxMerge represents the flow-box-merge function.
type BoxMerge struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *BoxMerge) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	self, ok := args[0].(*flavors.Instance)
	if !ok {
		slip.PanicType("box", args[0], "box")
	}
	_ = self.Receive(s, ":merge", args[1:], depth)

	return self
}

type boxMergeCaller struct{}

func (caller boxMergeCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	other, _ := args[0].(*flavors.Instance)
	if other == nil || other.Type != boxFlavor {
		slip.PanicType("box", args[0], "box")
	}
	obj.Any.(*box).merge(other.Any.(*box))

	return obj
}

func (caller boxMergeCaller) Docs() string {
	return methodDocFromFunc(":merge", "flow-box-merge", "flow-box-flavor", "box")
}

func (bx *box) merge(other *box) {
	if bx.frozen {
		bx.content = alt.Dup(bx.content)
		bx.frozen = false
	}
	bx.track.merge(&other.track)

	jp.Walk(other.content, func(path jp.Expr, value any) {
		defer func() { _ = recover() }() // ignore errors
		if !path.Has(bx.content) {
			_ = path.Set(bx.content, value)
		}
	}, true)
}
