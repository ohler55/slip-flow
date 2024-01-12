// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestBoxBagMethod(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :bag))`,
		Expect: `/#<bag-flavor [0-9a-f]+>/`,
	}).Test(t)
}

func TestBoxBagFrozen(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :freeze)
                  (send box :bag))`,
		Expect: `/#<bag-flavor [0-9a-f]+>/`,
	}).Test(t)
}

func TestBoxBagFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-bag box))`,
		Expect: `/#<bag-flavor [0-9a-f]+>/`,
	}).Test(t)
}

func TestBoxBagNotBox(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-box-bag t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
