// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"testing"

	"github.com/ohler55/ojg/tt"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestPkgTimeFormat(t *testing.T) {
	scope := slip.NewScope()
	orig := slip.ReadString(`*flow-box-time-format*`).Eval(scope, nil)
	defer scope.Set("*flow-box-time-format*", orig)

	_ = slip.ReadString(`(setq *flow-box-time-format* "2006-01-02")`).Eval(scope, nil)
	v := slip.ReadString(`*flow-box-time-format*`).Eval(scope, nil)
	tt.Equal(t, `"2006-01-02"`, slip.ObjectString(v))
	(&sliptest.Function{
		Source: `(send (make-flow-box :parse "[\"2023-12-13\"]") :native)`,
		Expect: `(@2023-12-13T00:00:00Z)`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(send (make-flow-box :parse "[\"2023-abcdef\"]") :native)`,
		Expect: `("2023-abcdef")`,
	}).Test(t)

	_ = slip.ReadString(`(setq *flow-box-time-format* *rfc3339nano*)`).Eval(scope, nil)
	(&sliptest.Function{
		Source: `(send (make-flow-box :parse "[\"2022-09-19T01:02:03.000Z\"]") :native)`,
		Expect: `(@2022-09-19T01:02:03Z)`,
	}).Test(t)

	_ = slip.ReadString(`(setq *flow-box-time-format* nil)`).Eval(scope, nil)
	v = slip.ReadString(`*flow-box-time-format*`).Eval(scope, nil)
	tt.Equal(t, `nil`, slip.ObjectString(v))

	_ = slip.ReadString(`(setq *flow-box-time-format* 'nano)`).Eval(scope, nil)
	v = slip.ReadString(`*flow-box-time-format*`).Eval(scope, nil)
	tt.Equal(t, `"nano"`, slip.ObjectString(v))

	_ = slip.ReadString(`(setq *flow-box-time-format* 'second)`).Eval(scope, nil)
	v = slip.ReadString(`*flow-box-time-format*`).Eval(scope, nil)
	tt.Equal(t, `"second"`, slip.ObjectString(v))
	(&sliptest.Function{
		Source: `(send (make-flow-box :parse "[1663549323.000000000]") :native)`,
		Expect: `(@2022-09-19T01:02:03Z)`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(send (make-flow-box :parse "[1.1]") :native)`,
		Expect: `(1.1)`,
	}).Test(t)

	(&sliptest.Function{
		Source:    `(setq *flow-box-time-format* t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestPkgTimeWrap(t *testing.T) {
	scope := slip.NewScope()
	orig := slip.ReadString(`*flow-box-time-wrap*`).Eval(scope, nil)
	defer scope.Set("*flow-box-time-wrap*", orig)

	_ = slip.ReadString(`(setq *flow-box-time-wrap* "time")`).Eval(scope, nil)
	v := slip.ReadString(`*flow-box-time-wrap*`).Eval(scope, nil)
	tt.Equal(t, `"time"`, slip.ObjectString(v))
	(&sliptest.Function{
		Source: `(send (make-flow-box :parse "{time:\"2023-12-13\"}") :native)`,
		Expect: `@2023-12-13T00:00:00Z`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(send (make-flow-box :parse "{time:1663549323000000000}") :native)`,
		Expect: `@2022-09-19T01:02:03Z`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(send (make-flow-box :parse "{time:1.25}") :native)`,
		Expect: `(("time" . 1.25))`,
	}).Test(t)

	_ = slip.ReadString(`(setq *flow-box-time-wrap* nil)`).Eval(scope, nil)
	v = slip.ReadString(`*flow-box-time-wrap*`).Eval(scope, nil)
	tt.Equal(t, `nil`, slip.ObjectString(v))

	_ = slip.ReadString(`(setq *flow-box-time-wrap* 'time)`).Eval(scope, nil)
	v = slip.ReadString(`*flow-box-time-wrap*`).Eval(scope, nil)
	tt.Equal(t, `"time"`, slip.ObjectString(v))

	(&sliptest.Function{
		Source:    `(setq *flow-box-time-wrap* t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
