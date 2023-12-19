// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestBoxHistoryMethod(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :tracking-id 1234 :parse "{x:3}")))
                  (send box :scan "flo" "tisk")
                  (send box :history))`,
		Expect: `/\(\(@\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\d.[\d]+Z "tisk" "flo"\)\)/`,
	}).Test(t)
}

func TestBoxHistoryFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :tracking-id 1234 :parse "{x:3}")))
                  (send box :scan "flo" "tisk")
                  (flow-box-history box))`,
		Expect: `/\(\(@\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\d.[\d]+Z "tisk" "flo"\)\)/`,
	}).Test(t)
	(&sliptest.Function{
		Source:    `(flow-box-history t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
