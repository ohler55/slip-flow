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

func TestWriteFileActorTextStatic(t *testing.T) {
	filename := "testdata/write-me.txt"
	_ = os.RemoveAll(filename)
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (flow-add-task flow
                :name "start"
                :actor (make-instance 'flow-write-file-actor
                                      :filename "testdata/write-me.txt"
                                      :content "Sample file"
                                      :overwrite t
                                      :permissions "-r--r--r--"))
 (flow-add-task flow
                :name "done"
                :actor (make-instance 'flow-exit-actor))

 (flow-link flow 'ok 'start "done")
 (flow-set-entry flow 'start)
 (flow-submit flow (make-flow-box :set nil :watch 'done))
 (channel-pop done))`,
		Expect: "/#<flow-box [0-9a-f]+>/",
	}).Test(t)
	fi, err := os.Stat(filename)
	tt.Nil(t, err)
	tt.Equal(t, "-r--r--r--", fi.Mode().String())
	content, _ := os.ReadFile(filename)
	tt.Equal(t, "Sample file", string(content))
}

func TestWriteFileActorNoFile(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (flow-add-task flow
                :name "start"
                :actor (make-instance 'flow-write-file-actor
                                      :filename "testdata/nothing/no-file.txt"
                                      :content "text"))
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
		Expect: `"{content: null error: "open testdata/nothing/no-file.txt: no such file or directory"}"`,
	}).Test(t)
}

func TestWriteFileActorBadPermissions(t *testing.T) {
	(&sliptest.Function{
		Source: `(make-instance 'flow-write-file-actor
                                :filename "testdata/write-me.text"
                                :content "Sample content"
                                :permissions 'rw-r--r)`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
	(&sliptest.Function{
		Source: `(make-instance 'flow-write-file-actor
                                :filename "testdata/write-me.text"
                                :content "Sample content"
                                :permissions t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
	(&sliptest.Function{
		Source: `(make-instance 'flow-write-file-actor
                                :filename "testdata/write-me.text"
                                :content "Sample content"
                                :permissions "-rz-r--r--")`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}

func TestWriteFileActorInitKeyValues(t *testing.T) {
	(&sliptest.Function{
		Source: `(send (make-instance 'flow-write-file-actor
                                      :filename "testdata/write-me.txt"
                                      :content "Sample write"
                                      :append t
                                      :overwrite t
                                      :permissions #o660) :init-key-values)`,
		Validate: func(t *testing.T, v slip.Object) {
			for i, x := range []string{
				`:filename`,
				`"testdata/write-me.txt"`,
				`:content`,
				`"Sample write"`,
				`:permissions`,
				`"-rw-rw----"`,
				`:overwrite`,
				`t`,
				`:append`,
				`t`,
			} {
				tt.Equal(t, x, slip.ObjectString(v.(slip.List)[i]), x)
			}
		},
	}).Test(t)
	(&sliptest.Function{
		Source: `(send (make-instance 'flow-write-file-actor
                                      :filename "testdata/write-me.txt"
                                      :content "Sample write"
                                      :append nil
                                      :overwrite nil
                                      :permissions 'rw-r--r--) :init-key-values)`,
		Validate: func(t *testing.T, v slip.Object) {
			for i, x := range []string{
				`:filename`,
				`"testdata/write-me.txt"`,
				`:content`,
				`"Sample write"`,
				`:permissions`,
				`"-rw-r--r--"`,
				`:overwrite`,
				`nil`,
				`:append`,
				`nil`,
			} {
				tt.Equal(t, x, slip.ObjectString(v.(slip.List)[i]), x)
			}
		},
	}).Test(t)
}

func TestWriteFileActorDocs(t *testing.T) {
	scope := slip.NewScope()
	var out strings.Builder
	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})

	for _, method := range []string{
		":init",
		":start",
		":perform",
		":init-key-values",
	} {
		_ = slip.ReadString(fmt.Sprintf(`(describe-method flow-write-file-actor %s out)`, method), scope).Eval(scope, nil)
		tt.Equal(t, true, strings.Contains(out.String(), method))
		out.Reset()
	}
}
