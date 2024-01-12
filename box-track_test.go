// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestBoxTrackMethod(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :tracking-id 1234 :parse "{x:3}")))
                  (send box :scan "flo" "tisk")
                  (list (send (flow-box-track box) :id) (send (flow-box-track box) :history)))`,
		Expect: `/\(1234 \(\(@\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\d.[\d]+Z "tisk" "flo"\)\)\)/`,
	}).Test(t)
}

func TestBoxTrackFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :tracking-id 1234 :parse "{x:3}")))
                  (send box :scan "flo" "tisk")
                  (list (send (flow-box-track box) :id) (send (flow-box-track box) :history)))`,
		Expect: `/\(1234 \(\(@\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\d.[\d]+Z "tisk" "flo"\)\)\)/`,
	}).Test(t)
	(&sliptest.Function{
		Source:    `(flow-box-track t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
