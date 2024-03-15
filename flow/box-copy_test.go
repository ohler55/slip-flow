// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestBoxCopyMethod(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :scan "flo" "tisk")
                  (send (send box :copy) :native))`,
		Expect: `(("x" . 3))`,
	}).Test(t)
}

func TestBoxCopyFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send (flow-box-copy box) :native))`,
		Expect: `(("x" . 3))`,
	}).Test(t)
}

func TestBoxCopyNotBox(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-box-copy t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
