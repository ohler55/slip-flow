// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := FlowEntry{Function: slip.Function{Name: "flow-entry", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-entry",
			Args: []*slip.DocArg{
				{
					Name: "flow",
					Type: "flow",
					Text: "to return the entry task from.",
				},
			},
			Return: "instance",
			Text:   `__flow-entry__ returns the entry task in the _flow_.`,
			Examples: []string{
				`(setq flow (make-instance 'flow-flavor :entry "flo"))`,
				`(flow-add-task flow :name 'tisk :actor (lambda (b) (list 'ok b)))`,
				`(flow-set-entry flow "tisk") => (#<flow-task-flavor 12345>)`,
				`(flow-entry flow) => (#<flow-task-flavor 12345>)`,
			},
		}, &Pkg)
}

// FlowEntry represents the flow-entry function.
type FlowEntry struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *FlowEntry) Call(s *slip.Scope, args slip.List, depth int) (entry slip.Object) {
	slip.ArgCountCheck(f, args, 1, 1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Type != flowFlavor {
		slip.PanicType("flow", args[0], "flow")
	}
	if self.Any.(*flow).entry != nil {
		entry = self.Any.(*flow).entry.self
	}
	return
}

type flowEntryCaller struct{}

func (caller flowEntryCaller) Call(s *slip.Scope, args slip.List, _ int) (entry slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	if obj.Any.(*flow).entry != nil {
		entry = obj.Any.(*flow).entry.self
	}
	return
}

func (caller flowEntryCaller) Docs() string {
	return methodDocFromFunc(":entry", "flow-entry", "flow-flavor", "flow")
}
