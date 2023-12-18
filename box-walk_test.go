// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

import (
	"testing"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestBoxWalkMethod(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3 y: 4}"))
                      (values nil))
                  (send box :walk (lambda (x) (setq values (add values x))) "x")
                  values)`,
		Expect: `(3)`,
	}).Test(t)
}

func TestBoxWalkFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3 y: 4}"))
                      (values nil))
                  (flow-box-walk box (lambda (x) (setq values (add values x))) (make-bag-path "x"))
                  values)`,
		Expect: `(3)`,
	}).Test(t)
}

func TestBoxWalkBadPath(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3 y: 4}")))
                  (flow-box-walk box (lambda (x) nil) t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestBoxWalkNotBox(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-box-walk t (lambda (x) nil))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
