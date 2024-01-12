// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestBoxHas(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-has box "x"))`,
		Expect: "t",
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :has (make-bag-path "x")))`,
		Expect: "t",
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :has "y"))`,
		Expect: "nil",
	}).Test(t)
}

func TestBoxHasNotBox(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-box-has t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestBoxHasArgCount(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :has "x" t))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}

func TestBoxHasBadPath(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :has t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
