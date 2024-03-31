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

func TestForeachActorList(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (flow-add-task flow
                :name "start"
                :actor (make-instance 'flow-foreach-actor
                                      :list '(1 3 5)
                                      :destination "value"))
 (flow-add-task flow
                :name "done"
                :actor (make-instance 'flow-exit-actor))

 (flow-link flow 'ok 'start "done")
 (flow-set-entry flow 'start)
 (flow-submit flow (make-flow-box :set nil :watch 'done))
 (list (flow-box-write (channel-pop done)) (flow-box-write (channel-pop done)) (flow-box-write (channel-pop done))))`,
		Validate: func(t *testing.T, v slip.Object) {
			list := v.(slip.List)
			for i, x := range []string{
				`{value: 1}`,
				`{value: 3}`,
				`{value: 5}`,
			} {
				tt.Equal(t, x, string(list[i].(slip.String)))
			}
		},
	}).Test(t)
}

func TestForeachActorLambda(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (flow-add-task flow
                :name "start"
                :actor (make-instance 'flow-foreach-actor
                                      :list (lambda (b) '(1 3 5))
                                      :destination "value"))
 (flow-add-task flow
                :name "done"
                :actor (make-instance 'flow-exit-actor))

 (flow-link flow 'ok 'start "done")
 (flow-set-entry flow 'start)
 (flow-submit flow (make-flow-box :set nil :watch 'done))
 (list (flow-box-write (channel-pop done)) (flow-box-write (channel-pop done)) (flow-box-write (channel-pop done))))`,
		Validate: func(t *testing.T, v slip.Object) {
			list := v.(slip.List)
			for i, x := range []string{
				`{value: 1}`,
				`{value: 3}`,
				`{value: 5}`,
			} {
				tt.Equal(t, x, string(list[i].(slip.String)))
			}
		},
	}).Test(t)
}

func TestReadCSVActorBadLambda(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (flow-add-task flow
                :name "start"
                :actor (make-instance 'flow-foreach-actor
                                      :list (lambda (b) t)
                                      :destination "value"))
 (flow-add-task flow
                :name "done"
                :actor (make-instance 'flow-exit-actor))
 (flow-add-task flow
                :name "error"
                :actor (make-instance 'flow-exit-actor))

 (flow-link flow 'ok 'start "done")
 (flow-set-entry flow 'start)
 (flow-submit flow (make-flow-box :set nil :watch 'done))
 (flow-box-write (channel-pop done)))`,
		Expect: `"{content: null error: "list must be a list not t, a t."}"`,
	}).Test(t)
}

func TestForeachActorBadDestination(t *testing.T) {
	(&sliptest.Function{
		Source: `(make-instance 'flow-foreach-actor
                                :list '(1 3 5)
                                :destination t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestForeachActorBadList(t *testing.T) {
	(&sliptest.Function{
		Source: `(make-instance 'flow-foreach-actor
                                :list t
                                :destination "value")`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestForeachActorDocs(t *testing.T) {
	scope := slip.NewScope()
	var out strings.Builder
	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})

	for _, method := range []string{
		":init",
		":start",
		":perform",
		":init-key-values",
	} {
		_ = slip.ReadString(fmt.Sprintf(`(describe-method flow-foreach-actor %s out)`, method)).Eval(scope, nil)
		tt.Equal(t, true, strings.Contains(out.String(), method))
		out.Reset()
	}
}

func TestForeachActorInitKeyValues(t *testing.T) {
	(&sliptest.Function{
		Source: `(send (make-instance 'flow-foreach-actor
                                      :list '(1 3 5)
                                      :destination 'value) :init-key-values)`,
		Validate: func(t *testing.T, v slip.Object) {
			for i, x := range []string{
				`:list`,
				`(1 3 5)`,
				`:destination`,
				`"value"`,
			} {
				tt.Equal(t, x, slip.ObjectString(v.(slip.List)[i]), x)
			}
		},
	}).Test(t)
}
