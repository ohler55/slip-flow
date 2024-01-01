// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
	"github.com/ohler55/slip/sliptest"
)

// Verify the actor is called but not that the result as that is checked with
// multiple tasks in the flow tests.

func TestTaskReceiveSync(t *testing.T) {
	scope := slip.NewScope()
	_ = slip.ReadString(
		`(defvar task-receive-test-box (make-flow-box :tracking-id 123 :set '((a . 1)(b . 2))))`).Eval(scope, nil)
	(&sliptest.Function{
		Scope: scope,
		Source: `(let ((task (make-instance 'flow-task-flavor
                                            :name "tisk"
                                            :actor (lambda (b)
                                                    (flow-box-set b 3 "c")
                                                    (setq task-receive-test-box b)
                                                    (list 'ok b)))))
                  (flow-task-receive task task-receive-test-box)
                  (flow-box-native task-receive-test-box))`,
		Expect: `/\("c" \. 3\)/`,
	}).Test(t)
	_ = slip.ReadString(
		`(setq task-receive-test-box (make-flow-box :tracking-id 123 :set '((a . 1)(b . 2))))`).Eval(scope, nil)
	(&sliptest.Function{
		Scope: scope,
		Source: `(let ((task (make-instance 'flow-task-flavor
                                            :name 'tisk
                                            :actor (lambda (b)
                                                    (flow-box-set b 3 "c")
                                                    (setq task-receive-test-box b)
                                                    (list 'ok b)))))
                  (send task :receive task-receive-test-box)
                  (send task-receive-test-box :native))`,
		Expect: `/\("c" \. 3\)/`,
	}).Test(t)
}

func TestTaskReceiveAsync(t *testing.T) {
	scope := slip.NewScope()
	_ = slip.ReadString(
		`(defvar task-receive-test-box (make-flow-box :tracking-id 123 :set '((a . 1)(b . 2))))`).Eval(scope, nil)
	(&sliptest.Function{
		Scope: scope,
		Source: `(let ((task (make-instance 'flow-task-flavor
                                            :name 'tisk
                                            :workers 3
                                            :depth 5
                                            :actor (lambda (b)
                                                    (flow-box-set b 3 "c")
                                                    (setq task-receive-test-box b)
                                                    (list 'ok b)))))
                  (flow-task-start task)
                  (flow-task-receive task task-receive-test-box)
                  (flow-task-shutdown task)
                  (flow-box-native task-receive-test-box))`,
		Expect: `/\("c" \. 3\)/`,
	}).Test(t)
	_ = slip.ReadString(
		`(setq task-receive-test-box (make-flow-box :tracking-id 123 :set '((a . 1)(b . 2))))`).Eval(scope, nil)
	(&sliptest.Function{
		Scope: scope,
		Source: `(let ((task (make-instance 'flow-task-flavor
                                            :name 'tisk
                                            :workers 3
                                            :depth 5
                                            :actor (lambda (b)
                                                    (flow-box-set b 3 "c")
                                                    (setq task-receive-test-box b)
                                                    (list 'ok b)))))
                  (send task :start)
                  (send task :receive task-receive-test-box)
                  (send task :shutdown)
                  (send task-receive-test-box :native))`,
		Expect: `/\("c" \. 3\)/`,
	}).Test(t)
}

func TestTaskReceiveFunction(t *testing.T) {
	scope := slip.NewScope()
	_ = slip.ReadString(
		`(defvar task-receive-test-box (make-flow-box :tracking-id 123 :set '((a . 1)(b . 2))))`).Eval(scope, nil)
	_ = slip.ReadString(`(defun task-receive-func (b)
                          (flow-box-set b 3 "c")
                          (setq task-receive-test-box b)
                          (list 'ok b))`).Eval(scope, nil)
	(&sliptest.Function{
		Scope: scope,
		Source: `(let ((task (make-instance 'flow-task-flavor
                                            :name "tisk"
                                            :actor 'task-receive-func)))
                  (flow-task-receive task task-receive-test-box)
                  (flow-box-native task-receive-test-box))`,
		Expect: `/\("c" \. 3\)/`,
	}).Test(t)
	_ = slip.ReadString(
		`(setq task-receive-test-box (make-flow-box :tracking-id 123 :set '((a . 1)(b . 2))))`).Eval(scope, nil)
	(&sliptest.Function{
		Scope: scope,
		Source: `(let ((task (make-instance 'flow-task-flavor
                                            :name 'tisk
                                            :actor 'task-receive-func)))
                  (send task :receive task-receive-test-box)
                  (send task-receive-test-box :native))`,
		Expect: `/\("c" \. 3\)/`,
	}).Test(t)
}

