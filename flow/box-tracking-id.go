// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := BoxTrackingID{Function: slip.Function{Name: "flow-box-tracking-id", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-box-tracking-id",
			Args: []*slip.DocArg{
				{
					Name: "box",
					Type: "flow-box",
					Text: "to get the tracking-id of.",
				},
			},
			Return: "object",
			Text:   `__flow-box-tracking-id__ returns the tracking identifier of the box.`,
			Examples: []string{
				`(setq box (make-instance 'flow-box :tracking-id "abc"))`,
				`(flow-box-tracking-id box) => "abc"`,
			},
		}, &Pkg)
}

// BoxTrackingID represents the flow-box-tracking-id function.
type BoxTrackingID struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *BoxTrackingID) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	self, ok := args[0].(*flavors.Instance)
	if !ok {
		slip.PanicType("box", args[0], "box")
	}
	return self.Receive(s, ":tracking-id", args[1:], depth)
}

type boxTrackingIDCaller struct{}

func (caller boxTrackingIDCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)

	return obj.Any.(*box).track.id
}

func (caller boxTrackingIDCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":tracking-id", "flow-box-tracking-id", "flow-box", "box")
}
