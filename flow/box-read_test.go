// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestBoxReadPath(t *testing.T) {
	scope := slip.NewScope()
	orig := slip.ReadString(`*flow-box-time-format*`).Eval(scope, nil)
	defer scope.Set("*flow-box-time-format*", orig)
	_ = slip.ReadString(`(setq *flow-box-time-format* *rfc3339nano*)`).Eval(scope, nil)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-read box (make-string-input-stream "[7]") "x")
                  (send box :native))`,
		Expect: `(("x" . (7)))`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :read (make-string-input-stream "[7]") (make-bag-path "x"))
                  (send box :native))`,
		Expect: `(("x" . (7)))`,
	}).Test(t)
}

func TestBoxReadAll(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-read box (make-string-input-stream "[7]"))
                  (send box :native))`,
		Expect: `(7)`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :read (make-string-input-stream "[7]"))
                  (send box :native))`,
		Expect: `(7)`,
	}).Test(t)
}

func TestBoxReadFrozen(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :freeze)
                  (flow-box-read box (make-string-input-stream "[7]") "x")
                  (list (send box :native) (send box :frozen)))`,
		Expect: `((("x" . (7))) nil)`,
	}).Test(t)
}

func TestBoxReadNotBox(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-box-read t (make-string-input-stream "7") "x")`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestBoxReadBadPath(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-read box (make-string-input-stream "7") t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestBoxReadBadStream(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-read box 7 "x"))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestBoxReadBadArgCount(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :read))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}
