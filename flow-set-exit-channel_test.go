// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestFlowSetExitChannelFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `(flow-set-exit-channel (make-instance 'flow-flavor) (make-channel 3))`,
		Expect: "#<channel 3>",
	}).Test(t)
	(&sliptest.Function{
		Source: `(flow-set-exit-channel (make-instance 'flow-flavor) nil)`,
		Expect: "nil",
	}).Test(t)
}

func TestFlowSetExitChannelSend(t *testing.T) {
	(&sliptest.Function{
		Source: `(send (make-instance 'flow-flavor) :set-exit-channel (make-channel 3))`,
		Expect: "#<channel 3>",
	}).Test(t)
	(&sliptest.Function{
		Source: `(send (make-instance 'flow-flavor) :set-exit-channel nil)`,
		Expect: "nil",
	}).Test(t)
}

func TestFlowSetExitChannelNotFlow(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-set-exit-channel t nil)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestFlowSetExitChannelArgCount(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-set-exit-channel (make-instance 'flow-flavor))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}

func TestFlowSetExitChannelNotChannel(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-set-exit-channel (make-instance 'flow-flavor) t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
	(&sliptest.Function{
		Source:    `(send (make-instance 'flow-flavor) :set-exit-channel t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
