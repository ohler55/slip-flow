// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

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
	(&sliptest.Function{
		Source: `
(let ((flow (make-flow :name 'flo)))
  (flow-add-task flow
                 :name "start"
                 :x 100
                 :y 0
                 :svg "<svg></svg>"
                 :actor (lambda (b) (list 'ok b)))
  (flow-add-task flow
                 :name "done"
                 :x 200
                 :y 0
                 :actor (make-instance 'flow-exit-actor))
  (flow-link flow 'ok 'start "done")
  (flow-set-entry flow 'start)
  (flow-write flow nil t))`,
		Expect: `"(let ((flow (make-flow :name "flo")))
  (flow-add-task flow
                 :name "done"
                 :x 200
                 :y 0
                 :actor (make-instance 'flow-exit-actor))
  (flow-add-task flow
                 :name "start"
                 :x 100
                 :y 0
                 :svg "<svg></svg>"
                 :actor (lambda (b) (list 'ok b)))
  (flow-link flow "ok" "start" "done")
  (flow-set-entry flow "start")
  flow)
"`,
	}).Test(t)
}

func TestFlowWriteBadX(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((flow (make-flow :name 'flo)))
  (flow-add-task flow :name "start" :x t :y 0 :actor (lambda (b) (list 'ok b))))
`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowWriteBadY(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((flow (make-flow :name 'flo)))
  (flow-add-task flow :name "start" :x 0 :y t :actor (lambda (b) (list 'ok b))))
`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowWriteSvg(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((flow (make-flow :name 'flo)))
  (flow-add-task flow :name "start" :svg "<sgv></svg>" :actor (lambda (b) (list 'ok b)))
  flow)
`,
		Expect: "/#<flow-flavor [0-9a-f]+>/",
	}).Test(t)
	(&sliptest.Function{
		Source: `
(let ((flow (make-flow :name 'flo)))
  (flow-add-task flow :name "start" :svg t :actor (lambda (b) (list 'ok b))))
`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowWriteBadFlow(t *testing.T) {
	(&sliptest.Function{
		Source:    "(flow-write t)",
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
