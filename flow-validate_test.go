// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestFlowValidateFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((flow (make-flow :name 'flo)))
  (flow-add-task flow
                 :name "start"
                 :actor (lambda (b) (list 'ok b)))
  (flow-add-task flow
                 :name "inspect"
                 :actor (make-instance 'flow-inspect-actor))
  (flow-add-task flow
                 :name "done"
                 :actor (make-instance 'flow-exit-actor))
  (flow-link flow 'ok 'start "inspect")
  (flow-link flow 'ok "inspect" 'done)
  (flow-set-entry flow 'start)
  (flow-validate flow))`,
		Expect: "nil",
	}).Test(t)
}

func TestFlowValidateSend(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((flow (make-flow :name 'flo)))
  (send flow :add-task
             :name "start"
             :actor (lambda (b) (list 'ok b)))
  (send flow :add-task
             :name "inspect"
             :actor (make-instance 'flow-inspect-actor))
  (send flow :add-task
             :name "done"
             :actor (make-instance 'flow-exit-actor))
  (send flow :link 'ok 'start "inspect")
  (send flow :link 'bad "inspect" 'done)
  (send flow :validate))`,
		Expect: `("no entry task" "start is not reachable"
                 "task inspect can not return a bad link")`,
	}).Test(t)
}

func TestFlowValidateNotFlow(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-validate t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
