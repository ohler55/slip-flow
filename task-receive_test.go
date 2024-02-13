// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"bytes"
	"testing"

	"github.com/ohler55/ojg/tt"
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
		Source: `(let ((task (make-instance 'flow-task
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
		Source: `(let ((task (make-instance 'flow-task
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
		Source: `(let ((task (make-instance 'flow-task
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
		Source: `(let ((task (make-instance 'flow-task
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
		Source: `(let ((task (make-instance 'flow-task
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
		Source: `(let ((task (make-instance 'flow-task
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
		Source: `(let ((task (make-instance 'flow-task
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
		Source: `(let ((task (make-instance 'flow-task
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
                          (unless (string= "flow-task" (send (send task :flavor) :name))
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
		Source: `(let ((task (make-instance 'flow-task :name 'tisk)))
                  (flow-task-receive task))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}

func TestTaskReceiveNotBox(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task :name 'tisk)))
                  (flow-task-receive task t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((task (make-instance 'flow-task :name 'tisk)))
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
		Source: `(make-instance 'flow-task
                                            :name "tisk"
                                            :actor (make-instance 'vanilla-flavor))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
	(&sliptest.Function{
		Scope: scope,
		Source: `(make-instance 'flow-task
                                            :name "tisk"
                                            :actor (list
                                                    (make-instance 'vanilla-flavor)
                                                    (make-instance 'vanilla-flavor)))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskReceiveNoErrorTask(t *testing.T) {
	testTaskReceive(t, `
(let* ((lg (make-instance 'logger-flavor))
       (flow (make-instance 'flow :name 'flo :logger lg)))
  (flow-add-task flow
                 :name "start"
                 :actor (lambda (b) (list 'ok)))
  (flow-set-entry flow 'start)
  (flow-submit flow (make-flow-box :set '(1)))
  (send lg :shutdown))`,
		"/^E flo:start box .+: /")
}

func TestTaskReceiveErrorTask(t *testing.T) {
	testTaskReceive(t, `
(let* ((lg (make-instance 'logger-flavor))
       (flow (make-instance 'flow :name 'flo :logger lg)))
  (flow-add-task flow
                 :name "start"
                 :actor (lambda (b) (list 'ok)))
  (flow-add-task flow
                 :name "error"
                 :actor (lambda (b) (format t "~A~%" (send b :get "error")) (list nil nil)))
  (flow-set-entry flow 'start)
  (flow-submit flow (make-flow-box :set '(1)))
  (send lg :shutdown))`,
		"Actor did not return a list of link name and box instance.\n")
}

func TestTaskReceiveErrorLink(t *testing.T) {
	testTaskReceive(t, `
(let* ((lg (make-instance 'logger-flavor))
       (flow (make-instance 'flow :name 'flo :logger lg)))
  (flow-add-task flow
                 :name "start"
                 :actor (lambda (b) (list 'ok)))
  (flow-add-task flow
                 :name "bad"
                 :actor (lambda (b) (format t "~A~%" (send b :get "error")) (list nil nil)))
  (flow-link flow 'error 'start "bad")
  (flow-set-entry flow 'start)
  (flow-submit flow (make-flow-box :set '(1)))
  (send lg :shutdown))`,
		"Actor did not return a list of link name and box instance.\n")
}

func TestTaskReceiveLogInfo(t *testing.T) {
	testTaskReceive(t, `
(let* ((lg (make-instance 'logger-flavor))
       (flow (make-instance 'flow :name 'flo :logger lg)))
  (flow-add-task flow
                 :name "start"
                 :actor (lambda (b) (list "ok" b)))
  (flow-add-task flow
                 :name "done"
                 :actor (make-instance 'flow-exit-actor))
  (flow-link flow 'ok 'start 'done)
  (flow-set-entry flow 'start)
  (send flow :set-level 'info)
  (flow-submit flow (make-flow-box :set '(1)))
  (send lg :shutdown))`,
		"/^I flo:start received box /")
}

func TestTaskReceiveEmptyLinkName(t *testing.T) {
	testTaskReceive(t, `
(let* ((lg (make-instance 'logger-flavor))
       (flow (make-instance 'flow :name 'flo :logger lg)))
  (flow-add-task flow
                 :name "start"
                 :actor (lambda (b) (list nil b)))
  (flow-add-task flow
                 :name "done"
                 :actor (make-instance 'flow-exit-actor))
  (flow-link flow "" 'start 'done)
  (flow-set-entry flow 'start)
  (send flow :set-level 'info)
  (flow-submit flow (make-flow-box :set '(1)))
  (send lg :shutdown))`,
		"/^I flo:start received box /")
}

func testTaskReceive(t *testing.T, code, expect string) {
	var b bytes.Buffer
	scope := slip.NewScope()
	orig := scope.Get("*standard-output*")
	defer scope.Set("*standard-output*", orig)
	scope.Set("*standard-output*", &slip.OutputStream{Writer: &b})

	(&sliptest.Function{
		Scope:  scope,
		Source: code,
		Expect: "nil",
	}).Test(t)
	tt.Equal(t, expect, b.String())
}
