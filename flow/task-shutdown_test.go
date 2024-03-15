// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestTaskShutdownOk(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task :name 'tisk :workers 3 :actor (lambda (b) nil))))
                  (flow-task-start task)
                  (flow-task-shutdown task)
                  (flow-task-running task))`,
		Expect: "nil",
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task :name 'tisk :workers 3 :actor (lambda (b) nil))))
                  (send task :start)
                  (send task :shutdown)
                  (send task :running))`,
		Expect: "nil",
	}).Test(t)
}

func TestTaskShutdownNotTask(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-task-shutdown t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskShutdownArgCount(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task :name 'tisk :workers 3)))
                  (flow-task-shutdown task t))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}
