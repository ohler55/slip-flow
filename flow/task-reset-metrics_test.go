// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestTaskResetMetricsOk(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task :name 'tisk :actor (lambda (b) (list 'ok b)))))
                  (flow-task-receive task (make-flow-box :set '(1 2 3)))
                  (flow-task-reset-metrics task)
                  (flow-task-metrics task))`,
		Expect: `((received . 0) (processed . 0) (errors . 0))`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task :name 'tisk :actor (lambda (b) (list 'ok b)))))
                  (send task :receive (make-flow-box :set '(1 2 3)))
                  (send task :reset-metrics)
                  (send task :metrics))`,
		Expect: `((received . 0) (processed . 0) (errors . 0))`,
	}).Test(t)
}

func TestTaskResetMetricsNotTask(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-task-reset-metrics t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskResetMetricsArgCount(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task :name 'tisk)))
                  (flow-task-reset-metrics task t))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}
