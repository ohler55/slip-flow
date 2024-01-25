// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/gi"
	"github.com/ohler55/slip/sliptest"
)

// The flow used for these tests is:
//
//                    ┏━━━━━━━┓
// ┏━━━━━━━┓── one ──>┃ + x:1 ┃── ok ──>┏━━━━━━━┓         ┏━━━━━━┓
// ┃ split ┃          ┗━━━━━━━┛         ┃ merge ┃── ok ──>┃ exit ┃
// ┃       ┃          ┏━━━━━━━┓         ┃       ┃         ┗━━━━━━┛
// ┗━━━━━━━┛── two ──>┃ + y:2 ┃── ok ──>┗━━━━━━━┛
//                    ┗━━━━━━━┛

func TestMergeActorOk(t *testing.T) {
	exitChan := make(gi.Channel, 5)
	scope := slip.NewScope()
	scope.Let("exit-channel", exitChan)

	(&sliptest.Function{
		Scope: scope,
		Source: `
(let ((flow (make-flow :name 'flo :exit-channel exit-channel)))
  (flow-add-task flow
                 :name "split"
                 :actor (make-instance 'flow-split-actor :links '(one two)))
  (flow-add-task flow
                 :name "x"
                 :actor (lambda (b) (flow-box-get b "x" 1) (list 'ok b)))
  (flow-add-task flow
                 :name "y"
                 :actor (lambda (b) (flow-box-get b "y" 2) (list 'ok b)))
  (flow-add-task flow
                 :name "merge"
                 :actor (make-instance 'flow-merge-actor) :number 2 :timeout 2)
  (flow-add-task flow
                 :name "exit"
                 :actor (make-instance 'flow-exit-actor))

  (flow-link flow 'one 'split 'x)
  (flow-link flow 'two 'split 'y)
  (flow-link flow 'ok "x" 'merge)
  (flow-link flow 'ok "y" 'merge)
  (flow-link flow 'ok 'merge 'exit)
  (flow-set-entry flow 'split)
  (send flow :set-level 'warn)
  (flow-submit flow (make-flow-box :parse "{z:0}"))
  (flow-shutdown flow))`,
		Expect: "nil",
	}).Test(t)

	out := <-exitChan
	scope.Let("merge-test-out", out)

	history := slip.ReadString(
		`(mapcar (lambda (ev) (cadr ev))(send (send merge-test-out :track) :history))`).Eval(scope, nil)
	fmt.Printf("*** history: %s\n", history)

	value := slip.ReadString(`(send merge-test-out :native)`).Eval(scope, nil)
	fmt.Printf("*** value: %s\n", value)
}
