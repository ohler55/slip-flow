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

func TestBoxWalkLambdaList(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3 y: 4}"))
                      (values nil))
                  (flow-box-walk box '(lambda (x) (setq values (add values x))) "x")
                  values)`,
		Expect: `(3)`,
	}).Test(t)
}

func TestBoxWalkAsBag(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3 y: 4}"))
                      (values nil))
                  (send box :walk (lambda (x) (setq values (add values (send x :native)))) "x" :as-lisp nil)
                  values)`,
		Expect: `(3)`,
	}).Test(t)
}

func TestBoxWalkDefun(t *testing.T) {
	(&sliptest.Function{
		Source: `(defvar box-walk-func-values nil)
                 (defun  box-walk-func-test (x) (setq box-walk-func-values (add box-walk-func-values x)))
                 (let ((box (make-flow-box :parse "{x:3 y: 4}")))
                  (send box :walk 'box-walk-func-test "x")
                  box-walk-func-values)`,
		Expect: `(3)`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(defun  box-walk-func-test-bag (x)
                  (setq box-walk-func-values (add box-walk-func-values (send x :native))))
                 (let ((box (make-flow-box :parse "{x:3 y: 4}")))
                  (setq box-walk-func-values nil)
                  (send box :walk 'box-walk-func-test-bag "x" :as-lisp nil)
                  box-walk-func-values)`,
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

func TestBoxWalkNotFunction(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3 y: 4}")))
                  (flow-box-walk box t "x"))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestBoxWalkNotBox(t *testing.T) {
	(&sliptest.Function{
		Source:    `(flow-box-walk t (lambda (x) nil))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
