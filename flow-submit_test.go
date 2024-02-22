// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/ojg/tt"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestFlowSubmitFunction(t *testing.T) {
	scope := slip.NewScope()
	tf := sliptest.Function{
		Scope: scope,
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
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
                 :actor (make-instance 'flow-exit-actor :notifiers 'done))
  (flow-add-task flow
                 :name "odd"
                 :actor (make-instance 'flow-exit-actor :notifiers "done"))
  (flow-link flow 'ok 'start "odd-or-even")
  (flow-link flow 'odd "odd-or-even" 'odd)
  (flow-link flow 'even "odd-or-even" 'even)
  (flow-set-entry flow 'start)
  (send flow :set-level 'warn)
  (flow-submit flow (make-flow-box :set '(1) :watch 'done))
  (channel-pop done))`,
		Expect: "/#<flow-box [0-9a-f]+>/",
	}
	tf.Test(t)
	scope.Let("submit-test-out", tf.Result)

	history := slip.ReadString(
		`(mapcar (lambda (ev) (cadr ev))(send (send submit-test-out :track) :history))`).Eval(scope, nil)
	tt.Equal(t, `("start" "odd-or-even" "odd")`, slip.ObjectString(history))

	value := slip.ReadString(`(send submit-test-out :native)`).Eval(scope, nil)
	tt.Equal(t, `(3)`, slip.ObjectString(value))
}

func TestFlowSubmitBag(t *testing.T) {
	scope := slip.NewScope()
	tf := sliptest.Function{
		Scope: scope,
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
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
                 :actor (make-instance 'flow-exit-actor :notifiers '(done)))
  (flow-link flow 'ok 'start "odd-or-even")
  (flow-link flow 'odd "odd-or-even" 'odd)
  (flow-link flow 'even "odd-or-even" 'even)
  (flow-set-entry flow 'start)
  (send flow :set-level 'warn)
  (flow-submit flow (make-bag '(1)) 'done)
  (channel-pop done))`,
		Expect: "/#<flow-box [0-9a-f]+>/",
	}
	tf.Test(t)
	scope.Let("submit-test-out", tf.Result)

	history := slip.ReadString(
		`(mapcar (lambda (ev) (cadr ev))(send (send submit-test-out :track) :history))`).Eval(scope, nil)
	tt.Equal(t, `("start" "odd-or-even" "odd")`, slip.ObjectString(history))

	value := slip.ReadString(`(send submit-test-out :native)`).Eval(scope, nil)
	tt.Equal(t, `(3)`, slip.ObjectString(value))
}

func TestFlowSubmitSend(t *testing.T) {
	scope := slip.NewScope()
	tf := sliptest.Function{
		Scope: scope,
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
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
  (send flow :submit '(1) 'done)
  (channel-pop done))`,
		Expect: "/#<flow-box [0-9a-f]+>/",
	}
	tf.Test(t)
	scope.Let("submit-test-out", tf.Result)

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

func TestFlowSubmitNoEntry(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((flow (make-instance 'flow :name 'flo)))
  (flow-add-task flow
                 :name "start"
                 :actor (lambda (b)
                          (flow-box-set b (* 3 (flow-box-get b "[0]")) "[0]")
                          (list 'ok b)))
  (flow-submit flow (make-flow-box :set '(1))))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}

func TestFlowSubmitNotBox(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((flow (make-instance 'flow :name 'flo)))
  (flow-add-task flow
                 :name "start"
                 :actor (lambda (b)
                          (flow-box-set b (* 3 (flow-box-get b "[0]")) "[0]")
                          (list 'ok b)))
  (flow-set-entry flow 'start)
  (flow-submit flow (make-instance 'vanilla-flavor)))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowSubmitBadWatch(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((flow (make-instance 'flow :name 'flo)))
  (flow-add-task flow
                 :name "start"
                 :actor (lambda (b)
                          (flow-box-set b (* 3 (flow-box-get b "[0]")) "[0]")
                          (list 'ok b)))
  (flow-set-entry flow 'start)
  (flow-submit flow (make-flow-box :set '(1)) t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
