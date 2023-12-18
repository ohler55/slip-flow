// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestBoxGetPath(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-get box "x"))`,
		Expect: "3",
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :get (make-bag-path "x")))`,
		Expect: "3",
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :get "y"))`,
		Expect: "nil",
	}).Test(t)
}

func TestBoxGetAll(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-get box))`,
		Expect: `(("x" . 3))`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :get))`,
		Expect: `(("x" . 3))`,
	}).Test(t)
}

func TestBoxGetAsBag(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-get box "x" t))`,
		Expect: "/#<bag-flavor [0-9a-f]+>/",
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :get (make-bag-path "x") t))`,
		Expect: "/#<bag-flavor [0-9a-f]+>/",
	}).Test(t)
}

func TestBoxGetFrozen(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (send box :freeze)
                  (flow-box-get box "x" t))`,
		Expect: "/#<bag-flavor [0-9a-f]+>/",
	}).Test(t)
}

func TestBoxGetNotBox(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-box-get t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestBoxGetBadPath(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-get box t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestBoxGetBadArgCount(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3}")))
                  (flow-box-get box "x" t t))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}
