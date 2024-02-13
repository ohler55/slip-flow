// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestFlowSetEntryFoundFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow :name 'flo)))
                  (flow-add-task flow :name "tisk" :actor (lambda (b) (list 'ok b)))
                  (flow-set-entry flow "tisk"))`,
		Expect: "/#<flow-task [0-9a-f]+>/",
	}).Test(t)
}

func TestFlowSetEntryFoundSend(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow :name 'flo)))
                  (send flow :add-task :name "tisk" :actor (lambda (b) (list 'ok b)))
                  (send flow :set-entry 'tisk))`,
		Expect: "/#<flow-task [0-9a-f]+>/",
	}).Test(t)
}

func TestFlowSetEntryNil(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow :name 'flo)))
                  (send flow :add-task :name "tisk" :actor (lambda (b) (list 'ok b)))
                  (send flow :set-entry 'tisk)
                  (send flow :set-entry nil))`,
		Expect: "nil",
	}).Test(t)
}

func TestFlowSetEntryNotFound(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow :name 'flo)))
                  (flow-add-task flow :name "tisk" :actor (lambda (b) (list 'ok b)))
                  (flow-set-entry flow "who"))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}

func TestFlowSetEntryNotFlow(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-set-entry t "tisk")`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowSetEntryBadTask(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow :name 'flo)))
                  (flow-set-entry flow t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
