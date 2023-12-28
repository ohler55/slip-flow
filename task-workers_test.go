// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestTaskWorkersOk(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task-flavor :name 'tisk :workers 3)))
                  (flow-task-workers task))`,
		Expect: "3",
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task-flavor :name 'tisk :workers 3)))
                  (send task :workers))`,
		Expect: "3",
	}).Test(t)
}

func TestTaskWorkersNotTask(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-task-workers t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskWorkersArgCount(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task-flavor :name 'tisk :workers 3)))
                  (flow-task-workers task t))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}
