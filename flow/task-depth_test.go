// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestTaskDepthOk(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task :name 'tisk :depth 3)))
                  (flow-task-depth task))`,
		Expect: "3",
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task :name 'tisk :depth 3)))
                  (send task :depth))`,
		Expect: "3",
	}).Test(t)
}

func TestTaskDepthNotTask(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-task-depth t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskDepthArgCount(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task :name 'tisk :depth 3)))
                  (flow-task-depth task t))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}
