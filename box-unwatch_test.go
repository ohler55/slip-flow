// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestBoxUnwatchFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let* ((chan (make-channel 3))
       (box (make-flow-box :parse "{x:3}")))
 (flow-box-watch box "test" chan)
 (flow-box-unwatch box "test")
 (flow-box-notify box)
 (channel-push chan 7) ;; push a non-box value
 (channel-pop chan))`,
		Expect: "7",
	}).Test(t)
}

func TestBoxUnwatchAllFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let* ((chan (make-channel 3))
       (box (make-flow-box :parse "{x:3}")))
 (flow-box-watch box "test" chan)
 (flow-box-unwatch box)
 (flow-box-notify box)
 (channel-push chan 7) ;; push a non-box value
 (channel-pop chan))`,
		Expect: "7",
	}).Test(t)
}

func TestBoxUnwatchSend(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let* ((chan (make-channel 3))
       (box (make-flow-box :parse "{x:3}")))
 (send box :watch "test" chan)
 (send box :unwatch 'test)
 (send box :notify)
 (channel-push chan 7) ;; push a non-box value
 (channel-pop chan))`,
		Expect: "7",
	}).Test(t)
}

func TestBoxUnwatchAllSend(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let* ((chan (make-channel 3))
       (box (make-flow-box :parse "{x:3}")))
 (send box :watch "test" chan)
 (send box :unwatch)
 (send box :notify)
 (channel-push chan 7) ;; push a non-box value
 (channel-pop chan))`,
		Expect: "7",
	}).Test(t)
}

func TestBoxUnwatchNotBox(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-box-unwatch 7)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestBoxUnwatchArgCount(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-box-unwatch (make-flow-box "[1]") 'test t)`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}

func TestBoxUnwatchBadName(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-box-unwatch (make-flow-box "[1]") t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
