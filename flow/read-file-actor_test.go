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

func TestReadFileActorTextStatic(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (flow-add-task flow
                :name "start"
                :actor (make-instance 'flow-read-file-actor
                                      :filename "testdata/read-me.txt"
                                      :destination "content"))
 (flow-add-task flow
                :name "done"
                :actor (make-instance 'flow-exit-actor))

 (flow-link flow 'ok 'start "done")
 (flow-set-entry flow 'start)
 (flow-submit flow (make-flow-box :set nil :watch 'done))
 (flow-box-write (channel-pop done)))`,
		Expect: `"{content: "This is just
text in a file.
"}"`,
	}).Test(t)
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (send flow :add-task
            :name "start"
            :actor (make-instance 'flow-read-file-actor
                                  :filename "testdata/read-me.txt"
                                  :destination "content"))
 (send flow :add-task
            :name "done"
            :actor (make-instance 'flow-exit-actor))

 (send flow :link 'ok 'start "done")
 (send flow :set-entry 'start)
 (send flow :submit (make-flow-box :set nil :watch 'done))
 (send (channel-pop done) :write))`,
		Expect: `"{content: "This is just
text in a file.
"}"`,
	}).Test(t)
}

func TestReadFileActorTextFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (flow-add-task flow
                :name "start"
                :actor (make-instance 'flow-read-file-actor
                                      :filename (lambda (b) "testdata/read-me.txt")
                                      :destination "content"))
 (flow-add-task flow
                :name "done"
                :actor (make-instance 'flow-exit-actor))

 (flow-link flow 'ok 'start "done")
 (flow-set-entry flow 'start)
 (flow-submit flow (make-flow-box :set nil :watch 'done))
 (flow-box-write (channel-pop done)))`,
		Expect: `"{content: "This is just
text in a file.
"}"`,
	}).Test(t)
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (send flow :add-task
            :name "start"
            :actor (make-instance 'flow-read-file-actor
                                  :filename (lambda (b) "testdata/read-me.txt")
                                  :destination "content"))
 (send flow :add-task
            :name "done"
            :actor (make-instance 'flow-exit-actor))

 (send flow :link 'ok 'start "done")
 (send flow :set-entry 'start)
 (send flow :submit (make-flow-box :set nil :watch 'done))
 (send (channel-pop done) :write))`,
		Expect: `"{content: "This is just
text in a file.
"}"`,
	}).Test(t)
}

func TestReadFileActorNoFile(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (flow-add-task flow
                :name "start"
                :actor (make-instance 'flow-read-file-actor
                                      :filename "testdata/no-file.txt"
                                      :destination "text"))
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
		Expect: `"{content: null error: "open testdata/no-file.txt: no such file or directory"}"`,
	}).Test(t)
}

func TestReadFileActorInitKeyValues(t *testing.T) {
	(&sliptest.Function{
		Source: `(send (make-instance 'flow-read-file-actor
                                      :filename "testdata/read-me.txt"
                                      :destination 'content) :init-key-values)`,
		Validate: func(t *testing.T, v slip.Object) {
			for i, x := range []string{
				`:filename`,
				`"testdata/read-me.txt"`,
				`:destination`,
				`"content"`,
			} {
				tt.Equal(t, x, slip.ObjectString(v.(slip.List)[i]), x)
			}
		},
	}).Test(t)
}

func TestReadFileActorBadDestination(t *testing.T) {
	(&sliptest.Function{
		Source: `(make-instance 'flow-read-file-actor
                                :filename "testdata/read-me.txt"
                                :destination t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestReadFileActorDocs(t *testing.T) {
	scope := slip.NewScope()
	var out strings.Builder
	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})

	for _, method := range []string{
		":init",
		":start",
		":perform",
		":init-key-values",
	} {
		_ = slip.ReadString(fmt.Sprintf(`(describe-method flow-read-file-actor %s out)`, method), scope).Eval(scope, nil)
		tt.Equal(t, true, strings.Contains(out.String(), method))
		out.Reset()
	}
}
