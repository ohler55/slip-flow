// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestTaskStartOk(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task :name 'tisk :workers 3)))
                  (flow-task-start task)
                  (flow-task-running task))`,
		Expect: "t",
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task :name 'tisk :workers 3)))
                  (send task :start)
                  (send task :running))`,
		Expect: "t",
	}).Test(t)
}

func TestTaskStartNotTask(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-task-start t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskStartArgCount(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task :name 'tisk :workers 3)))
                  (flow-task-start task t))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}
