// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestBoxFrozenMethod(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :frozen))`,
		Expect: "nil",
	}).Test(t)
}

func TestBoxFrozenFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-frozen box))`,
		Expect: "nil",
	}).Test(t)
	(&sliptest.Function{
		Source:    `(flow-box-frozen t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
