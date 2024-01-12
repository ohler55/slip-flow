// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

// Partially tested in flow-start_test.go.

func TestFlowRunningFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow-flavor :name 'flo)))
                  (flow-add-task flow :name "tisk" :workers 1 :actor (lambda (b) (list 'ok b)))
                  (flow-running flow))`,
		Expect: "nil",
	}).Test(t)
}

func TestFlowRunningSend(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow-flavor :name 'flo)))
                  (send flow :add-task :name "tisk" :workers 1 :actor (lambda (b) (list 'ok b)))
                  (send flow :running))`,
		Expect: "nil",
	}).Test(t)
}

func TestFlowRunningNotFlow(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-running t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowRunningArgCount(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-running (make-instance 'flow-flavor :name 'flo) t)`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}
