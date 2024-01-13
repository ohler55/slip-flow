// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestFlowResetMetricsFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((flow (make-instance 'flow-flavor :name 'flo)))
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
  (flow-add-task flow
                 :name "error"
                 :actor (make-instance 'flow-exit-actor))
  (flow-link flow 'ok 'start "odd-or-even")
  (flow-link flow 'odd "odd-or-even" 'odd)
  (flow-link flow 'even "odd-or-even" 'even)
  (flow-set-entry flow 'start)
  (send flow :set-level 'warn)
  (flow-submit flow (make-flow-box :set '(1)))
  (flow-reset-metrics flow)
  (flow-metrics flow))`,
		Expect: `((received . 0) (processed . 0) (errors . 0))`,
	}).Test(t)
}

func TestFlowResetMetricsSend(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((flow (make-instance 'flow-flavor :name 'flo)))
  (send flow :add-task
             :name "start"
             :actor (lambda (b)
                      (flow-box-set b (* 3 (flow-box-get b "[0]")) "[0]")
                      (list 'ok b)))
  (send flow :add-task
             :name "odd-or-even"
             :actor (lambda (b)
                      (cond ((evenp (flow-box-get b "[0]")) (list 'even b))
                            (t (list 'odd b)))))
  (send flow :add-task
             :name "even"
             :actor (make-instance 'flow-exit-actor))
  (send flow :add-task
             :name "odd"
             :actor (make-instance 'flow-exit-actor))
  (send flow :link 'ok 'start "odd-or-even")
  (send flow :link 'odd "odd-or-even" 'odd)
  (send flow :link 'even "odd-or-even" 'even)
  (send flow :set-entry 'start)
  (send flow :set-level 'warn)
  (send flow :submit (make-flow-box :set '(1)))
  (send flow :reset-metrics)
  (send flow :metrics))`,
		Expect: `((received . 0) (processed . 0) (errors . 0))`,
	}).Test(t)
}

func TestFlowResetMetricsNotFlow(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-reset-metrics t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
