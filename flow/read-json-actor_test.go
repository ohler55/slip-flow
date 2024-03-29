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

func TestReadJSONActorStatic(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (flow-add-task flow
                :name "start"
                :actor (make-instance 'flow-read-json-actor
                                      :filename "testdata/read-me.json"
                                      :destination "content"
                                      :count "count"))
 (flow-add-task flow
                :name "done"
                :actor (make-instance 'flow-exit-actor))

 (flow-link flow 'ok 'start "done")
 (flow-set-entry flow 'start)
 (flow-submit flow (make-flow-box :set nil :watch 'done))
 (list (flow-box-write (channel-pop done)) (flow-box-write (channel-pop done))))`,
		Expect: `("{content: {quux: "this is json"} count: 0}"
 "{content: {quack: "like a duck"} count: 1}")`,
	}).Test(t)
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (send flow :add-task
            :name "start"
            :actor (make-instance 'flow-read-json-actor
                                  :filename "testdata/read-me.sen"
                                  :destination "content"
                                  :count "count"))
 (send flow :add-task
            :name "done"
            :actor (make-instance 'flow-exit-actor))

 (send flow :link 'ok 'start "done")
 (send flow :set-entry 'start)
 (send flow :submit (make-flow-box :set nil :watch 'done))
 (list (send (channel-pop done) :write) (send (channel-pop done) :write)))`,
		Expect: `("{content: {quux: "this is json"} count: 0}"
 "{content: {quack: "like a duck"} count: 1}")`,
	}).Test(t)
}

func TestReadJSONActorFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (flow-add-task flow
                :name "start"
                :actor (make-instance 'flow-read-json-actor
                                      :filename (lambda (b) "testdata/read-me.json")
                                      :destination "content"
                                      :count "count"))
 (flow-add-task flow
                :name "done"
                :actor (make-instance 'flow-exit-actor))

 (flow-link flow 'ok 'start "done")
 (flow-set-entry flow 'start)
 (flow-submit flow (make-flow-box :set nil :watch 'done))
 (list (flow-box-write (channel-pop done)) (flow-box-write (channel-pop done))))`,
		Expect: `("{content: {quux: "this is json"} count: 0}"
 "{content: {quack: "like a duck"} count: 1}")`,
	}).Test(t)
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (send flow :add-task
            :name "start"
            :actor (make-instance 'flow-read-json-actor
                                  :filename (lambda (b) "testdata/read-me.json")
                                  :destination 'content
                                  :count 'count))
 (send flow :add-task
            :name "done"
            :actor (make-instance 'flow-exit-actor))

 (send flow :link 'ok 'start "done")
 (send flow :set-entry 'start)
 (send flow :submit (make-flow-box :set nil :watch 'done))
 (list (send (channel-pop done) :write) (send (channel-pop done) :write)))`,
		Expect: `("{content: {quux: "this is json"} count: 0}"
 "{content: {quack: "like a duck"} count: 1}")`,
	}).Test(t)
}

func TestReadJSONActorNoFile(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (flow-add-task flow
                :name "start"
                :actor (make-instance 'flow-read-json-actor
                                      :filename "testdata/no-file.json"
                                      :destination "json"))
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
		Expect: `"{content: null error: "open testdata/no-file.json: no such file or directory"}"`,
	}).Test(t)
}

func TestReadJSONActorBadDestination(t *testing.T) {
	(&sliptest.Function{
		Source: `(make-instance 'flow-read-json-actor
                                :filename "testdata/read-me.sen"
                                :destination t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestReadJSONActorBadCount(t *testing.T) {
	(&sliptest.Function{
		Source: `(make-instance 'flow-read-json-actor
                                :filename "testdata/read-me.sen"
                                :destination 'json
                                :count t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestReadJSONActorDocs(t *testing.T) {
	scope := slip.NewScope()
	var out strings.Builder
	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})

	for _, method := range []string{
		":init",
		":start",
		":perform",
		":init-key-values",
	} {
		_ = slip.ReadString(fmt.Sprintf(`(describe-method flow-read-json-actor %s out)`, method)).Eval(scope, nil)
		tt.Equal(t, true, strings.Contains(out.String(), method))
		out.Reset()
	}
}

func TestReadJSONActorInitKeyValues(t *testing.T) {
	(&sliptest.Function{
		Source: `(send (make-instance 'flow-read-json-actor
                                      :filename "testdata/read-me.txt"
                                      :destination 'content
                                      :count 'count) :init-key-values)`,
		Validate: func(t *testing.T, v slip.Object) {
			for i, x := range []string{
				`:filename`,
				`"testdata/read-me.txt"`,
				`:destination`,
				`"content"`,
				`:count`,
				`"count"`,
			} {
				tt.Equal(t, x, slip.ObjectString(v.(slip.List)[i]), x)
			}
		},
	}).Test(t)
}
