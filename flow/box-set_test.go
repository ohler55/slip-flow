// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestBoxSetPath(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-set box 7 "x")
                  (send box :write nil))`,
		Expect: `"{x: 7}"`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :set 7 (make-bag-path "x"))
                  (send box :write nil))`,
		Expect: `"{x: 7}"`,
	}).Test(t)
}

func TestBoxSetAll(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-set box 7)
                  (send box :write nil))`,
		Expect: `"7"`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :set 7)
                  (send box :write nil))`,
		Expect: `"7"`,
	}).Test(t)
}

func TestBoxSetFrozen(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :freeze)
                  (flow-box-set box 7 "x")
                  (list (send box :write nil) (send box :frozen)))`,
		Expect: `("{x: 7}" nil)`,
	}).Test(t)
}

func TestBoxSetNotBox(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-box-set t 7 "x")`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestBoxSetBadPath(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-set box 7 t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestBoxSetBadArgCount(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :set))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}
