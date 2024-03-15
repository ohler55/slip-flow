// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := FlowLink{Function: slip.Function{Name: "flow-link", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-link",
			Args: []*slip.DocArg{
				{
					Name: "flow",
					Type: "flow",
					Text: "to create a link in.",
				},
				{
					Name: "link-name",
					Type: "string",
					Text: "name of the link to create.",
				},
				{
					Name: "from",
					Type: "string",
					Text: "name of the originating task of a link.",
				},
				{
					Name: "to",
					Type: "string",
					Text: "name of the destination task of a link.",
				},
				{Name: "&optional"},
				{
					Name: "mid-points",
					Type: "list",
					Text: `a list of mid points for drawing the link in an SVG or in an editor. A list
of x and y pairs are expected. e.g., ((100 100) (100 200))`,
				},
			},
			Return: "nil",
			Text:   `__flow-link__ creates a named link between the _from_ task to the _to_ task in a flow.`,
			Examples: []string{
				`(setq flow (make-instance 'flow-flavor :name "flo"))`,
				`(flow-add-task flow :name 'tick :actor (lambda (b) (list 'ok b)))`,
				`(flow-add-task flow :name 'tock :actor (lambda (b) (list 'ok b)))`,
				`(flow-link flow "ok" 'tick 'tock) => nil`,
			},
		}, &Pkg)
}

// FlowLink represents the flow-link function.
type FlowLink struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *FlowLink) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.ArgCountCheck(f, args, 4, 5)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Flavor != flowFlavor {
		slip.PanicType("flow", args[0], "flow")
	}
	self.Any.(*flow).link(args[1:])

	return nil
}

type flowLinkCaller struct{}

func (caller flowLinkCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)

	obj.Any.(*flow).link(args)

	return nil
}

func (caller flowLinkCaller) Docs() string {
	return methodDocFromFunc(":link", "flow-link", "flow-flavor", "flow")
}
