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

func TestGlobActorNames(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (flow-add-task flow
                :name "start"
                :actor (make-instance 'flow-glob-actor
                                      :pattern "testdata/*.csv"
                                      :destination "paths"))
 (flow-add-task flow
                :name "done"
                :actor (make-instance 'flow-exit-actor))

 (flow-link flow 'ok 'start "done")
 (flow-set-entry flow 'start)
 (flow-submit flow (make-flow-box :set nil :watch 'done))
 (flow-box-write (channel-pop done)))`,
		Expect: `"{paths: ["testdata/bad.csv" "testdata/read-me.csv"]}"`,
	}).Test(t)
}

func TestGlobActorInfo(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (flow-add-task flow
                :name "start"
                :actor (make-instance 'flow-glob-actor
                                      :pattern "testdata/*.csv"
                                      :with-info t
                                      :destination "paths"))
 (flow-add-task flow
                :name "done"
                :actor (make-instance 'flow-exit-actor))

 (flow-link flow 'ok 'start "done")
 (flow-set-entry flow 'start)
 (flow-submit flow (make-flow-box :set nil :watch 'done))
 (flow-box-write (channel-pop done) :time-format "2006-01-02T15:04:05Z07:00"))`,
		Expect: `/"{
  paths: \[
    {is-dir: false mode: -rw-r--r-- modified-time: ".*Z" name: "testdata/bad.csv" size: 23}
    {is-dir: false mode: -rw-r--r-- modified-time: ".*Z" name: "testdata/read-me.csv" size: 30}
  \]
}"/`,
	}).Test(t)
}

func TestGlobActorBadPattern(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (flow-add-task flow
                :name "start"
                :actor (make-instance 'flow-glob-actor
                                      :pattern "testdata/[]"
                                      :destination "dir"))
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
		Expect: `"{content: null error: "syntax error in pattern"}"`,
	}).Test(t)
}

func TestGlobActorDocs(t *testing.T) {
	scope := slip.NewScope()
	var out strings.Builder
	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})

	for _, method := range []string{
		":init",
		":start",
		":perform",
		":init-key-values",
	} {
		_ = slip.ReadString(fmt.Sprintf(`(describe-method flow-glob-actor %s out)`, method)).Eval(scope, nil)
		tt.Equal(t, true, strings.Contains(out.String(), method))
		out.Reset()
	}
}

func TestGlobActorInitKeyValues(t *testing.T) {
	(&sliptest.Function{
		Source: `(send (make-instance 'flow-glob-actor
                                      :pattern "testdata/*.csv"
                                      :destination 'dir
                                      :with-info t) :init-key-values)`,
		Validate: func(t *testing.T, v slip.Object) {
			for i, x := range []string{
				`:pattern`,
				`"testdata/*.csv"`,
				`:destination`,
				`"dir"`,
				`:with-info`,
				`t`,
			} {
				tt.Equal(t, x, slip.ObjectString(v.(slip.List)[i]), x)
			}
		},
	}).Test(t)
	(&sliptest.Function{
		Source: `(send (make-instance 'flow-glob-actor
                                      :pattern "testdata/*.csv"
                                      :destination 'dir
                                      :with-info nil) :init-key-values)`,
		Validate: func(t *testing.T, v slip.Object) {
			for i, x := range []string{
				`:pattern`,
				`"testdata/*.csv"`,
				`:destination`,
				`"dir"`,
				`:with-info`,
				`nil`,
			} {
				tt.Equal(t, x, slip.ObjectString(v.(slip.List)[i]), x)
			}
		},
	}).Test(t)
}
