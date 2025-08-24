// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := FlowShutdown{Function: slip.Function{Name: "flow-shutdown", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-shutdown",
			Args: []*slip.DocArg{
				{
					Name: "flow",
					Type: "flow",
					Text: "to shutdown.",
				},
			},
			Return: "nil",
			Text:   `__flow-shutdown__ shutdowns all the _tasks_ in the _flow_.`,
			Examples: []string{
				`(setq flow (make-instance 'flow :shutdown "flo"))`,
				`(flow-add-task flow :name 'tisk :actor (lambda (b) (list 'ok b)))`,
				`(flow-start flow) => nil`,
				`(flow-shutdown flow) => nil`,
			},
		}, &Pkg)
}

// FlowShutdown represents the flow-shutdown function.
type FlowShutdown struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *FlowShutdown) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.ArgCountCheck(f, args, 1, 1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Type != flowFlavor {
		slip.TypePanic(s, depth, "flow", args[0], "flow")
	}
	self.Any.(*flow).shutdown(s)

	return nil
}

type flowShutdownCaller struct{}

func (caller flowShutdownCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)

	obj.Any.(*flow).shutdown(s)

	return nil
}

func (caller flowShutdownCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":shutdown", "flow-shutdown", "flow", "flow")
}
