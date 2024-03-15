// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/bag"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := BoxBag{Function: slip.Function{Name: "flow-box-bag", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-box-bag",
			Args: []*slip.DocArg{
				{
					Name: "box",
					Type: "flow-box",
					Text: "to convert to a bag.",
				},
			},
			Return: "object",
			Text: `__flow-box-bag__ converts the content of an instance of the
_flow-box-flavor_ to an instance of the _bag-flavor_ with the contents of the box.`,
			Examples: []string{
				`(setq box (make-instance 'flow-box-flavor :parse "{a:7}"))`,
				`(flow-box-bag box) => #<bag-flavor 12345> ;; with content {a:7}`,
			},
		}, &Pkg)
}

// BoxBag represents the flow-box-bag function.
type BoxBag struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *BoxBag) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	self, ok := args[0].(*flavors.Instance)
	if !ok {
		slip.PanicType("box", args[0], "box")
	}
	return self.Receive(s, ":bag", args[1:], depth)
}

type boxBagCaller struct{}

func (caller boxBagCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	bg := bag.Flavor().MakeInstance().(*flavors.Instance)
	bx := obj.Any.(*box)
	if bx.frozen {
		bx.content = alt.Dup(bx.content)
		bx.frozen = false
	}
	bg.Any = bx.content

	return bg
}

func (caller boxBagCaller) Docs() string {
	return methodDocFromFunc(":bag", "flow-box-bag", "flow-box-flavor", "box")
}
