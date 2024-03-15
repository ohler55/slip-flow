// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestBoxTrackingIDMethod(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :tracking-id 1234 :parse "{x:3}")))
                  (send box :tracking-id))`,
		Expect: "1234",
	}).Test(t)
}

func TestBoxTrackingIDFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :tracking-id 1234 :parse "{x:3}")))
                  (flow-box-tracking-id box))`,
		Expect: "1234",
	}).Test(t)
	(&sliptest.Function{
		Source:    `(flow-box-tracking-id t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
