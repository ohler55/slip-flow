// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/ojg/tt"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/gi"
	"github.com/ohler55/slip/sliptest"
)

func TestFlowSubmitFunction(t *testing.T) {
	exitChan := make(gi.Channel, 5)
	scope := slip.NewScope()
	scope.Let("exit-channel", exitChan)

	(&sliptest.Function{
		Scope: scope,
		Source: `
(let ((flow (make-instance 'flow-flavor :name 'flo :exit-channel exit-channel)))
  (flow-add-task flow
                 :name "start"
                 :actor (lambda (b)
                          (flow-box-set b (* 3 (flow-box-get b "[0]")) "[0]")
                          (list 'ok b)))
  (flow-add-task flow
                 :name "odd-or-even"
                 :actor (lambda (b)
                          (cond ((evenp (flow-box-get b "[0]")) (list 'even b))
                                (t (list 'odd b)))))
  (flow-add-task flow
                 :name "even"
                 :actor (make-instance 'flow-exit-actor))
  (flow-add-task flow
                 :name "odd"
                 :actor (make-instance 'flow-exit-actor))
  (flow-link flow 'ok 'start "odd-or-even")
  (flow-link flow 'odd "odd-or-even" 'odd)
  (flow-link flow 'even "odd-or-even" 'even)
  (flow-set-entry flow 'start)
  (flow-start flow)
  (send flow :set-level 'warn)
  (flow-submit flow (make-flow-box :set '(1))))`,
		Expect: "nil",
	}).Test(t)

	out := <-exitChan
	scope.Let("submit-test-out", out)

	history := slip.ReadString(
		`(mapcar (lambda (ev) (cadr ev))(send (send submit-test-out :track) :history))`).Eval(scope, nil)
	tt.Equal(t, `("start" "odd-or-even" "odd")`, slip.ObjectString(history))

	value := slip.ReadString(`(send submit-test-out :native)`).Eval(scope, nil)
	tt.Equal(t, `(3)`, slip.ObjectString(value))
}

func TestFlowSubmitSend(t *testing.T) {
	exitChan := make(gi.Channel, 5)
	scope := slip.NewScope()
	scope.Let("exit-channel", exitChan)

	(&sliptest.Function{
		Scope: scope,
		Source: `
(let ((flow (make-instance 'flow-flavor :name 'flo :exit-channel exit-channel)))
  (flow-add-task flow
                 :name "start"
                 :actor (lambda (b)
                          (flow-box-set b (* 3 (flow-box-get b "[0]")) "[0]")
                          (list 'ok b)))
  (flow-add-task flow
                 :name "odd-or-even"
                 :actor (lambda (b)
                          (cond ((evenp (flow-box-get b "[0]")) (list 'even b))
                                (t (list 'odd b)))))
  (flow-add-task flow
                 :name "even"
                 :actor (make-instance 'flow-exit-actor))
  (flow-add-task flow
                 :name "odd"
                 :actor (make-instance 'flow-exit-actor))
  (flow-link flow 'ok 'start "odd-or-even")
  (flow-link flow 'odd "odd-or-even" 'odd)
  (flow-link flow 'even "odd-or-even" 'even)
  (flow-set-entry flow 'start)
  (flow-start flow)
  (send flow :set-level 'warn)
  (send flow :submit (make-flow-box :set '(1))))`,
		Expect: "nil",
	}).Test(t)

	out := <-exitChan
	scope.Let("submit-test-out", out)

	history := slip.ReadString(
		`(mapcar (lambda (ev) (cadr ev))(send (send submit-test-out :track) :history))`).Eval(scope, nil)
	tt.Equal(t, `("start" "odd-or-even" "odd")`, slip.ObjectString(history))

	value := slip.ReadString(`(send submit-test-out :native)`).Eval(scope, nil)
	tt.Equal(t, `(3)`, slip.ObjectString(value))
}

func TestFlowSubmitNotFlow(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-submit t (make-flow-box))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
