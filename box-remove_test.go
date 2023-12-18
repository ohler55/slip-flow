// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestBoxRemoveOk(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3 y:4}")))
                  (flow-box-remove box "x")
                  (send box :native))`,
		Expect: `(("y" . 4))`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3 y:4}")))
                  (send box :remove (make-bag-path "x"))
                  (send box :native))`,
		Expect: `(("y" . 4))`,
	}).Test(t)
}

func TestBoxRemoveFrozen(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3 y:4}")))
                  (send box :freeze)
                  (flow-box-remove box "x")
                  (send box :native))`,
		Expect: `(("y" . 4))`,
	}).Test(t)
}

func TestBoxRemoveNil(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-remove box nil)
                  (send box :native))`,
		Expect: "nil",
	}).Test(t)
}

func TestBoxRemoveNotBox(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-box-remove t "x")`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestBoxRemoveArgCount(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3 y:4}")))
                  (flow-box-remove box "x" t))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}

func TestBoxRemoveBadPath(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3 y:4}")))
                  (send box :remove t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
