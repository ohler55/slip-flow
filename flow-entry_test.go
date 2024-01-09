// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestFlowEntryFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow-flavor :name 'flo)))
                  (flow-add-task flow :name "tisk" :actor (lambda (b) (list 'ok b)))
                  (flow-set-entry flow "tisk")
                  (flow-entry flow))`,
		Expect: "/#<flow-task-flavor [0-9a-f]+>/",
	}).Test(t)
}

func TestFlowEntryFoundSend(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow-flavor :name 'flo)))
                  (send flow :add-task :name "tisk" :actor (lambda (b) (list 'ok b)))
                  (send flow :set-entry 'tisk)
                  (send flow :entry))`,
		Expect: "/#<flow-task-flavor [0-9a-f]+>/",
	}).Test(t)
}

func TestFlowEntryNil(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow-flavor :name 'flo)))
                  (send flow :add-task :name "tisk" :actor (lambda (b) (list 'ok b)))
                  (send flow :entry))`,
		Expect: "nil",
	}).Test(t)
}

func TestFlowEntryNotFound(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow-flavor :name 'flo)))
                  (flow-add-task flow :name "tisk" :actor (lambda (b) (list 'ok b)))
                  (flow-set-entry flow "who"))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}

func TestFlowEntryNotFlow(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-entry t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
