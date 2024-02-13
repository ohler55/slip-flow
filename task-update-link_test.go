// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestTaskUpdateLinkFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let* ((flow (make-flow :name 'flo))
       (task (flow-add-task flow
                            :name "start"
                            :actor (lambda (b) (list 'ok b)))))
  (flow-add-task flow
                 :name "done"
                 :actor (make-instance 'flow-exit-actor))
  (flow-link flow 'ok 'start "done")
  (flow-task-update-link task 'ok '((100 200)))
  (flow-task-links task))
`,
		Expect: `/\(\("ok" #<flow-task [0-9a-f]+> \(100 200\)\)\)/`,
	}).Test(t)
}

func TestTaskUpdateLinkSend(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let* ((flow (make-flow :name 'flo))
       (task (send flow :add-task
                        :name "start"
                        :actor (lambda (b) (list 'ok b)))))
  (send flow :add-task
             :name "done"
             :actor (make-instance 'flow-exit-actor))
  (send flow :link 'ok 'start "done")
  (send task :update-link "ok" '((100 200)))
  (send task :links))
`,
		Expect: `/\(\("ok" #<flow-task [0-9a-f]+> \(100 200\)\)\)/`,
	}).Test(t)
}

func TestTaskUpdateLinkNotTask(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let* ((flow (make-flow :name 'flo))
       (task (send flow :add-task
                        :name "start"
                        :actor (lambda (b) (list 'ok b)))))
  (send flow :add-task
             :name "done"
             :actor (make-instance 'flow-exit-actor))
  (send flow :link 'ok 'start "done")
  (flow-task-update-link t 'ok nil))
`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskUpdateLinkBadName(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let* ((flow (make-flow :name 'flo))
       (task (send flow :add-task
                        :name "start"
                        :actor (lambda (b) (list 'ok b)))))
  (send flow :add-task
             :name "done"
             :actor (make-instance 'flow-exit-actor))
  (send flow :link 'ok 'start "done")
  (send task :update-link t '((100 200))))
`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskUpdateLinkNoLink(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let* ((flow (make-flow :name 'flo))
       (task (send flow :add-task
                        :name "start"
                        :actor (lambda (b) (list 'ok b)))))
  (send flow :add-task
             :name "done"
             :actor (make-instance 'flow-exit-actor))
  (send flow :link 'ok 'start "done")
  (send task :update-link "not-me" '((100 200))))
`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}

func TestTaskUpdateLinkBadMidPoints(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let* ((flow (make-flow :name 'flo))
       (task (send flow :add-task
                        :name "start"
                        :actor (lambda (b) (list 'ok b)))))
  (send flow :add-task
             :name "done"
             :actor (make-instance 'flow-exit-actor))
  (send flow :link 'ok 'start "done")
  (send task :update-link 'ok 100))
`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
