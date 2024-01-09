// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := FlowSetEntry{Function: slip.Function{Name: "flow-set-entry", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-set-entry",
			Args: []*slip.DocArg{
				{
					Name: "flow",
					Type: "flow",
					Text: "to return the set-entry task from.",
				},
				{
					Name: "task-name",
					Type: "string",
					Text: "of the task to set as the entry to the flow.",
				},
			},
			Return: "instance",
			Text:   `__flow-set-entry__ returns the task set to be the entry task of the _flow_.`,
			Examples: []string{
				`(setq flow (make-instance 'flow-flavor :set-entry "flo"))`,
				`(flow-add-task flow :name 'tisk :actor (lambda (b) (list 'ok b)))`,
				`(flow-set-entry flow "tisk") => (#<flow-task-flavor 12345>)`,
				`(flow-entry flow) => (#<flow-task-flavor 12345>)`,
			},
		}, &Pkg)
}

// FlowSetEntry represents the flow-flow-set-entry function.
type FlowSetEntry struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *FlowSetEntry) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.ArgCountCheck(f, args, 2, 2)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Flavor != flowFlavor {
		slip.PanicType("flow", args[0], "flow")
	}
	return self.Any.(*flow).setEntry(args[1])
}

type flowSetEntryCaller struct{}

func (caller flowSetEntryCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)

	return obj.Any.(*flow).setEntry(args[0])
}

func (caller flowSetEntryCaller) Docs() string {
	return methodDocFromFunc(":set-entry", "flow-set-entry", "flow-flavor", "flow")
}
