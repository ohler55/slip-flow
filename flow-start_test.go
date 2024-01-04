// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

// func TestFlowStartOk(t *testing.T) {
// 	(&sliptest.Function{
// 		Source: `(let ((flow (make-instance 'flow-flavor :name 'flo))
//                        running)
//                   (flow-add-task flow :name "tisk" :actor (lambda (b) (list 'ok b)))
//                   (flow-start flow)
//                   (setq running (flow-running flow))
//                   (flow-shutdown flow)
//                   running)`,
// 		Expect: "t",
// 	}).Test(t)
// }

func TestFlowStartNotFlow(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-start t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowStartArgCount(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-start (make-instance 'flow-flavor :name 'flo) t)`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}
