// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestBoxParsePath(t *testing.T) {
	scope := slip.NewScope()
	orig := slip.ReadString(`*flow-box-time-format*`).Eval(scope, nil)
	defer scope.Set("*flow-box-time-format*", orig)
	_ = slip.ReadString(`(setq *flow-box-time-format* *rfc3339nano*)`).Eval(scope, nil)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-parse box "[7]" "x")
                  (send box :native))`,
		Expect: `(("x" . (7)))`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :parse "[7]" (make-bag-path "x"))
                  (send box :native))`,
		Expect: `(("x" . (7)))`,
	}).Test(t)
}

func TestBoxParseAll(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-parse box "[7]")
                  (send box :native))`,
		Expect: `(7)`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :parse "[7]")
                  (send box :native))`,
		Expect: `(7)`,
	}).Test(t)
}

func TestBoxParseFrozen(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :freeze)
                  (flow-box-parse box "[7]" "x")
                  (list (send box :native) (send box :frozen)))`,
		Expect: `((("x" . (7))) nil)`,
	}).Test(t)
}

func TestBoxParseNotBox(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-box-parse t "7" "x")`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestBoxParseBadPath(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-parse box "7" t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestBoxParseBadString(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-parse box 7 "x"))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestBoxParseBadArgCount(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :parse))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}
