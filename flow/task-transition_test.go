// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"strings"
	"testing"

	"github.com/ohler55/ojg/tt"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestTaskTransitionFunction(t *testing.T) {
	scope := slip.NewScope()
	_ = slip.ReadString(`
(defflavor task-transition-test-actor () (flow-task-actor))
(defmethod (task-transition-test-actor :perform) (box)
 (flow-task-transition task box 'one)
 (flow-task-transition task box "two")
 (list nil nil))
`).Eval(scope, nil)

	tf := sliptest.Function{
		Scope: scope,
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (flow-add-task flow
                :name "start"
                :actor (make-instance 'task-transition-test-actor))
 (flow-add-task flow
                :name "branch-one"
                :actor (make-instance 'flow-exit-actor))
 (flow-add-task flow
                :name "branch-two"
                :actor (make-instance 'flow-exit-actor))

 (flow-link flow 'one 'start "branch-one")
 (flow-link flow 'two 'start "branch-two")
 (flow-set-entry flow 'start)
 (flow-submit flow (make-flow-box :set '(1) :watch 'done))
 (list (channel-pop done) (channel-pop done)))`,
		Expect: `/\(#<flow-box [0-9a-f]+> #<flow-box [0-9a-f]+>\)/`,
	}
	tf.Test(t)
	scope.Let("split-out", tf.Result.(slip.List)[0])
	history := slip.ReadString(
		`(mapcar (lambda (ev) (cadr ev))(send (send split-out :track) :history))`).Eval(scope, nil).(slip.List)

	scope.Let("split-out", tf.Result.(slip.List)[1])
	history = append(history,
		slip.ReadString(
			`(mapcar (lambda (ev) (cadr ev))(send (send split-out :track) :history))`).Eval(scope, nil).(slip.List)...,
	)
	hstr := slip.ObjectString(history)
	tt.Equal(t, true, strings.Contains(hstr, "branch-one"))
	tt.Equal(t, true, strings.Contains(hstr, "branch-two"))
}

func TestTaskTransitionSend(t *testing.T) {
	scope := slip.NewScope()
	_ = slip.ReadString(`
(defflavor task-transition-test2-actor () (flow-task-actor))
(defmethod (task-transition-test2-actor :perform) (box)
 (send task :transition box 'one)
 (send task :transition box "two")
 (list nil nil))
`).Eval(scope, nil)

	tf := sliptest.Function{
		Scope: scope,
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (flow-add-task flow
                :name "start"
                :actor (make-instance 'task-transition-test2-actor))
 (flow-add-task flow
                :name "branch-one"
                :actor (make-instance 'flow-exit-actor))
 (flow-add-task flow
                :name "branch-two"
                :actor (make-instance 'flow-exit-actor))

 (flow-link flow 'one 'start "branch-one")
 (flow-link flow 'two 'start "branch-two")
 (flow-set-entry flow 'start)
 (flow-submit flow (make-flow-box :set '(1) :watch 'done))
 (list (channel-pop done) (channel-pop done)))`,
		Expect: `/\(#<flow-box [0-9a-f]+> #<flow-box [0-9a-f]+>\)/`,
	}
	tf.Test(t)
	scope.Let("split-out", tf.Result.(slip.List)[0])
	history := slip.ReadString(
		`(mapcar (lambda (ev) (cadr ev))(send (send split-out :track) :history))`).Eval(scope, nil).(slip.List)

	scope.Let("split-out", tf.Result.(slip.List)[1])
	history = append(history,
		slip.ReadString(
			`(mapcar (lambda (ev) (cadr ev))(send (send split-out :track) :history))`).Eval(scope, nil).(slip.List)...,
	)
	hstr := slip.ObjectString(history)
	tt.Equal(t, true, strings.Contains(hstr, "branch-one"))
	tt.Equal(t, true, strings.Contains(hstr, "branch-two"))
}

func TestTaskTransitionNotTask(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-task-transition t 'ok (make-instance 'flow-box))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskTransitionNotBox(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-task-transition (make-instance 'flow-task) t 'ok)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
	(&sliptest.Function{
		Source:    `(send (make-instance 'flow-task) :transition t 'ok)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskTransitionNotLink(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-task-transition (make-instance 'flow-task) (make-instance 'flow-box) t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
	(&sliptest.Function{
		Source:    `(send (make-instance 'flow-task) :transition (make-instance 'flow-box) t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
