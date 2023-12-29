// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

// Verify the actor is called but not that the result as that is checked with
// multiple tasks in the flow tests.

func TestTaskReceiveOk(t *testing.T) {
	scope := slip.NewScope()
	_ = slip.ReadString(
		`(defvar task-receive-test-box (make-flow-box :tracking-id 123 :set '((a . 1)(b . 2))))`).Eval(scope, nil)
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task-flavor
                                            :name 'tisk
                                            :actor (lambda (b)
                                                    (flow-box-set b 3 "c")
                                                    (setq task-receive-test-box b)
                                                    (list 'ok b)))))
                  (flow-task-receive task task-receive-test-box)
                  (flow-box-native task-receive-test-box))`,
		Expect: `/\("c" \. 3\)/`,
	}).Test(t)
}

func TestTaskReceiveNotTask(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-task-receive t (make-flow-box))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskReceiveArgCount(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task-flavor :name 'tisk)))
                  (flow-task-receive task))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}

// TBD
