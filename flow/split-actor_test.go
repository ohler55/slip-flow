// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ohler55/ojg/tt"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestSplitActor(t *testing.T) {
	scope := slip.NewScope()
	tf := sliptest.Function{
		Scope: scope,
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
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
 (flow-submit flow (make-flow-box :set '(1) :watch 'done))
 (list (channel-pop done) (channel-pop done)))`,
		Expect: `/\(#<flow-box [0-9a-f]+> #<flow-box [0-9a-f]+>\)/`,
	}
	tf.Test(t)
	scope.Let("split-out", tf.Result.(slip.List)[0])
	history := slip.ReadString(
		`(mapcar (lambda (ev) (cadr ev))(send (send split-out :track) :history))`).Eval(scope, nil).(slip.List)

	scope.Let("split-out", tf.Result.(slip.List)[1])
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
