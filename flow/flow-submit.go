// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := FlowSubmit{Function: slip.Function{Name: "flow-submit", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-submit",
			Args: []*slip.DocArg{
				{
					Name: "flow",
					Type: "flow",
					Text: "to submit a box to.",
				},
				{
					Name: "box",
					Type: "instance",
					Text: "instance of the flow-box-flavor to process by the flow.",
				},
				{Name: "&optional"},
				{
					Name: "watch",
					Type: "symbol bound to a gi:channel",
					Text: "Add a watcher to the box with the name of the symbol which must be bound to a gi:channel",
				},
			},
			Return: "box",
			Text:   `__flow-submit__ submits an instance of the _flow-box-flavor_ for processing by the _flow_.`,
			Examples: []string{
				`(setq flow (make-instance 'flow :submit "flo"))`,
				`(flow-add-task flow :name 'tisk :actor (lambda (b) (list 'ok b)))`,
				`(flow-submit flow "tisk") => nil`,
			},
		}, &Pkg)
}

// FlowSubmit represents the flow-submit function.
type FlowSubmit struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *FlowSubmit) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.ArgCountCheck(f, args, 2, 3)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Type != flowFlavor {
		slip.TypePanic(s, depth, "flow", args[0], "flow")
	}
	var watcher slip.Object
	if 2 < len(args) {
		watcher = args[2]
	}
	return self.Any.(*flow).submit(s, args[1], watcher, depth)
}

type flowSubmitCaller struct{}

func (caller flowSubmitCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	var watcher slip.Object
	if 1 < len(args) {
		watcher = args[1]
	}
	return obj.Any.(*flow).submit(s, args[0], watcher, depth)
}

func (caller flowSubmitCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":submit", "flow-submit", "flow", "flow")
}
