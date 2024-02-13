// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestGroupFlowsFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((group (make-flow-group))
                       (flow (make-instance 'flow :name 'flo)))
                  (flow-add-task flow :name "tisk" :workers 1 :actor (lambda (b) (list 'ok b)))
                  (flow-group-add group flow)
                  (flow-group-flows group))`,
		Expect: `/\(#<flow [0-9a-f]+>\)/`,
	}).Test(t)
}

func TestGroupFlowsSend(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((group (make-instance 'flow-group))
                       (flow (make-instance 'flow :name 'flo)))
                  (send flow :add-task :name "tisk" :workers 1 :actor (lambda (b) (list 'ok b)))
                  (send group :add flow)
                  (send group :flows))`,
		Expect: `/\(#<flow [0-9a-f]+>\)/`,
	}).Test(t)
}

func TestGroupFlowsNotGroup(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-group-flows t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
