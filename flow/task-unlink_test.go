// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestTaskUnlinkFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `(let* ((flow (make-instance 'flow :name 'flo))
                        (tick (flow-add-task flow :name "tick" :actor (lambda (b) (list 'ok b))))
                        (tock (flow-add-task flow :name "tock" :actor (lambda (b) (list 'ok b)))))
                  (flow-link flow 'ok 'tick "tock")
                  (flow-task-unlink tick 'ok)
                  (flow-task-links tick))`,
		Expect: "nil",
	}).Test(t)
}

func TestTaskUnlinkSend(t *testing.T) {
	(&sliptest.Function{
		Source: `(let* ((flow (make-instance 'flow :name 'flo))
                        (tick (send flow :add-task :name "tick" :actor (lambda (b) (list 'ok b))))
                        (tock (send flow :add-task :name "tock" :actor (lambda (b) (list 'ok b)))))
                  (send flow :link 'ok 'tick "tock")
                  (send tick :unlink "ok")
                  (send tick :links))`,
		Expect: "nil",
	}).Test(t)
}

func TestTaskUnlinkNotTask(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-task-unlink t 'ok)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskUnlinkBadName(t *testing.T) {
	(&sliptest.Function{
		Source: `(let* ((flow (make-instance 'flow :name 'flo))
                        (tick (send flow :add-task :name "tick" :actor (lambda (b) (list 'ok b)))))
                  (send tick :unlink t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
