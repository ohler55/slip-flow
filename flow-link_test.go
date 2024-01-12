// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestFlowLinkFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `(let* ((flow (make-instance 'flow-flavor :name 'flo))
                        (tick (flow-add-task flow :name "tick" :actor (lambda (b) (list 'ok b))))
                        (tock (flow-add-task flow :name "tock" :actor (lambda (b) (list 'ok b)))))
                  (flow-link flow 'ok 'tick "tock")
                  (flow-task-links tick))`,
		Expect: `/\(\("ok" . #<flow-task-flavor [0-9a-f]+>\)\)/`,
	}).Test(t)
}

func TestFlowLinkSend(t *testing.T) {
	(&sliptest.Function{
		Source: `(let* ((flow (make-instance 'flow-flavor :name 'flo))
                        (tick (flow-add-task flow :name "tick" :actor (lambda (b) (list 'ok b))))
                        (tock (flow-add-task flow :name "tock" :actor (lambda (b) (list 'ok b)))))
                  (send flow :link 'ok 'tick 'tock)
                  (send tick :links))`,
		Expect: `/\(\("ok" . #<flow-task-flavor [0-9a-f]+>\)\)/`,
	}).Test(t)
}

func TestFlowLinkNilName(t *testing.T) {
	(&sliptest.Function{
		Source: `(let* ((flow (make-instance 'flow-flavor :name 'flo))
                        (tick (flow-add-task flow :name "tick" :actor (lambda (b) (list 'ok b))))
                        (tock (flow-add-task flow :name "tock" :actor (lambda (b) (list 'ok b)))))
                  (flow-link flow nil 'tick "tock")
                  (flow-task-links tick))`,
		Expect: `/\(\("" . #<flow-task-flavor [0-9a-f]+>\)\)/`,
	}).Test(t)
}

func TestFlowLinkNotFlow(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-link t 'ok 'tick 'tock)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowLinkBadName(t *testing.T) {
	(&sliptest.Function{
		Source: `(let* ((flow (make-instance 'flow-flavor :name 'flo))
                        (tick (flow-add-task flow :name "tick" :actor (lambda (b) (list 'ok b))))
                        (tock (flow-add-task flow :name "tock" :actor (lambda (b) (list 'ok b)))))
                  (flow-link flow t 'tick "tock"))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowLinkBadFrom(t *testing.T) {
	(&sliptest.Function{
		Source: `(let* ((flow (make-instance 'flow-flavor :name 'flo))
                        (tick (flow-add-task flow :name "tick" :actor (lambda (b) (list 'ok b))))
                        (tock (flow-add-task flow :name "tock" :actor (lambda (b) (list 'ok b)))))
                  (flow-link flow 'ok 'bad "tock"))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}

func TestFlowLinkBadTo(t *testing.T) {
	(&sliptest.Function{
		Source: `(let* ((flow (make-instance 'flow-flavor :name 'flo))
                        (tick (flow-add-task flow :name "tick" :actor (lambda (b) (list 'ok b))))
                        (tock (flow-add-task flow :name "tock" :actor (lambda (b) (list 'ok b)))))
                  (flow-link flow 'ok 'tick 'bad))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}
