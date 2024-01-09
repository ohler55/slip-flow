// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

// Already tested in flow-link tests.

func TestTaskLinksNotTask(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-task-links t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
