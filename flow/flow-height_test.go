// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestFlowHeightOk(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow :name 'flo)))
                  (flow-add-task flow :name "tisk" :x 100 :y 150 :actor (lambda (b) (list 'ok b)))
                  (flow-add-task flow :name "tusk" :x 200 :y 300 :actor (lambda (b) (list 'ok b)))
                  (flow-height flow))`,
		Expect: "364",
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow :name 'flo)))
                  (send flow :add-task :name "tisk" :x 100 :y 150 :actor (lambda (b) (list 'ok b)))
                  (send flow :add-task :name "tusk" :x 200 :y 300 :actor (lambda (b) (list 'ok b)))
                  (send flow :height))`,
		Expect: "364",
	}).Test(t)
}

func TestFlowHeightNotFlow(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-height t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowHeightArgCount(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-height (make-instance 'flow :name 'flo) t)`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}
