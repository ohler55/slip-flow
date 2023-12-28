// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestTaskMetricsOk(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task-flavor :name 'tisk)))
                  (flow-task-metrics task))`,
		Expect: "((received . 0) (processed . 0) (errors . 0))",
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task-flavor :name 'tisk)))
                  (send task :metrics))`,
		Expect: "((received . 0) (processed . 0) (errors . 0))",
	}).Test(t)
}

func TestTaskMetricsNotTask(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-task-metrics t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskMetricsArgCount(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task-flavor :name 'tisk)))
                  (flow-task-metrics task t))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}
