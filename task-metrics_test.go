// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestTaskMetricsOk(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task-flavor :name 'tisk :actor (lambda (b) (list 'ok b)))))
                  (flow-task-receive task (make-flow-box :set '(1 2 3)))
                  (flow-task-metrics task))`,
		Expect: `/\(\(received . 1\) \(processed . 1\) \(errors . 0\) \(average . .+\)\)/`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task-flavor :name 'tisk :actor (lambda (b) (list 'ok b)))))
                  (send task :receive (make-flow-box :set '(1 2 3)))
                  (send task :metrics))`,
		Expect: `/\(\(received . 1\) \(processed . 1\) \(errors . 0\) \(average . .+\)\)/`,
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
