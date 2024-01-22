// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"strings"
	"testing"

	"github.com/ohler55/ojg/tt"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/gi"
	"github.com/ohler55/slip/sliptest"
)

func TestTaskTransitionFunction(t *testing.T) {
	exitChan := make(gi.Channel, 5)
	scope := slip.NewScope()
	scope.Let("exit-channel", exitChan)

	_ = slip.ReadString(`
(defflavor task-transition-test-actor () (flow-task-actor))
(defmethod (task-transition-test-actor :perform) (box)
 (flow-task-transition task box 'one)
 (flow-task-transition task box "two")
 (list nil nil))
`).Eval(scope, nil)

	(&sliptest.Function{
		Scope: scope,
		Source: `(let ((flow (make-flow :name 'flo :exit-channel exit-channel)))
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
  (flow-submit flow (make-flow-box :set '(1))))`,
		Expect: "nil",
	}).Test(t)

	out := <-exitChan
	scope.Let("split-out", out)
	history := slip.ReadString(
		`(mapcar (lambda (ev) (cadr ev))(send (send split-out :track) :history))`).Eval(scope, nil).(slip.List)

	out = <-exitChan
	scope.Let("split-out", out)
	history = append(history,
		slip.ReadString(
			`(mapcar (lambda (ev) (cadr ev))(send (send split-out :track) :history))`).Eval(scope, nil).(slip.List)...,
	)
	hstr := slip.ObjectString(history)
	tt.Equal(t, true, strings.Contains(hstr, "branch-one"))
	tt.Equal(t, true, strings.Contains(hstr, "branch-two"))
}

func TestTaskTransitionSend(t *testing.T) {
	exitChan := make(gi.Channel, 5)
	scope := slip.NewScope()
	scope.Let("exit-channel", exitChan)

	_ = slip.ReadString(`
(defflavor task-transition-test2-actor () (flow-task-actor))
(defmethod (task-transition-test2-actor :perform) (box)
 (send task :transition box 'one)
 (send task :transition box "two")
 (list nil nil))
`).Eval(scope, nil)

	(&sliptest.Function{
		Scope: scope,
		Source: `(let ((flow (make-flow :name 'flo :exit-channel exit-channel)))
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
  (flow-submit flow (make-flow-box :set '(1))))`,
		Expect: "nil",
	}).Test(t)

	out := <-exitChan
	scope.Let("split-out", out)
	history := slip.ReadString(
		`(mapcar (lambda (ev) (cadr ev))(send (send split-out :track) :history))`).Eval(scope, nil).(slip.List)

	out = <-exitChan
	scope.Let("split-out", out)
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
		Source:    `(flow-task-transition t 'ok (make-instance 'flow-box-flavor))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskTransitionNotBox(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-task-transition (make-instance 'flow-task-flavor) t 'ok)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
	(&sliptest.Function{
		Source:    `(send (make-instance 'flow-task-flavor) :transition t 'ok)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskTransitionNotLink(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-task-transition (make-instance 'flow-task-flavor) (make-instance 'flow-box-flavor) t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
	(&sliptest.Function{
		Source:    `(send (make-instance 'flow-task-flavor) :transition (make-instance 'flow-box-flavor) t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
