// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestFlowInitBadName(t *testing.T) {
	(&sliptest.Function{
		Source:    `(make-instance 'flow :name t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowInitBadKeyword(t *testing.T) {
	(&sliptest.Function{
		Source:    `(make-instance 'flow t t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
