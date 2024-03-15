// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := BoxTrack{Function: slip.Function{Name: "flow-box-track", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-box-track",
			Args: []*slip.DocArg{
				{
					Name: "box",
					Type: "flow-box",
					Text: "to get the track from.",
				},
			},
			Return: "flow-track",
			Text:   `__flow-box-track__ returns the _flow-track-flavor_ instance of the box.`,
			Examples: []string{
				`(setq box (make-instance 'flow-box-flavor :tracking-id 123))`,
				`(flow-box-track box) => #<flow-track-flavor 12345>`,
			},
		}, &Pkg)
}

// BoxTrack represents the flow-box-track function.
type BoxTrack struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *BoxTrack) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	self, ok := args[0].(*flavors.Instance)
	if !ok {
		slip.PanicType("box", args[0], "box")
	}
	return self.Receive(s, ":track", args[1:], depth)
}

type boxTrackCaller struct{}

func (caller boxTrackCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	trk := trackFlavor.MakeInstance().(*flavors.Instance)
	trk.Any = &(obj.Any.(*box).track)

	return trk
}

func (caller boxTrackCaller) Docs() string {
	return methodDocFromFunc(":track", "flow-box-track", "flow-box-flavor", "box")
}
