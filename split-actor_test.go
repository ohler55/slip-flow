// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ohler55/ojg/tt"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/gi"
	"github.com/ohler55/slip/sliptest"
)

func TestSplitActor(t *testing.T) {
	exitChan := make(gi.Channel, 5)
	scope := slip.NewScope()
	scope.Let("exit-channel", exitChan)

	(&sliptest.Function{
		Scope: scope,
		Source: `
(let ((flow (make-flow :name 'flo :exit-channel exit-channel)))
  (flow-add-task flow
                 :name "start"
                 :actor (make-instance 'flow-split-actor))
  (flow-add-task flow
                 :name "branch-one"
                 :actor (make-instance 'flow-exit-actor))
  (flow-add-task flow
                 :name "branch-two"
                 :actor (make-instance 'flow-exit-actor))

  (flow-link flow 'one 'start "branch-one")
  (flow-link flow 'two 'start "branch-two")
  (flow-set-entry flow 'start)
  (flow-submit flow (make-flow-box :set '(1))))`,
		Expect: "nil",
	}).Test(t)

	out := <-exitChan
	scope.Let("split-out", out)
	history := slip.ReadString(
		`(mapcar (lambda (ev) (cadr ev))(send (send split-out :track) :history))`).Eval(scope, nil).(slip.List)

	out = <-exitChan
	scope.Let("split-out", out)
	history = append(history,
		slip.ReadString(
			`(mapcar (lambda (ev) (cadr ev))(send (send split-out :track) :history))`).Eval(scope, nil).(slip.List)...,
	)
	hstr := slip.ObjectString(history)
	tt.Equal(t, true, strings.Contains(hstr, "branch-one"))
	tt.Equal(t, true, strings.Contains(hstr, "branch-two"))
}

func TestSplitActorDocs(t *testing.T) {
	scope := slip.NewScope()
	var out strings.Builder
	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})

	for _, method := range []string{
		":start",
		":perform",
	} {
		_ = slip.ReadString(fmt.Sprintf(`(describe-method flow-split-actor %s out)`, method)).Eval(scope, nil)
		tt.Equal(t, true, strings.Contains(out.String(), method))
		out.Reset()
	}
}
