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

func TestReadCSVActorList(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (flow-add-task flow
                :name "start"
                :actor (make-instance 'flow-read-csv-actor
                                      :filename "testdata/read-me.csv"
                                      :destination "content"
                                      :count "count"))
 (flow-add-task flow
                :name "done"
                :actor (make-instance 'flow-exit-actor))

 (flow-link flow 'ok 'start "done")
 (flow-set-entry flow 'start)
 (flow-submit flow (make-flow-box :set nil :watch 'done))
 (list (flow-box-write (channel-pop done)) (flow-box-write (channel-pop done)) (flow-box-write (channel-pop done))))`,
		Validate: func(t *testing.T, v slip.Object) {
			list := v.(slip.List)
			for i, x := range []string{
				`{content: [first last] count: 0}`,
				`{content: [foo bar] count: 1}`,
				`{content: [quux ducks] count: 2}`,
			} {
				tt.Equal(t, x, string(list[i].(slip.String)))
			}
		},
	}).Test(t)
}

func TestReadCSVActorMap(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (flow-add-task flow
                :name "start"
                :actor (make-instance 'flow-read-csv-actor
                                      :filename "testdata/read-me.csv"
                                      :destination "content"
                                      :as-map t))
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
				`{content: {first: foo last: bar}}`,
				`{content: {first: quux last: ducks}}`,
			} {
				tt.Equal(t, x, string(list[i].(slip.String)))
			}
		},
	}).Test(t)
}

func TestReadCSVActorNoFile(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (flow-add-task flow
                :name "start"
                :actor (make-instance 'flow-read-csv-actor
                                      :filename "testdata/no-file.csv"
                                      :destination "csv"))
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
		Expect: `"{content: null error: "open testdata/no-file.csv: no such file or directory"}"`,
	}).Test(t)
}

func TestReadCSVActorBadFile(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let ((done (make-channel 3))
      (flow (make-flow :name 'flo)))
 (flow-add-task flow
                :name "start"
                :actor (make-instance 'flow-read-csv-actor
                                      :filename "testdata/bad.csv"
                                      :destination "csv"
                                      :as-map t))
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
		Expect: `"{content: null error: "record on line 2: wrong number of fields"}"`,
	}).Test(t)
}

func TestReadCSVActorBadSeparator(t *testing.T) {
	(&sliptest.Function{
		Source: `(make-instance 'flow-read-csv-actor
                                :filename "testdata/read-me.sen"
                                :destination "csv"
                                :separator t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestReadCSVActorBadComment(t *testing.T) {
	(&sliptest.Function{
		Source: `(make-instance 'flow-read-csv-actor
                                :filename "testdata/read-me.sen"
                                :destination "csv"
                                :comment t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestReadCSVActorDocs(t *testing.T) {
	scope := slip.NewScope()
	var out strings.Builder
	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})

	for _, method := range []string{
		":init",
		":start",
		":perform",
		":init-key-values",
	} {
		_ = slip.ReadString(fmt.Sprintf(`(describe-method flow-read-csv-actor %s out)`, method), scope).Eval(scope, nil)
		tt.Equal(t, true, strings.Contains(out.String(), method))
		out.Reset()
	}
}

func TestReadCSVActorInitKeyValues(t *testing.T) {
	(&sliptest.Function{
		Source: `(send (make-instance 'flow-read-csv-actor
                                      :filename "testdata/read-me.csv"
                                      :destination 'content
                                      :separator #\Tab
                                      :comment #\#
                                      :trim t
                                      :as-map t
                                      :count 'count) :init-key-values)`,
		Validate: func(t *testing.T, v slip.Object) {
			for i, x := range []string{
				`:filename`,
				`"testdata/read-me.csv"`,
				`:destination`,
				`"content"`,
				`:count`,
				`"count"`,
				`:separator`,
				`#\Tab`,
				`:comment`,
				`#\#`,
				`:trim`,
				`t`,
				`:as-map`,
				`t`,
			} {
				tt.Equal(t, x, slip.ObjectString(v.(slip.List)[i]), x)
			}
		},
	}).Test(t)
	(&sliptest.Function{
		Source: `(send (make-instance 'flow-read-csv-actor
                                      :filename "testdata/read-me.csv"
                                      :destination 'content
                                      :separator #\Tab
                                      :comment #\#
                                      :trim nil
                                      :as-map nil
                                      :count 'count) :init-key-values)`,
		Validate: func(t *testing.T, v slip.Object) {
			for i, x := range []string{
				`:filename`,
				`"testdata/read-me.csv"`,
				`:destination`,
				`"content"`,
				`:count`,
				`"count"`,
				`:separator`,
				`#\Tab`,
				`:comment`,
				`#\#`,
				`:trim`,
				`nil`,
				`:as-map`,
				`nil`,
			} {
				tt.Equal(t, x, slip.ObjectString(v.(slip.List)[i]), x)
			}
		},
	}).Test(t)
}
