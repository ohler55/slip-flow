// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

// Also tested in group-start_test

func TestGroupAddNotGroup(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-group-add t (make-flow :name 'flo))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
func TestGroupAddNotFlow(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-group-add (make-flow-group) t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
