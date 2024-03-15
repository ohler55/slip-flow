// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"testing"
	"time"

	"github.com/ohler55/ojg/tt"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestMakeFlowBoxParse(t *testing.T) {
	scope := slip.NewScope()
	scope.Set("*flow-box-time-format*", slip.String(time.RFC3339Nano))
	tf := sliptest.Function{
		Scope:  scope,
		Source: `(make-flow-box :parse "{x:7}")`,
		Expect: "/#<flow-box [0-9a-f]+>/",
	}
	tf.Test(t)
	scope.Let("box", tf.Result)
	out := slip.ReadString(`(send box :write nil)`).Eval(scope, nil).(slip.String)
	tt.Equal(t, "{x: 7}", string(out))

	(&sliptest.Function{
		Source:    `(make-flow-box :parse t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestMakeFlowBoxSet(t *testing.T) {
	scope := slip.NewScope()
	tf := sliptest.Function{
		Scope:  scope,
		Source: `(make-flow-box :set (make-instance 'bag-flavor :parse "{x:7}"))`,
		Expect: "/#<flow-box [0-9a-f]+>/",
	}
	tf.Test(t)
	scope.Let("box", tf.Result)
	out := slip.ReadString(`(send box :write nil)`).Eval(scope, nil).(slip.String)
	tt.Equal(t, "{x: 7}", string(out))

	(&sliptest.Function{
		Source: `(send (make-flow-box :set '((x . 3))) :write nil)`,
		Expect: `"{x: 3}"`,
	}).Test(t)

	(&sliptest.Function{
		Source:    `(make-flow-box :set (make-instance 'vanilla-flavor))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestMakeFlowBoxRead(t *testing.T) {
	scope := slip.NewScope()
	tf := sliptest.Function{
		Scope:  scope,
		Source: `(make-flow-box :read (make-string-input-stream "{x:7}"))`,
		Expect: "/#<flow-box [0-9a-f]+>/",
	}
	tf.Test(t)
	scope.Let("box", tf.Result)
	out := slip.ReadString(`(send box :write nil)`).Eval(scope, nil).(slip.String)
	tt.Equal(t, "{x: 7}", string(out))

	(&sliptest.Function{
		Source:    `(make-flow-box :read t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestMakeFlowBoxTrackingID(t *testing.T) {
	(&sliptest.Function{
		Source: `(send (make-flow-box :tracking-id "abcd") :tracking-id)`,
		Expect: `"abcd"`,
	}).Test(t)
	(&sliptest.Function{
		Source: `(send (make-flow-box :tracking-id 1234) :tracking-id)`,
		Expect: "1234",
	}).Test(t)
}

func TestMakeFlowBoxTrack(t *testing.T) {
	(&sliptest.Function{
		Source: `(send (make-flow-box :track (send (make-flow-box :tracking-id 1234) :track)) :tracking-id)`,
		Expect: "1234",
	}).Test(t)
	(&sliptest.Function{
		Source:    `(make-flow-box :track t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestMakeFlowBoxBadKeyword(t *testing.T) {
	(&sliptest.Function{
		Source:    `(make-flow-box :bad t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestMakeFlowBoxBadWatch(t *testing.T) {
	(&sliptest.Function{
		Source:    `(make-flow-box :watch t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
	(&sliptest.Function{
		Source:    `(let ((quux 7)) (make-flow-box :watch 'quux))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
