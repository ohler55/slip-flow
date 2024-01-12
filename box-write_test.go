// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"fmt"
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestBoxWriteString(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-write box :json t))`,
		Expect: `"{"x": 3}"`,
	}).Test(t)
}

func TestBoxWriteStream(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}"))
                       (out (make-string-output-stream)))
                  (send box :write out)
                  (get-output-stream-string out))`,
		Expect: `"{x: 3}"`,
	}).Test(t)
}

func TestBoxWriteStdOut(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}"))
                       (*standard-output* (make-string-output-stream)))
                  (send box :write t)
                  (get-output-stream-string *standard-output*))`,
		Expect: `"{x: 3}"`,
	}).Test(t)
}

func TestBoxWriteSEN(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-write box :json nil :pretty nil :indent 0))`,
		Expect: `"{x:3}"`,
	}).Test(t)
}

func TestBoxWriteJSON(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-write box :json t :pretty nil :depth 0 :time-format nil :time-wrap nil))`,
		Expect: `"{"x":3}"`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-write box :time-format "2006-01-02" :time-wrap "@"))`,
		Expect: `"{x: 3}"`,
	}).Test(t)
}

func TestBoxWriteFull(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :tracking-id 1234 :parse "{x:3}")))
                  (send box :scan "flo" "tisk")
                  (flow-box-write box :full t))`,
		Expect: `/"{
  content: {x: 3}
  track: {history: \[{flow: flo task: tisk when: [0-9]+}\] id: 1234}
}"/`,
	}).Test(t)
}

func TestBoxWriteOptions(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-write box :color nil :right-margin 80 :depth 4))`,
		Expect: `"{x: 3}"`,
	}).Test(t)
}

func TestBoxWriteNotBox(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-box-write t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestBoxWriteNotStream(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :write 7))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestBoxWriteBadKeyword(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :write :bad t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
	                  (send box :write :right-margin t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
	                  (send box :write :depth t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
	                  (send box :write :indent t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
	                  (send box :write :time-format t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
	                  (send box :write :time-wrap t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

type badWriter int

func (w badWriter) Write([]byte) (int, error) {
	return 0, fmt.Errorf("oops")
}

func TestBoxWriteOutputError(t *testing.T) {
	scope := slip.NewScope()
	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: badWriter(0)})
	(&sliptest.Function{
		Scope: scope,
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :write out))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}
