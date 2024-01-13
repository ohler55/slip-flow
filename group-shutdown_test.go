// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

// Also tested in group-start_test

func TestGroupShutdownNotGroup(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-group-shutdown t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
