// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := TaskMetrics{Function: slip.Function{Name: "flow-task-metrics", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-task-metrics",
			Args: []*slip.DocArg{
				{
					Name: "task",
					Type: "flow-task",
					Text: "to return the metrics for.",
				},
			},
			Return: "list",
			Text:   `__flow-task-metrics__ returns the metrics the _task_.`,
			Examples: []string{
				`(setq task (make-instance 'flow-task-flavor :name "tisk"))`,
				`(flow-task-metrics task) => ((received . 0) (processed . 0) (errors . 0))`,
			},
		}, &Pkg)
}

// TaskMetrics represents the flow-task-metrics function.
type TaskMetrics struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *TaskMetrics) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.ArgCountCheck(f, args, 1, 1)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Flavor != taskFlavor {
		slip.PanicType("task", args[0], "task")
	}
	return self.Any.(*task).metrics()
}

type taskMetricsCaller struct{}

func (caller taskMetricsCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	return obj.Any.(*task).metrics()
}

func (caller taskMetricsCaller) Docs() string {
	return methodDocFromFunc(":metrics", "flow-task-metrics", "flow-task-flavor", "task")
}
