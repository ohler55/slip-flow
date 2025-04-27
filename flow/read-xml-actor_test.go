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

func TestReadXMLActorOkay(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (flow-add-task flow
                :name "start"
                :actor (make-instance 'flow-read-xml-actor
                                      :filename "testdata/read-me.xml"
                                      :destination "content"
                                      :count "count"
                                      :trim t))
 (flow-add-task flow
                :name "done"
                :actor (make-instance 'flow-exit-actor))

 (flow-link flow 'ok 'start "done")
 (flow-set-entry flow 'start)
 (flow-submit flow (make-flow-box :set nil :watch 'done))
 (list (flow-box-write (channel-pop done)) (flow-box-write (channel-pop done))))`,
		Validate: func(t *testing.T, v slip.Object) {
			list := v.(slip.List)
			for i, x := range []string{
				`{
  content: [
    [":processing-instruction" xml "version=\"1.0\""]
    [":directive" "DOCTYPE sample PUBLIC \"sample.dtd\""]
    [top {id: "123"} [child {} "Some text."] [":comment" "a comment"] [blank {}]]
  ]
  count: 0
}`,
				`{
  content: [[":processing-instruction" xml "version=\"1.0\""] [top {id: "321"} [child {} "More text."]]]
  count: 1
}`,
			} {
				tt.Equal(t, x, string(list[i].(slip.String)))
			}
		},
	}).Test(t)
}

func TestReadXMLActorHTML(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (flow-add-task flow
                :name "start"
                :actor (make-instance 'flow-read-xml-actor
                                      :filename "testdata/read-me.html"
                                      :destination "content"
                                      :html t
                                      :trim t))
 (flow-add-task flow
                :name "done"
                :actor (make-instance 'flow-exit-actor))

 (flow-link flow 'ok 'start "done")
 (flow-set-entry flow 'start)
 (flow-submit flow (make-flow-box :set nil :watch 'done))
 (flow-box-write (channel-pop done)))`,
		Expect: `"{
  content: [[":directive" "DOCTYPE html"] [html {} [body {} "The HTML body."]]]
}"`,
	}).Test(t)
}

func TestReadXMLActorNoFile(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (flow-add-task flow
                :name "start"
                :actor (make-instance 'flow-read-xml-actor
                                      :filename "testdata/no-file.xml"
                                      :destination "xml"))
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
		Expect: `"{content: null error: "open testdata/no-file.xml: no such file or directory"}"`,
	}).Test(t)
}

func TestReadXMLActorBadFile(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (flow-add-task flow
                :name "start"
                :actor (make-instance 'flow-read-xml-actor
                                      :filename "testdata/bad.xml"
                                      :destination "xml"))
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
		Expect: `"{content: null error: "XML syntax error on line 3: expected attribute name in element"}"`,
	}).Test(t)
}

func TestReadXMLActorDocs(t *testing.T) {
	scope := slip.NewScope()
	var out strings.Builder
	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})

	for _, method := range []string{
		":init",
		":start",
		":perform",
		":init-key-values",
	} {
		_ = slip.ReadString(fmt.Sprintf(`(describe-method flow-read-xml-actor %s out)`, method), scope).Eval(scope, nil)
		tt.Equal(t, true, strings.Contains(out.String(), method))
		out.Reset()
	}
}

func TestReadXMLActorInitKeyValues(t *testing.T) {
	(&sliptest.Function{
		Source: `(send (make-instance 'flow-read-xml-actor
                                      :filename "testdata/read-me.xml"
                                      :destination 'content
                                      :trim t
                                      :strict t
                                      :html nil
                                      :count 'count) :init-key-values)`,
		Validate: func(t *testing.T, v slip.Object) {
			for i, x := range []string{
				`:filename`,
				`"testdata/read-me.xml"`,
				`:destination`,
				`"content"`,
				`:count`,
				`"count"`,
				`:strict`,
				`t`,
				`:trim`,
				`t`,
				`:html`,
				`nil`,
			} {
				tt.Equal(t, x, slip.ObjectString(v.(slip.List)[i]), x)
			}
		},
	}).Test(t)
	(&sliptest.Function{
		Source: `(send (make-instance 'flow-read-xml-actor
                                      :filename "testdata/read-me.xml"
                                      :destination 'content
                                      :trim nil
                                      :strict nil
                                      :html t
                                      :count 'count) :init-key-values)`,
		Validate: func(t *testing.T, v slip.Object) {
			for i, x := range []string{
				`:filename`,
				`"testdata/read-me.xml"`,
				`:destination`,
				`"content"`,
				`:count`,
				`"count"`,
				`:strict`,
				`nil`,
				`:trim`,
				`nil`,
				`:html`,
				`t`,
			} {
				tt.Equal(t, x, slip.ObjectString(v.(slip.List)[i]), x)
			}
		},
	}).Test(t)
}
