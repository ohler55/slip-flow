// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ohler55/ojg/tt"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestTaskActor(t *testing.T) {
	scope := slip.NewScope()
	_ = slip.ReadString(`
(defflavor task-actor-test () (flow-task-actor))
(defmethod (task-actor-test :perform) (box) (list 'ok box))
`).Eval(scope, nil)

	(&sliptest.Function{
		Scope:  scope,
		Source: `(let ((actor (make-instance 'task-actor-test))) (send actor :task))`,
		Expect: "nil",
	}).Test(t)

	(&sliptest.Function{
		Source: `(let* ((actor (make-instance 'task-actor-test))
                        (task (make-instance 'flow-task :actor actor)))
                  (send actor :start task)
                  (send actor :task))`,
		Expect: "/#<flow-task [0-9a-f]+>/",
	}).Test(t)
}

func TestTaskActorDocs(t *testing.T) {
	scope := slip.NewScope()
	var out strings.Builder
	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})

	for _, method := range []string{
		":start",
	} {
		_ = slip.ReadString(fmt.Sprintf(`(describe-method flow-task-actor %s out)`, method)).Eval(scope, nil)
		tt.Equal(t, true, strings.Contains(out.String(), method))
		out.Reset()
	}
}
