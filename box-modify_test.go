// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestBoxModifyAsLisp(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:[1 2 3]}")))
                  (flow-box-modify box 'reverse "x")
                  (send box :native))`,
		Expect: `(("x" . (3 2 1)))`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:[1 2 3]}")))
	                  (send box :modify 'reverse (make-bag-path "x"))
	                  (send box :native))`,
		Expect: `(("x" . (3 2 1)))`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:[1 2 3]}")))
	                  (send box :modify (lambda (x) (make-instance 'bag-flavor :parse "[2 4 6]")) (make-bag-path "x"))
	                  (send box :native))`,
		Expect: `(("x" . (2 4 6)))`,
	}).Test(t)
}

func TestBoxModifyAsBag(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:[1 2 3]}")))
                  (flow-box-modify box (lambda (b) (bag-set b 4 "[1]")) "x" :as-bag t)
                  (send box :native))`,
		Expect: `(("x" . (1 4 3)))`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:[1 2 3]}")))
                  (flow-box-modify box (lambda (b) (bag-set b 4 "x[1]")) nil :as-bag t)
                  (send box :native))`,
		Expect: `(("x" . (1 4 3)))`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:[1 2 3]}")))
                  (flow-box-modify box (lambda (b) 7) "x" :as-bag t)
                  (send box :native))`,
		Expect: `(("x" . 7))`,
	}).Test(t)
}

func TestBoxModifyFrozen(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:[1 2 3]}")))
                  (send box :freeze)
                  (flow-box-modify box 'reverse "x")
                  (send box :native))`,
		Expect: `(("x" . (3 2 1)))`,
	}).Test(t)
}

func TestBoxModifyNotBox(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-box-modify t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestBoxModifyBadPath(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:[1 2 3]}")))
                  (flow-box-modify box 'reverse t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestBoxModifyBadKeyword(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:[1 2 3]}")))
                  (flow-box-modify box 'reverse "x" :bad t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:[1 2 3]}")))
                  (flow-box-modify box 'reverse "x" :as-bag))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:[1 2 3]}")))
                  (flow-box-modify box 'reverse "x" t t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
