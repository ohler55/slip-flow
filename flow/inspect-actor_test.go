// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/ohler55/ojg/tt"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestInspectActorStandardOutput(t *testing.T) {
	scope := slip.NewScope()

	var b bytes.Buffer
	orig := scope.Get("*standard-output*")
	defer scope.Set("*standard-output*", orig)
	scope.Set("*standard-output*", &slip.OutputStream{Writer: &b})

	(&sliptest.Function{
		Scope: scope,
		Source: `
(let* ((lg (make-instance 'logger-flavor))
       (flow (make-flow :name 'flo :logger lg)))
  (flow-add-task flow
                 :name "start"
                 :actor (lambda (b)
                          (flow-box-set b (* 3 (flow-box-get b "[0]")) "[0]")
                          (list 'ok b)))
  (flow-add-task flow
                 :name "show-me"
                 :actor (make-instance 'flow-inspect-actor :output nil))
  (flow-add-task flow
                 :name "done"
                 :actor (make-instance 'flow-exit-actor))
  (flow-link flow 'ok 'start 'show-me)
  (flow-link flow 'ok 'show-me 'done)
  (flow-set-entry flow 'start)
  (flow-submit flow (make-flow-box :set '(1)))
  (flow-shutdown flow))`,
		Expect: "nil",
	}).Test(t)

	tt.Equal(t, `[3]`, b.String())
}

func TestInspectActorLog(t *testing.T) {
	scope := slip.NewScope()

	var b bytes.Buffer
	orig := scope.Get("*standard-output*")
	defer scope.Set("*standard-output*", orig)
	scope.Set("*standard-output*", &slip.OutputStream{Writer: &b})

	(&sliptest.Function{
		Scope: scope,
		Source: `
(let* ((lg (make-instance 'logger-flavor))
       (flow (make-flow :name 'flo :logger lg)))
  (flow-add-task flow
                 :name "start"
                 :actor (lambda (b)
                          (flow-box-set b (* 3 (flow-box-get b "[0]")) "[0]")
                          (list 'ok b)))
  (flow-add-task flow
                 :name "show-me"
                 :actor (make-instance 'flow-inspect-actor :output ':warn))
  (flow-add-task flow
                 :name "done"
                 :actor (make-instance 'flow-exit-actor))
  (flow-link flow 'ok 'start 'show-me)
  (flow-link flow 'ok 'show-me 'done)
  (send flow :set-level 'warn)
  (flow-set-entry flow 'start)
  (flow-submit flow (make-flow-box :set '(1)))
  (sleep 0.1)
  (flow-shutdown flow))`,
		Expect: "nil",
	}).Test(t)

	tt.Equal(t, "W [3]\n", b.String())
}

func TestInspectActorLinks(t *testing.T) {
	(&sliptest.Function{
		Source: `(send (make-instance 'flow-inspect-actor) :links)`,
		Expect: `("ok")`,
	}).Test(t)
}

func TestInspectActorBadOutput(t *testing.T) {
	scope := slip.NewScope()
	var b bytes.Buffer
	orig := scope.Get("*standard-output*")
	defer scope.Set("*standard-output*", orig)
	scope.Set("*standard-output*", &slip.OutputStream{Writer: &b})
	(&sliptest.Function{
		Scope: scope,
		Source: `
(let* ((flow (make-flow :name 'flo)))
  (flow-add-task flow
                 :name "start"
                 :actor (lambda (b)
                          (flow-box-set b (* 3 (flow-box-get b "[0]")) "[0]")
                          (list 'ok b)))
  (flow-add-task flow
                 :name "show-me"
                 :actor (make-instance 'flow-inspect-actor :output t))

  (flow-set-entry flow 'show-me)
  (flow-submit flow (make-flow-box :set '(1)))
  (sleep 0.1))`,
		Expect: "nil",
	}).Test(t)
	tt.Equal(t, "/E flo:show-me/", b.String())
}

func TestInspectActorDocs(t *testing.T) {
	scope := slip.NewScope()
	var out strings.Builder
	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})

	for _, method := range []string{
		":init",
		":start",
		":perform",
		":links",
		":init-key-values",
	} {
		_ = slip.ReadString(fmt.Sprintf(`(describe-method flow-inspect-actor %s out)`, method), scope).Eval(scope, nil)
		tt.Equal(t, true, strings.Contains(out.String(), method))
		out.Reset()
	}
}

func TestInspectActorInitKeyValues(t *testing.T) {
	(&sliptest.Function{
		Source: `(send (make-instance 'flow-inspect-actor
                                      :depth 4
                                      :right-margin 80
                                      :indent 2
                                      :time-format "2006-01-02"
                                      :time-wrap "@"
                                      :json t
                                      :color t
                                      :pretty t
                                      :full t
                                      :output :info) :init-key-values)`,
		Validate: func(t *testing.T, v slip.Object) {
			for i, x := range []string{
				`:depth`,
				`4`,
				`:right-margin`,
				`80`,
				`:indent`,
				`2`,
				`:time-format`,
				`"2006-01-02"`,
				`:time-wrap`,
				`"@"`,
				`:output`,
				`:info`,
				`:json`,
				`t`,
				`:color`,
				`t`,
				`:pretty`,
				`t`,
				`:full`,
				`t`,
			} {
				tt.Equal(t, x, slip.ObjectString(v.(slip.List)[i]), x)
			}
		},
	}).Test(t)
	(&sliptest.Function{
		Source: `(send (make-instance 'flow-inspect-actor
                                      :depth 4
                                      :right-margin 80
                                      :indent 2
                                      :time-format "2006-01-02"
                                      :time-wrap "@"
                                      :json nil
                                      :color nil
                                      :pretty nil
                                      :full nil
                                      :output nil) :init-key-values)`,
		Validate: func(t *testing.T, v slip.Object) {
			for i, x := range []string{
				`:depth`,
				`4`,
				`:right-margin`,
				`80`,
				`:indent`,
				`2`,
				`:time-format`,
				`"2006-01-02"`,
				`:time-wrap`,
				`"@"`,
				`:output`,
				`nil`,
				`:json`,
				`nil`,
				`:color`,
				`nil`,
				`:pretty`,
				`nil`,
				`:full`,
				`nil`,
			} {
				tt.Equal(t, x, slip.ObjectString(v.(slip.List)[i]), x)
			}
		},
	}).Test(t)
}