func TestTaskReceiveInstanceSync(t *testing.T) {
	scope := slip.NewScope()
	_ = slip.ReadString(
		`(defvar task-receive-test-box (make-flow-box :tracking-id 123 :set '((a . 1)(b . 2))))`).Eval(scope, nil)
	assureTaskRecieveTestActor(scope)
	(&sliptest.Function{
		Scope: scope,
		Source: `(let ((task (make-instance 'flow-task-flavor
                                            :name "tisk"
                                            :actor (make-instance 'task-receiver-test-actor))))
                  (send task :start)
                  (flow-task-receive task task-receive-test-box)
                  (send task :shutdown)
                  (flow-box-native task-receive-test-box))`,
		Expect: `/\("c" \. 3\)/`,
	}).Test(t)
}

func TestTaskReceiveInstanceAsync(t *testing.T) {
	scope := slip.NewScope()
	_ = slip.ReadString(
		`(defvar task-receive-test-box (make-flow-box :tracking-id 123 :set '((a . 1)(b . 2))))`).Eval(scope, nil)
	assureTaskRecieveTestActor(scope)

	(&sliptest.Function{
		Scope: scope,
		Source: `(let ((task (make-instance 'flow-task-flavor
                                            :name "tisk"
                                            :workers 3
                                            :actor (list
                                                    (make-instance 'task-receiver-test-actor)
                                                    (make-instance 'task-receiver-test-actor)
                                                    (make-instance 'task-receiver-test-actor)))))
                  (send task :start)
                  (flow-task-receive task task-receive-test-box)
                  (send task :shutdown)
                  (flow-box-native task-receive-test-box))`,
		Expect: `/\("c" \. 3\)/`,
	}).Test(t)
}

func assureTaskRecieveTestActor(scope *slip.Scope) {
	if flavors.Find("task-receiver-test-actor") == nil {
		_ = slip.ReadString(`(defflavor task-receiver-test-actor (task)
                                                             ()
                                                             :gettable-instance-variables
                                                             :settable-instance-variables)`).Eval(scope, nil)
		_ = slip.ReadString(`(defmethod (task-receiver-test-actor :start) (tsk) (setq task tsk))`).Eval(scope, nil)
		_ = slip.ReadString(`(defmethod (task-receiver-test-actor :shutdown) ()
                          (unless (string= "flow-task-flavor" (send (send task :flavor) :name))
                                  (panic "task not set")))`).Eval(scope, nil)
		_ = slip.ReadString(`(defmethod (task-receiver-test-actor :perform) (b)
                                     (flow-box-set b 3 "c")
                                     (setq task-receive-test-box b)
                                     (list 'ok b))`).Eval(scope, nil)
	}
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

func TestTaskReceiveNotBox(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task-flavor :name 'tisk)))
                  (flow-task-receive task t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task-flavor :name 'tisk)))
                  (send task :receive t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskReceiveInstanceNoPerform(t *testing.T) {
	scope := slip.NewScope()
	_ = slip.ReadString(
		`(defvar task-receive-test-box (make-flow-box :tracking-id 123 :set '((a . 1)(b . 2))))`).Eval(scope, nil)
	(&sliptest.Function{
		Scope: scope,
		Source: `(make-instance 'flow-task-flavor
                                            :name "tisk"
                                            :actor (make-instance 'vanilla-flavor))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
	(&sliptest.Function{
		Scope: scope,
		Source: `(make-instance 'flow-task-flavor
                                            :name "tisk"
                                            :actor (list
                                                    (make-instance 'vanilla-flavor)
                                                    (make-instance 'vanilla-flavor)))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
