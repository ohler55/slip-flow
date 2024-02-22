// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

// Mostly tested in box-watch_test.go

func TestBoxNotifyNotBox(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-box-notify 7 'test)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestBoxNotifyBadName(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-box-notify (make-flow-box) 7)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
