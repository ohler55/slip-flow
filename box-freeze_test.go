// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestBoxFreezeMethod(t *testing.T) {
	(&sliptest.Function{
		Source: `(let* ((box (make-flow-box :parse "{x:3}"))
                        (result (list (send box :frozen))))
                  (send box :freeze)
                  (add result (send box :frozen)))`,
		Expect: "(nil t)",
	}).Test(t)
}

func TestBoxFreezeFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `(let* ((box (make-flow-box :parse "{x:3}"))
                        (result (list (send box :frozen))))
                  (flow-box-freeze box)
                  (add result (send box :frozen)))`,
		Expect: "(nil t)",
	}).Test(t)
	(&sliptest.Function{
		Source:    `(flow-box-freeze t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
