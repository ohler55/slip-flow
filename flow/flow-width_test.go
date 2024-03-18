// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestFlowWidthOk(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow :name 'flo)))
                  (flow-add-task flow :name "tisk" :x 100 :y 150 :actor (lambda (b) (list 'ok b)))
                  (flow-add-task flow :name "tusk" :x 200 :y 300 :actor (lambda (b) (list 'ok b)))
                  (flow-width flow))`,
		Expect: "264",
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow :name 'flo)))
                  (send flow :add-task :name "tisk" :x 100 :y 150 :actor (lambda (b) (list 'ok b)))
                  (send flow :add-task :name "tusk" :x 200 :y 300 :actor (lambda (b) (list 'ok b)))
                  (send flow :width))`,
		Expect: "264",
	}).Test(t)
}

func TestFlowWidthNotFlow(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-width t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowWidthArgCount(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-width (make-instance 'flow :name 'flo) t)`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}
