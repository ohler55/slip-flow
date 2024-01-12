// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestFlowFindTaskFoundFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow-flavor :name 'flo)))
                  (flow-add-task flow :name "tisk" :actor (lambda (b) (list 'ok b)))
                  (flow-find-task flow "tisk"))`,
		Expect: "/#<flow-task-flavor [0-9a-f]+>/",
	}).Test(t)
}

func TestFlowFindTaskFoundSend(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow-flavor :name 'flo)))
                  (send flow :add-task :name "tisk" :actor (lambda (b) (list 'ok b)))
                  (send flow :find-task 'tisk))`,
		Expect: "/#<flow-task-flavor [0-9a-f]+>/",
	}).Test(t)
}

func TestFlowFindTaskNotFound(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow-flavor :name 'flo)))
                  (flow-add-task flow :name "tisk" :actor (lambda (b) (list 'ok b)))
                  (flow-find-task flow "who"))`,
		Expect: "nil",
	}).Test(t)
}

func TestFlowFindTaskNotFlow(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-find-task t "tisk")`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowFindTaskBadTask(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow-flavor :name 'flo)))
                  (flow-find-task flow t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
