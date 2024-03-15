// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestBoxNativeMethod(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :native))`,
		Expect: `(("x" . 3))`,
	}).Test(t)
}

func TestBoxNativeFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-native box))`,
		Expect: `(("x" . 3))`,
	}).Test(t)
}

func TestBoxNativeNotBox(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-box-native t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestBoxNativeArgCount(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-native box t))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}
