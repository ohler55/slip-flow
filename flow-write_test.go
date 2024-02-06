// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestFlowWriteSend(t *testing.T) {
	scope := slip.NewScope()
	scope.Let("*print-right-margin*", slip.Fixnum(80))
	(&sliptest.Function{
		Scope: scope,
		Source: `
(let ((flow (make-flow :name 'flo)))
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
  (flow-link flow 'ok 'start "odd-or-even")
  (flow-link flow 'odd "odd-or-even" 'odd)
  (flow-link flow 'even "odd-or-even" 'even)
  (flow-set-entry flow 'start)
  (send flow :write nil))`,
		Expect: `"(let ((flow (make-flow :name "flo")))
  (send flow :add-task
        :name "even"
        :actor (make-instance 'flow-exit-actor)
  (send flow :add-task
        :name "odd"
        :actor (make-instance 'flow-exit-actor)
  (send flow :add-task
        :name "odd-or-even"
        :actor (lambda (b)
                       (cond ((evenp (flow-box-get b "[0]")) (list 'even b))
                             (t (list 'odd b))))
  (send flow :add-task
        :name "start"
        :actor (lambda (b) (flow-box-set b (* 3 (flow-box-get b "[0]")) "[0]")
                       (list 'ok b))
  (send flow :link "even" "odd-or-even" "even")
  (send flow :link "odd" "odd-or-even" "odd")
  (send flow :link "ok" "start" "odd-or-even")
  (send flow :set-entry "start")
  flow)
"`,
	}).Test(t)
}
