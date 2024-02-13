// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestGroupStartFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((group (make-flow-group))
                       (flow (make-instance 'flow :name 'flo))
                       running)
                  (flow-add-task flow :name "tisk" :workers 1 :actor (lambda (b) (list 'ok b)))
                  (flow-group-add group flow)
                  (flow-group-start group)
                  (setq running (flow-running flow))
                  (flow-group-shutdown group)
                  running)`,
		Expect: "t",
	}).Test(t)
}

func TestGroupStartSend(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((group (make-instance 'flow-group))
                       (flow (make-instance 'flow :name 'flo))
                       running)
                  (send flow :add-task :name "tisk" :workers 1 :actor (lambda (b) (list 'ok b)))
                  (send group :add flow)
                  (send group :start)
                  (setq running (send flow :running))
                  (send group :shutdown)
                  running)`,
		Expect: "t",
	}).Test(t)
}

func TestGroupStartNotGroup(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-group-start t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestGroupStartArgCount(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-group-start (make-instance 'flow-group) t)`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}
