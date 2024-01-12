// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestBoxThawMethod(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-freeze box)
                  (send box :thaw)
                  (send box :frozen))`,
		Expect: "nil",
	}).Test(t)
}

func TestBoxThawFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-freeze box)
                  (flow-box-thaw box)
                  (send box :frozen))`,
		Expect: "nil",
	}).Test(t)
	(&sliptest.Function{
		Source:    `(flow-box-thaw t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
