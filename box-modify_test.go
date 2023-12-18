// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

/*
func TestBoxModifyOk(t *testing.T) {
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:[1 2 3]}")))
                  (flow-box-modify box 'reverse "x")
                  (send box :native))`,
		Expect: `(("x" . (3 2 1)))`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(let ((box (make-flow-box :parse "{x:3 y:4}")))
                  (send box :modify 'reverse (make-bag-path "x"))
                  (send box :native))`,
		Expect: `(("x" . (3 2 1)))`,
	}).Test(t)
}
*/
