// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestFlowExitChannelOk(t *testing.T) {
	(&sliptest.Function{
		Source: `(flow-exit-channel (make-instance 'flow :exit-channel (make-channel 3)))`,
		Expect: "#<channel 3>",
	}).Test(t)
	(&sliptest.Function{
		Source: `(send (make-instance 'flow :exit-channel (make-channel 3)) :exit-channel)`,
		Expect: "#<channel 3>",
	}).Test(t)
}

func TestFlowExitChannelNotFlow(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-exit-channel t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowExitChannelArgCount(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-exit-channel (make-instance 'flow :exit-channel (make-channel 3)) t)`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}
