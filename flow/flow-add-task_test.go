// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

// Partially tested in flow-start_test.go.

func TestFlowAddTaskNotFlow(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-add-task t :name "tisk")`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowAddTaskExists(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow :name 'flo)))
                  (flow-add-task flow :name "tisk" :workers 1 :actor (lambda (b) (list 'ok b)))
                  (flow-add-task flow :name "tisk" :workers 2 :actor (lambda (b) (list 'ok b))))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}

func TestFlowAddTaskBadX(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((flow (make-flow :name 'flo)))
  (flow-add-task flow :name "start" :x t :y 0 :actor (lambda (b) (list 'ok b))))
`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowAddTaskBadY(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((flow (make-flow :name 'flo)))
  (flow-add-task flow :name "start" :x 0 :y t :actor (lambda (b) (list 'ok b))))
`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowAddTaskSvg(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((flow (make-flow :name 'flo)))
  (flow-add-task flow :name "start" :svg "<sgv></svg>" :actor (lambda (b) (list 'ok b)))
  flow)
`,
		Expect: "/#<flow [0-9a-f]+>/",
	}).Test(t)
	(&sliptest.Function{
		Source: `
(let ((flow (make-flow :name 'flo)))
  (flow-add-task flow :name "start" :svg t :actor (lambda (b) (list 'ok b))))
`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
