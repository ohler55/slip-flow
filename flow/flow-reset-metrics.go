// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := FlowResetMetrics{Function: slip.Function{Name: "flow-reset-metrics", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-reset-metrics",
			Args: []*slip.DocArg{
				{
					Name: "flow",
					Type: "flow",
					Text: "to reset the metrics of.",
				},
			},
			Return: "nil",
			Text:   `__flow-reset-metrics__ resets the metrics of the _flow_.`,
			Examples: []string{
				`(setq flow (make-instance 'flow :name "flo"))`,
				`(flow-reset-metrics flow) => nil`,
				`(flow-metrics flow) => ((received . 0) (processed . 0) (errors . 0))`,
			},
		}, &Pkg)
}

// FlowResetMetrics represents the flow-reset-metrics function.
type FlowResetMetrics struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *FlowResetMetrics) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.ArgCountCheck(f, args, 1, 1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Type != flowFlavor {
		slip.TypePanic(s, depth, "flow", args[0], "flow")
	}
	self.Any.(*flow).resetMetrics()

	return nil
}

type flowResetMetricsCaller struct{}

func (caller flowResetMetricsCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	obj.Any.(*flow).resetMetrics()

	return nil
}

func (caller flowResetMetricsCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":reset-metrics", "flow-reset-metrics", "flow", "flow")
}
