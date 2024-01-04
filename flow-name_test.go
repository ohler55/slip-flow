// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestFlowNameOk(t *testing.T) {
	(&sliptest.Function{
		Source: `(flow-name (make-instance 'flow-flavor :name 'flo))`,
		Expect: `"flo"`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(send (make-instance 'flow-flavor :name "flo") :name)`,
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
		Source:    `(flow-name (make-instance 'flow-flavor :name 'flo) t)`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}
