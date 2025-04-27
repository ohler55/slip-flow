// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/ohler55/ojg/tt"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestDeleteFileActorTextStatic(t *testing.T) {
	filename := "testdata/delete-me.txt"
	_ = os.WriteFile(filename, []byte("delete-me"), 0666)
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (flow-add-task flow
                :name "start"
                :actor (make-instance 'flow-delete-file-actor
                                      :filename "testdata/delete-me.txt"))
 (flow-add-task flow
                :name "done"
                :actor (make-instance 'flow-exit-actor))

 (flow-link flow 'ok 'start "done")
 (flow-set-entry flow 'start)
 (flow-submit flow (make-flow-box :set nil :watch 'done))
 (channel-pop done))`,
		Expect: "/#<flow-box [0-9a-f]+>/",
	}).Test(t)
	_, err := os.Stat(filename)
	tt.NotNil(t, err)
}

func TestDeleteFileActorInitKeyValues(t *testing.T) {
	(&sliptest.Function{
		Source: `(send (make-instance 'flow-delete-file-actor
                                      :filename "testdata/delete-me.txt") :init-key-values)`,
		Validate: func(t *testing.T, v slip.Object) {
			for i, x := range []string{
				`:filename`,
				`"testdata/delete-me.txt"`,
			} {
				tt.Equal(t, x, slip.ObjectString(v.(slip.List)[i]), x)
			}
		},
	}).Test(t)
}

func TestDeleteFileActorDocs(t *testing.T) {
	scope := slip.NewScope()
	var out strings.Builder
	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})

	for _, method := range []string{
		":init",
		":start",
		":perform",
		":init-key-values",
	} {
		_ = slip.ReadString(fmt.Sprintf(`(describe-method flow-delete-file-actor %s out)`, method), scope).Eval(scope, nil)
		tt.Equal(t, true, strings.Contains(out.String(), method))
		out.Reset()
	}
}
