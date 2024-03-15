// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestGroupFindFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((group (make-flow-group))
                       (flow (make-instance 'flow :name 'flo)))
                  (flow-add-task flow :name "tisk" :workers 1 :actor (lambda (b) (list 'ok b)))
                  (flow-group-add group flow)
                  (send group :set-level 'warn)
                  (flow-group-find group 'flo))`,
		Expect: "/#<flow [0-9a-f]+>/",
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((group (make-flow-group))
                       (flow (make-instance 'flow :name 'flo)))
                  (flow-add-task flow :name "tisk" :workers 1 :actor (lambda (b) (list 'ok b)))
                  (flow-group-add group flow)
                  (flow-group-find group 'flu))`,
		Expect: "nil",
	}).Test(t)
}

func TestGroupFindSend(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((group (make-instance 'flow-group))
                       (flow (make-instance 'flow :name 'flo)))
                  (send flow :add-task :name "tisk" :workers 1 :actor (lambda (b) (list 'ok b)))
                  (send group :add flow)
                  (send group :find "flo"))`,
		Expect: "/#<flow [0-9a-f]+>/",
	}).Test(t)
}

func TestGroupFindNotGroup(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-group-find t "flo")`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestGroupFindArgCount(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-group-find (make-instance 'flow-group))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}

func TestGroupFindBadName(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-group-find (make-flow-group) t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
