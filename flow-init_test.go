// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestFlowInitBadName(t *testing.T) {
	(&sliptest.Function{
		Source:    `(make-instance 'flow-flavor :name t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowInitBadExitChannel(t *testing.T) {
	(&sliptest.Function{
		Source:    `(make-instance 'flow-flavor :exit-channel t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
