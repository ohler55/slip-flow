// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := BoxCopy{Function: slip.Function{Name: "flow-box-copy", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-box-copy",
			Args: []*slip.DocArg{
				{
					Name: "box",
					Type: "flow-box",
					Text: "to copy.",
				},
			},
			Return: "box",
			Text: `__flow-box-copy__ makes a copy of the box with shared content.
Both the box and the copy are frozen.`,
			Examples: []string{
				`(setq box (make-instance 'flow-box-flavor :parse "{a:7}"))`,
				`(flow-box-copy box) => #<flow-box-flavor 12346> ;; with the same content as box`,
			},
		}, &Pkg)
}

// BoxCopy represents the flow-box-copy function.
type BoxCopy struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *BoxCopy) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	self, ok := args[0].(*flavors.Instance)
	if !ok {
		slip.PanicType("box", args[0], "box")
	}
	return self.Receive(s, ":copy", args[1:], depth)
}

type boxCopyCaller struct{}

func (caller boxCopyCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	orig := obj.Any.(*box)
	inst := boxFlavor.MakeInstance().(*flavors.Instance)
	bx := &box{track: track{id: orig.track.id}, frozen: true, content: orig.content}
	orig.frozen = true
	bx.track.history = make([]*event, len(orig.track.history))
	for i, ev := range orig.track.history {
		nev := *ev
		bx.track.history[i] = &nev
	}
	inst.Any = bx

	return inst
}

func (caller boxCopyCaller) Docs() string {
	return methodDocFromFunc(":copy", "flow-box-copy", "flow-box-flavor", "box")
}
