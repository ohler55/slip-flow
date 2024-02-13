// Copyright (c) 2024, Peter Ohler, All rights reserved.

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
	task, _ := flow.MakeTask(slip.Symbol(":name"), slip.String("tisk"))
	tt.Equal(t, "/#<flow-task [0-9a-f]+>/", task.String())

	task, _ = flow.MakeTask(slip.Symbol(":name"), slip.Symbol("tisk"))
	tt.Equal(t, "/#<flow-task [0-9a-f]+>/", task.String())
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
		":reset-metrics",
		":links",
		":unlink",
		":transition",
		":update-link",
	} {
		_ = slip.ReadString(fmt.Sprintf(`(describe-method flow-task %s out)`, method)).Eval(scope, nil)
		tt.Equal(t, true, strings.Contains(out.String(), method))
		out.Reset()
	}
}

func TestTaskInitLogger(t *testing.T) {
	(&sliptest.Function{
		Source: `(make-instance 'flow-task :name "tisk" :logger (make-instance 'logger-flavor))`,
		Expect: "/#<flow-task [0-9a-f]+>/",
	}).Test(t)
}

func TestTaskInitBadName(t *testing.T) {
	(&sliptest.Function{
		Source:    `(make-instance 'flow-task :name 123)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskInitBadWorkers(t *testing.T) {
	(&sliptest.Function{
		Source:    `(make-instance 'flow-task :workers t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskInitBadDepth(t *testing.T) {
	(&sliptest.Function{
		Source:    `(make-instance 'flow-task :depth t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskInitBadActor(t *testing.T) {
	(&sliptest.Function{
		Source:    `(make-instance 'flow-task :actor t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestTaskInitBadKeyword(t *testing.T) {
	(&sliptest.Function{
		Source:    `(make-instance 'flow-task :bad t)`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}

// func TestTaskInitDoc(t *testing.T) {
// 	scope := slip.NewScope()
// 	var out strings.Builder
// 	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})
// 	// _ = slip.ReadString(`(describe-method flow-task :init out)`).Eval(scope, nil)
// 	// fmt.Printf("*** docs: \n%s\n", out.String())
// 	_ = slip.ReadString(`(describe-flavor flow-task)`).Eval(scope, nil)
// }
