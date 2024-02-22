// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestBoxWatchNamedFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let* ((chan (make-channel 3))
       (box (make-flow-box :parse "{x:3}")))
 (flow-box-watch box "test" chan)
 (flow-box-notify box "test" "nothing")
 (flow-box-native (channel-pop chan)))`,
		Expect: `(("x" . 3))`,
	}).Test(t)
}

func TestBoxWatchNamedSend(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let* ((chan (make-channel 3))
       (box (make-flow-box :parse "{x:3}")))
 (send box :watch "test" chan)
 (send box :notify 'test "nothing")
 (send (channel-pop chan) :native))`,
		Expect: `(("x" . 3))`,
	}).Test(t)
}

func TestBoxWatchAllFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let* ((chan (make-channel 3))
       (box (make-flow-box :parse "{x:3}")))
 (flow-box-watch box 'test chan)
 (flow-box-notify box)
 (flow-box-native (channel-pop chan)))`,
		Expect: `(("x" . 3))`,
	}).Test(t)
}

func TestBoxWatchAllSend(t *testing.T) {
	(&sliptest.Function{
		Source: `
(let* ((chan (make-channel 3))
       (box (make-flow-box :parse "{x:3}")))
 (send box :watch 'test chan)
 (send box :notify)
 (send (channel-pop chan) :native))`,
		Expect: `(("x" . 3))`,
	}).Test(t)
}

func TestBoxWatchNotBox(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-box-watch 7 'test nil)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestBoxWatchArgCount(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-box-watch (make-flow-box "[1]") 'test)`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}

func TestBoxWatchBadName(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-box-watch (make-flow-box "[1]") t (make-channel 1))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestBoxWatchNotChannel(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-box-watch (make-flow-box "[1]") 'test t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
