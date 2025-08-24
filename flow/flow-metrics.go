// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := FlowMetrics{Function: slip.Function{Name: "flow-metrics", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-metrics",
			Args: []*slip.DocArg{
				{
					Name: "flow",
					Type: "flow",
					Text: "to return the metrics for.",
				},
			},
			Return: "list",
			Text:   `__flow-metrics__ returns the metrics the _flow_.`,
			Examples: []string{
				`(setq flow (make-instance 'flow :name "flo"))`,
				`(flow-metrics flow) => ((received . 0) (processed . 0) (errors . 0))`,
			},
		}, &Pkg)
}

// FlowMetrics represents the flow-metrics function.
type FlowMetrics struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *FlowMetrics) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.CheckArgCount(s, depth, f, args, 1, 1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Type != flowFlavor {
		slip.TypePanic(s, depth, "flow", args[0], "flow")
	}
	return self.Any.(*flow).metrics()
}

type flowMetricsCaller struct{}

func (caller flowMetricsCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	return obj.Any.(*flow).metrics()
}

func (caller flowMetricsCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":metrics", "flow-metrics", "flow", "flow")
}
