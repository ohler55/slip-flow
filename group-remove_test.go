// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestGroupRemoveFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((group (make-flow-group))
                       (flow (make-instance 'flow-flavor :name 'flo)))
                  (flow-add-task flow :name "tisk" :workers 1 :actor (lambda (b) (list 'ok b)))
                  (flow-group-add group flow)
                  (flow-group-remove group 'flo)
                  (flow-group-flows group))`,
		Expect: "nil",
	}).Test(t)
}

func TestGroupRemoveSend(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((group (make-instance 'flow-group-flavor))
                       (flow (make-instance 'flow-flavor :name 'flo)))
                  (send flow :add-task :name "tisk" :workers 1 :actor (lambda (b) (list 'ok b)))
                  (send group :add flow)
                  (send group :remove "flo")
                  (send group :flows))`,
		Expect: "nil",
	}).Test(t)
}

func TestGroupRemoveNotGroup(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-group-remove t "flo")`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestGroupRemoveBadString(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-group-remove (make-flow-group) t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
