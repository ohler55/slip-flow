// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestFlowRemoveTaskNotFlow(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-remove-task t "tisk")`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowRemoveTaskBadTask(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow-flavor :name 'flo)))
                  (flow-remove-task flow t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowRemoveTaskOk(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow-flavor :name 'flo)))
                  (flow-add-task flow :name "tick" :actor (lambda (b) (list 'ok b)))
                  (flow-add-task flow :name "tock" :actor (lambda (b) (list 'ok b)))
                  (flow-remove-task flow "tick")
                  (send flow :remove-task 'tock)
                  (flow-tasks flow))`,
		Expect: "nil",
	}).Test(t)
}
