// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := TaskResetMetrics{Function: slip.Function{Name: "flow-task-reset-metrics", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-task-reset-metrics",
			Args: []*slip.DocArg{
				{
					Name: "task",
					Type: "flow-task",
					Text: "to reset the metrics of.",
				},
			},
			Return: "nil",
			Text:   `__flow-task-reset-metrics__ resets the metrics of the _task_.`,
			Examples: []string{
				`(setq task (make-instance 'flow-task :name "tisk"))`,
				`(flow-task-reset-metrics task) => nil`,
				`(flow-task-metrics task) => ((received . 0) (processed . 0) (errors . 0))`,
			},
		}, &Pkg)
}

// TaskResetMetrics represents the flow-task-reset-metrics function.
type TaskResetMetrics struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *TaskResetMetrics) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.CheckArgCount(s, depth, f, args, 1, 1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Type != taskFlavor {
		slip.TypePanic(s, depth, "task", args[0], "task")
	}
	self.Any.(*task).resetMetrics()

	return nil
}

type taskResetMetricsCaller struct{}

func (caller taskResetMetricsCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	obj.Any.(*task).resetMetrics()

	return nil
}

func (caller taskResetMetricsCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":reset-metrics", "flow-task-reset-metrics", "flow-task", "task")
}
