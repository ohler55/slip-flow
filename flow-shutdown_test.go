// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

// Partially tested in flow-start_test.go.

func TestFlowShutdownNotFlow(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-shutdown t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowShutdownArgCount(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-shutdown (make-instance 'flow-flavor :name 'flo) t)`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}
