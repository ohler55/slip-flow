// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestTaskRunningOk(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task :name 'tisk :workers 3)))
                  (flow-task-running task))`,
		Expect: "nil",
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task :name 'tisk :workers 3)))
                  (flow-task-start task)
                  (flow-task-running task))`,
		Expect: "t",
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task :name 'tisk :workers 3)))
                  (send task :running))`,
		Expect: "nil",
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task :name 'tisk :workers 3)))
                  (send task :start)
                  (send task :running))`,
		Expect: "t",
	}).Test(t)
}

func TestTaskRunningNotTask(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-task-running t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskRunningArgCount(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task :name 'tisk :workers 3)))
                  (flow-task-running task t))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}
