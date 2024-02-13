// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestFlowStartFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow :name 'flo))
                       running)
                  (flow-add-task flow :name "tisk" :workers 1 :actor (lambda (b) (list 'ok b)))
                  (flow-start flow)
                  (setq running (flow-running flow))
                  (flow-shutdown flow)
                  running)`,
		Expect: "t",
	}).Test(t)
}

func TestFlowStartSend(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((flow (make-instance 'flow :name 'flo))
                       running)
                  (send flow :add-task :name "tisk" :workers 1 :actor (lambda (b) (list 'ok b)))
                  (send flow :start)
                  (setq running (send flow :running))
                  (send flow :shutdown)
                  running)`,
		Expect: "t",
	}).Test(t)
}

func TestFlowStartNotFlow(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-start t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowStartArgCount(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-start (make-instance 'flow :name 'flo) t)`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}
