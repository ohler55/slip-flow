// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestFlowNameOk(t *testing.T) {
	(&sliptest.Function{
		Source: `(flow-name (make-instance 'flow :name 'flo))`,
		Expect: `"flo"`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(send (make-instance 'flow :name "flo") :name)`,
		Expect: `"flo"`,
	}).Test(t)
}

func TestFlowNameNotFlow(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-name t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowNameArgCount(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-name (make-instance 'flow :name 'flo) t)`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}
