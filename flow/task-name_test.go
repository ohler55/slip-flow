// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip-flow/flow"
	"github.com/ohler55/slip/sliptest"
)

func TestTaskNameOk(t *testing.T) {
	scope := slip.NewScope()
	task, _ := flow.MakeTask(scope, 0, slip.Symbol(":name"), slip.String("tisk"))
	scope.Let("task", task)
	(&sliptest.Function{
		Scope:  scope,
		Source: `(flow-task-name task)`,
		Expect: `"tisk"`,
	}).Test(t)
	(&sliptest.Function{
		Scope:  scope,
		Source: `(send task :name)`,
		Expect: `"tisk"`,
	}).Test(t)
}

func TestTaskNameNotTask(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-task-name t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskNameArgCount(t *testing.T) {
	scope := slip.NewScope()
	task, _ := flow.MakeTask(scope, 0, slip.Symbol(":name"), slip.String("tisk"))
	scope.Let("task", task)
	(&sliptest.Function{
		Scope:     scope,
		Source:    `(flow-task-name task t)`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}
