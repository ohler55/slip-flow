// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestBoxMergeFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}"))
                       (other (make-flow-box :parse "{y:4}")))
                  (flow-box-merge box other)
                  (send box :write nil))`,
		Expect: `"{x: 3 y: 4}"`,
	}).Test(t)
}

func TestBoxMergeSend(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}"))
                       (other (make-flow-box :parse "{y:4}")))
                  (send box :freeze)
                  (send box :merge other)
                  (send box :write nil))`,
		Expect: `"{x: 3 y: 4}"`,
	}).Test(t)
}

func TestBoxMergeNotBox(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-merge box t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
	(&sliptest.Function{
		Source:    `(flow-box-merge t t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
