// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ohler55/ojg/tt"
	"github.com/ohler55/slip"
	flow "github.com/ohler55/slip-flow"
	"github.com/ohler55/slip/sliptest"
)

func TestMakeTask(t *testing.T) {
	task, _ := flow.MakeTask(slip.String("tisk"))
	tt.Equal(t, "/#<flow-task-flavor [0-9a-f]+>/", task.String())

	task, _ = flow.MakeTask(slip.Symbol("tisk"))
	tt.Equal(t, "/#<flow-task-flavor [0-9a-f]+>/", task.String())
}

func TestTaskDocs(t *testing.T) {
	scope := slip.NewScope()
	var out strings.Builder
	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})

	for _, method := range []string{
		":init",
		":name",
		":workers",
		":start",
		":shutdown",
		":running",
		":receive",
		":metrics",
	} {
		_ = slip.ReadString(fmt.Sprintf(`(describe-method flow-task-flavor %s out)`, method)).Eval(scope, nil)
		tt.Equal(t, true, strings.Contains(out.String(), method))
		out.Reset()
	}
}

func TestTaskInitLogger(t *testing.T) {
	(&sliptest.Function{
		Source: `(make-instance 'flow-task-flavor :name "tisk" :logger (make-instance 'logger-flavor))`,
		Expect: "/#<flow-task-flavor [0-9a-f]+>/",
	}).Test(t)
}

func TestTaskInitBadName(t *testing.T) {
	(&sliptest.Function{
		Source:    `(make-instance 'flow-task-flavor :name 123)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskInitBadWorkers(t *testing.T) {
	(&sliptest.Function{
		Source:    `(make-instance 'flow-task-flavor :workers t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskInitBadDepth(t *testing.T) {
	(&sliptest.Function{
		Source:    `(make-instance 'flow-task-flavor :depth t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskInitBadActor(t *testing.T) {
	(&sliptest.Function{
		Source:    `(make-instance 'flow-task-flavor :actor t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskInitBadKeyword(t *testing.T) {
	(&sliptest.Function{
		Source:    `(make-instance 'flow-task-flavor :bad t)`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}

// func TestTaskInitDoc(t *testing.T) {
// 	scope := slip.NewScope()
// 	var out strings.Builder
// 	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})
// 	// _ = slip.ReadString(`(describe-method flow-task-flavor :init out)`).Eval(scope, nil)
// 	// fmt.Printf("*** docs: \n%s\n", out.String())
// 	_ = slip.ReadString(`(describe-flavor flow-task-flavor)`).Eval(scope, nil)
// }
