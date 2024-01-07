// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestFlowTasksFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow-flavor :name 'flo)))
                  (flow-add-task flow :name "tick" :actor (lambda (b) (list 'ok b)))
                  (flow-add-task flow :name "tock" :actor (lambda (b) (list 'ok b)))
                  (mapcar (lambda (task) (send task :name)) (flow-tasks flow)))`,
		Expect: `("tick" "tock")`,
	}).Test(t)
}

func TestFlowTasksSend(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow-flavor :name 'flo)))
                  (send flow :add-task :name "tick" :actor (lambda (b) (list 'ok b)))
                  (send flow :add-task :name "tock" :actor (lambda (b) (list 'ok b)))
                  (mapcar (lambda (task) (send task :name)) (send flow :tasks)))`,
		Expect: `("tick" "tock")`,
	}).Test(t)
}

func TestFlowTasksNotFlow(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-tasks t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
