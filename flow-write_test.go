// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"bytes"
	"testing"

	"github.com/ohler55/ojg/tt"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestFlowWriteSend(t *testing.T) {
	scope := slip.NewScope()
	scope.Let("*print-right-margin*", slip.Fixnum(80))
	(&sliptest.Function{
		Scope: scope,
		Source: `
(let ((flow (make-flow :name 'flo :width 300 :height 280 :task-width 40 :task-height 40)))
  (flow-add-task flow
                 :name "start"
                 :x 100
                 :y 0
                 :actor (lambda (b)
                          (flow-box-set b (* 3 (flow-box-get b "[0]")) "[0]")
                          (list 'ok b)))
  (flow-add-task flow
                 :name "odd-or-even"
                 :x 100
                 :y 100
                 :svg "<svg></svg>"
                 :workers 2
                 :depth 5
                 :actor (lambda (b)
                          (cond ((evenp (flow-box-get b "[0]")) (list 'even b))
                                (t (list 'odd b)))))
  (flow-add-task flow
                 :name "even"
                 :x 0
                 :y 200
                 :actor (make-instance 'flow-exit-actor))
  (flow-add-task flow
                 :name "odd"
                 :x 200
                 :y 200
                 :actor (make-instance 'flow-exit-actor))
  (flow-link flow 'ok 'start "odd-or-even")
  (flow-link flow 'odd "odd-or-even" 'odd '((220 50)))
  (flow-link flow 'even "odd-or-even" 'even '((20 50)))
  (flow-set-entry flow 'start)
  (send flow :write nil))`,
		Expect: `"(let ((flow (make-flow :name "flo"
                       :width 300
                       :height 280
                       :task-width 40
                       :task-height 40)))
  (send flow :add-task
        :name "even"
        :x 0
        :y 200
        :actor (make-instance 'flow-exit-actor))
  (send flow :add-task
        :name "odd"
        :x 200
        :y 200
        :actor (make-instance 'flow-exit-actor))
  (send flow :add-task
        :name "odd-or-even"
        :x 100
        :y 100
        :svg "<svg></svg>"
        :workers 2
        :depth 5
        :actor (lambda (b)
                       (cond ((evenp (flow-box-get b "[0]")) (list 'even b))
                             (t (list 'odd b)))))
  (send flow :add-task
        :name "start"
        :x 100
        :y 0
        :actor (lambda (b) (flow-box-set b (* 3 (flow-box-get b "[0]")) "[0]")
                       (list 'ok b)))
  (send flow :link "even" "odd-or-even" "even" '((20 50)))
  (send flow :link "odd" "odd-or-even" "odd" '((220 50)))
  (send flow :link "ok" "start" "odd-or-even")
  (send flow :set-entry "start")
  flow)
"`,
	}).Test(t)
}

func TestFlowWriteFunction(t *testing.T) {
	scope := slip.NewScope()

	var b bytes.Buffer
	orig := scope.Get("*standard-output*")
	defer scope.Set("*standard-output*", orig)
	scope.Set("*standard-output*", &slip.OutputStream{Writer: &b})
	scope.Let("*print-right-margin*", slip.Fixnum(80))

	_ = slip.ReadString("(defun flow-write-test-perform (b) (list 'ok b))").Eval(scope, nil)

	(&sliptest.Function{
		Scope: scope,
		Source: `
(let ((flow (make-flow :name 'flo)))
  (flow-add-task flow
                 :name "start"
                 :x 100
                 :y 0
                 :svg "<svg></svg>"
                 :actor 'flow-write-test-perform)
  (flow-add-task flow
                 :name "done"
                 :x 200
                 :y 0
                 :actor (list (make-instance 'flow-exit-actor) (make-instance 'flow-exit-actor)))
  (flow-link flow 'ok 'start "done")
  (flow-set-entry flow 'start)
  (flow-write flow t t))`,
		Expect: "nil",
	}).Test(t)

	tt.Equal(t, `(let ((flow (make-flow :name "flo")))
  (flow-add-task flow
                 :name "done"
                 :x 200
                 :y 0
                 :actor (list
                         (make-instance 'flow-exit-actor)
                         (make-instance 'flow-exit-actor)))
  (flow-add-task flow
                 :name "start"
                 :x 100
                 :y 0
                 :svg "<svg></svg>"
                 :actor 'flow-write-test-perform)
  (flow-link flow "ok" "start" "done")
  (flow-set-entry flow "start")
  flow)
`, b.String())

}

func TestFlowWriteStream(t *testing.T) {
	scope := slip.NewScope()

	var b bytes.Buffer
	scope.Set("out", &slip.OutputStream{Writer: &b})
	scope.Let("*print-right-margin*", slip.Fixnum(80))

	(&sliptest.Function{
		Scope: scope,
		Source: `
(let ((flow (make-flow :name 'flo)))
  (flow-add-task flow
                 :name "start"
                 :x 100
                 :y 0
                 :svg "<svg></svg>"
                 :actor (make-instance 'flow-http-client-actor
                                       :method 'get
                                       :timeout (lambda (b) 1)
                                       :header '((Accept . "text/html"))
                                       :url "http://localhost:7777"))
  (flow-add-task flow
                 :name "done"
                 :x 200
                 :y 0
                 :actor (list (make-instance 'flow-exit-actor) (make-instance 'flow-exit-actor)))
  (flow-link flow 'ok 'start "done")
  (flow-set-entry flow 'start)
  (flow-write flow out t))`,
		Expect: "nil",
	}).Test(t)

	tt.Equal(t, `(let ((flow (make-flow :name "flo")))
  (flow-add-task flow
                 :name "done"
                 :x 200
                 :y 0
                 :actor (list
                         (make-instance 'flow-exit-actor)
                         (make-instance 'flow-exit-actor)))
  (flow-add-task flow
                 :name "start"
                 :x 100
                 :y 0
                 :svg "<svg></svg>"
                 :actor (make-instance 'flow-http-client-actor
                                       :method "get"
                                       :url "http://localhost:7777"
                                       :timeout (lambda (b) 1))
                                       :header '(("Accept" "text/html"))))
  (flow-link flow "ok" "start" "done")
  (flow-set-entry flow "start")
  flow)
`, b.String())
}

func TestFlowWriteBadFlow(t *testing.T) {
	(&sliptest.Function{
		Source:    "(flow-write t)",
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowWriteBadStream(t *testing.T) {
	(&sliptest.Function{
		Source:    "(flow-write (make-flow :name 'flo) 7)",
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowWriteOutputError(t *testing.T) {
	scope := slip.NewScope()
	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: badWriter(0)})
	(&sliptest.Function{
		Scope:     scope,
		Source:    `(send (make-flow) :write out)`,
		PanicType: slip.Symbol("stream-error"),
	}).Test(t)
}
