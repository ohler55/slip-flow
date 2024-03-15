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

func TestJumpActorBasic(t *testing.T) {
	scope := slip.NewScope()
	var b bytes.Buffer
	orig := scope.Get("*standard-output*")
	defer scope.Set("*standard-output*", orig)
	scope.Set("*standard-output*", &slip.OutputStream{Writer: &b})

	(&sliptest.Function{
		Scope: scope,
		Source: `
(let* ((done (make-channel 3))
       (lg (make-instance 'logger-flavor))
       (group (make-flow-group :logger lg))
       (s0 (make-instance 'flow :name 'first-stage))
       (s1 (make-instance 'flow :name 'second-stage)))
 (flow-add-task s0
                :name "triple"
                :actor (lambda (b)
                         (flow-box-set b (* 3 (flow-box-get b "[0]")) "[0]")
                         (list 'ok b)))
 (flow-add-task s0
                :name "jump"
                :actor (make-instance 'flow-jump-actor :target 'second-stage))
 (flow-link s0 'ok 'triple 'jump)
 (flow-set-entry s0 'triple)
 (flow-group-add group s0)

 (flow-add-task s1
                :name "double"
                :actor (lambda (b)
                         (flow-box-set b (* 2 (flow-box-get b "[0]")) "[0]")
                         (list 'ok b)))
 (flow-add-task s1
                :name "done"
                :actor (make-instance 'flow-exit-actor))
 (flow-link s1 'ok 'double 'done)
 (flow-set-entry s1 'double)
 (flow-group-add group s1)

 (send group :set-level 'info)
 (flow-group-start group)
 (flow-submit s0 (make-flow-box :set '(1) :watch 'done))
 (channel-pop done)
 (flow-group-shutdown group)
 (send lg :shutdown))
`,
		Expect: "nil",
	}).Test(t)

	tt.Equal(t, `/I first-stage:triple received box .*
I first-stage:triple following ok with box .*
I first-stage:jump received box .*
I second-stage:double received box .*
I second-stage:double following ok with box .*
I second-stage:done received box .*
/`, b.String())
}

func TestJumpActorDocs(t *testing.T) {
	scope := slip.NewScope()
	var out strings.Builder
	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})

	for _, method := range []string{
		":init",
		":start",
		":perform",
	} {
		_ = slip.ReadString(fmt.Sprintf(`(describe-method flow-jump-actor %s out)`, method)).Eval(scope, nil)
		tt.Equal(t, true, strings.Contains(out.String(), method))
		out.Reset()
	}
}

func TestJumpActorBadTarget(t *testing.T) {
	(&sliptest.Function{
		Source:    `(make-instance 'flow-jump-actor :target t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestJumpActorNotInGroup(t *testing.T) {
	scope := slip.NewScope()
	var b bytes.Buffer
	orig := scope.Get("*standard-output*")
	defer scope.Set("*standard-output*", orig)
	scope.Set("*standard-output*", &slip.OutputStream{Writer: &b})

	(&sliptest.Function{
		Scope: scope,
		Source: `
(let ((s0 (make-instance 'flow :name 'first-stage)))
  (flow-add-task s0
                 :name "triple"
                 :actor (lambda (b)
                          (flow-box-set b (* 3 (flow-box-get b "[0]")) "[0]")
                          (list 'ok b)))
  (flow-add-task s0
                 :name "jump"
                 :actor (make-instance 'flow-jump-actor :target "second-stage"))
  (flow-link s0 'ok 'triple 'jump)
  (flow-set-entry s0 'triple)
  (flow-submit s0 (make-flow-box :set '(1)))
  (flow-shutdown s0)
  (sleep 0.1))
`,
		Expect: "nil",
	}).Test(t)

	tt.Equal(t, `/E first-stage:jump .* task is not in a flow that is in a group/`, b.String())
}

func TestJumpActorNotInSameGroup(t *testing.T) {
	scope := slip.NewScope()
	var b bytes.Buffer
	orig := scope.Get("*standard-output*")
	defer scope.Set("*standard-output*", orig)
	scope.Set("*standard-output*", &slip.OutputStream{Writer: &b})

	(&sliptest.Function{
		Scope: scope,
		Source: `
(let ((group (make-flow-group))
      (s0 (make-instance 'flow :name 'first-stage)))
  (flow-add-task s0
                 :name "triple"
                 :actor (lambda (b)
                          (flow-box-set b (* 3 (flow-box-get b "[0]")) "[0]")
                          (list 'ok b)))
  (flow-add-task s0
                 :name "jump"
                 :actor (make-instance 'flow-jump-actor :target "second-stage"))
  (flow-link s0 'ok 'triple 'jump)
  (flow-set-entry s0 'triple)
  (flow-group-add group s0)

  (flow-group-start group)
  (flow-submit s0 (make-flow-box :set '(1)))
  (flow-group-shutdown group)
  (sleep 0.1))
`,
		Expect: "nil",
	}).Test(t)

	tt.Equal(t, `/E first-stage:jump .* flow second-stage is not in the same group that first-stage is in/`, b.String())
}
