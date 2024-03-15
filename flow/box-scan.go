// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := BoxScan{Function: slip.Function{Name: "flow-box-scan", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-box-scan",
			Args: []*slip.DocArg{
				{
					Name: "box",
					Type: "flow-box",
					Text: "to add the scan event to.",
				},
				{
					Name: "flow-name",
					Type: "string",
					Text: "of the flow the scan was initiated from.",
				},
				{
					Name: "task-name",
					Type: "string",
					Text: "of the task the scan was initiated from.",
				},
			},
			Return: "",
			Text: `__flow-box-scan__ adds an event consisting of the current time,
task name, and flow name to the history of the box.`,
			Examples: []string{
				`(setq box (make-instance 'flow-box-flavor :tracking-id 123))`,
				`(send box :scan "flo" "tisk")`,
				`(flow-box-history box) => ((@2023-12-15T19:23:17Z "tisk" "flo"))`,
			},
		}, &Pkg)
}

// BoxScan represents the flow-box-scan function.
type BoxScan struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *BoxScan) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	self, ok := args[0].(*flavors.Instance)
	if !ok {
		slip.PanicType("box", args[0], "box")
	}
	return self.Receive(s, ":scan", args[1:], depth)
}

type boxScanCaller struct{}

func (caller boxScanCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	var (
		flowName string
		taskName string
	)
	if ss, ok := args[0].(slip.String); ok {
		flowName = string(ss)
	} else {
		slip.PanicType("flow-name", args[0], "string")
	}
	if ss, ok := args[1].(slip.String); ok {
		taskName = string(ss)
	} else {
		slip.PanicType("task-name", args[1], "string")
	}
	obj.Any.(*box).track.Scan(flowName, taskName)

	return nil
}

func (caller boxScanCaller) Docs() string {
	return methodDocFromFunc(":scan", "flow-box-scan", "flow-box-flavor", "box")
}
