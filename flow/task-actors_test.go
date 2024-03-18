// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestTaskActorsLambda(t *testing.T) {
	(&sliptest.Function{
		Source: `(let* ((flow (make-instance 'flow :name 'flo))
                        (tick (flow-add-task flow :name "tick" :actor (lambda (b) (list 'ok b)))))
                  (flow-task-actors tick))`,
		Expect: `/\(#<function \(lambda \(b\)\) \{[0-9a-f]+\}>\)/`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(let* ((flow (make-instance 'flow :name 'flo))
                        (tick (send flow :add-task :name "tick" :actor (lambda (b) (list 'ok b)))))
                  (send tick :actors))`,
		Expect: `/\(#<function \(lambda \(b\)\) \{[0-9a-f]+\}>\)/`,
	}).Test(t)
}

func TestTaskActorsSymbol(t *testing.T) {
	(&sliptest.Function{
		// the function symbol isn't evaluated so anything is fine
		Source: `(let* ((flow (make-instance 'flow :name 'flo))
                        (tick (flow-add-task flow :name "tick" :actor 'car)))
                  (flow-task-actors tick))`,
		Expect: `(car)`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(let* ((flow (make-instance 'flow :name 'flo))
                        (tick (send flow :add-task :name "tick" :actor 'car)))
                  (send tick :actors))`,
		Expect: `(car)`,
	}).Test(t)
}

func TestTaskActorsInstance(t *testing.T) {
	(&sliptest.Function{
		Source: `(let* ((flow (make-instance 'flow :name 'flo))
                        (tick (flow-add-task flow :name "tick" :actor (make-instance 'flow-exit-actor))))
                  (flow-task-actors tick))`,
		Expect: `/\(#<flow-exit-actor [0-9a-f]+>\)/`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(let* ((flow (make-instance 'flow :name 'flo))
                        (tick (send flow :add-task :name "tick" :actor (make-instance 'flow-exit-actor))))
                  (send tick :actors))`,
		Expect: `/\(#<flow-exit-actor [0-9a-f]+>\)/`,
	}).Test(t)
}

func TestTaskActorsNotTask(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-task-actors t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
