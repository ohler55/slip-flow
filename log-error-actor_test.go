// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/ohler55/ojg/tt"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/gi"
	"github.com/ohler55/slip/sliptest"
)

func TestLogErrorActorExit(t *testing.T) {
	exitChan := make(gi.Channel, 5)
	scope := slip.NewScope()
	scope.Let("exit-channel", exitChan)

	var b bytes.Buffer
	orig := scope.Get("*standard-output*")
	defer scope.Set("*standard-output*", orig)
	scope.Set("*standard-output*", &slip.OutputStream{Writer: &b})

	(&sliptest.Function{
		Scope: scope,
		Source: `
(let* ((lg (make-instance 'logger-flavor))
       (flow (make-flow :name 'flo :exit-channel exit-channel :logger lg)))
  (flow-add-task flow
                 :name "start"
                 :actor (lambda (b)
                          (flow-box-set b (* 3 (flow-box-get b "[0]")) "[0]")
                          (list 'ok b)))
  (flow-add-task flow
                 :name "fail"
                 :actor (lambda (b) (list 'x b nil)))
  (flow-add-task flow
                 :name "error"
                 :actor (make-instance 'flow-log-error-actor))
  (flow-link flow 'ok 'start "fail")
  (flow-set-entry flow 'start)
  (send flow :set-level 'warn)
  (flow-submit flow (make-flow-box :set '(1)))
  (sleep 0.1)
  (send lg :shutdown))`,
		Expect: "nil",
	}).Test(t)

	out := <-exitChan
	scope.Let("log-error-test-out", out)

	history := slip.ReadString(
		`(mapcar (lambda (ev) (cadr ev))(send (send log-error-test-out :track) :history))`).Eval(scope, nil)
	tt.Equal(t, `("start" "fail" "error")`, slip.ObjectString(history))

	content := slip.ReadString(`(cdr (assoc "content" (send log-error-test-out :native)))`).Eval(scope, nil)
	tt.Equal(t, `(3)`, slip.ObjectString(content))

	err := slip.ReadString(`(cdr (assoc "error" (send log-error-test-out :native)))`).Eval(scope, nil)
	tt.Equal(t, `"Actor did not return a list of link name and box instance."`, slip.ObjectString(err))

	tt.Equal(t,
		`/E flo:fail #<uuid [0-9a-f-]+> - Actor did not return a list of link name and box instance./`,
		b.String())
}

func TestLogErrorActorLink(t *testing.T) {
	exitChan := make(gi.Channel, 5)
	scope := slip.NewScope()
	scope.Let("exit-channel", exitChan)

	var b bytes.Buffer
	orig := scope.Get("*standard-output*")
	defer scope.Set("*standard-output*", orig)
	scope.Set("*standard-output*", &slip.OutputStream{Writer: &b})

	(&sliptest.Function{
		Scope: scope,
		Source: `
(let* ((lg (make-instance 'logger-flavor))
       (flow (make-flow :name 'flo :exit-channel exit-channel :logger lg)))
  (flow-add-task flow
                 :name "start"
                 :actor (lambda (b)
                          (flow-box-set b (* 3 (flow-box-get b "[0]")) "[0]")
                          (list 'ok b)))
  (flow-add-task flow
                 :name "fail"
                 :actor (lambda (b) (list 'ok b)))
  (flow-add-task flow
                 :name "error"
                 :actor (make-instance 'flow-log-error-actor))
  (flow-add-task flow
                 :name "error2"
                 :actor (make-instance 'flow-exit-actor))
  (flow-link flow 'ok 'start "fail")
  (flow-link flow 'ok 'fail "error")
  (flow-link flow 'ok 'error "error2")
  (flow-set-entry flow 'start)
  (send flow :set-level 'warn)
  (flow-submit flow (make-flow-box :set '(1)))
  (sleep 0.1)
  (send lg :shutdown))`,
		Expect: "nil",
	}).Test(t)

	out := <-exitChan
	scope.Let("log-error-test-out", out)

	history := slip.ReadString(
		`(mapcar (lambda (ev) (cadr ev))(send (send log-error-test-out :track) :history))`).Eval(scope, nil)
	tt.Equal(t, `("start" "fail" "error" "error2")`, slip.ObjectString(history))

	tt.Equal(t, "E [3]\n", b.String())
}

func TestLogErrorActorDocs(t *testing.T) {
	scope := slip.NewScope()
	var out strings.Builder
	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})

	for _, method := range []string{
		":start",
		":perform",
	} {
		_ = slip.ReadString(fmt.Sprintf(`(describe-method flow-log-error-actor %s out)`, method)).Eval(scope, nil)
		tt.Equal(t, true, strings.Contains(out.String(), method))
		out.Reset()
	}
}
