// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ohler55/ojg/tt"
	"github.com/ohler55/slip"
	flow "github.com/ohler55/slip-flow"
)

func TestMakeTask(t *testing.T) {
	task, _ := flow.MakeTask(slip.Fixnum(123))
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
		":metrics",
	} {
		_ = slip.ReadString(fmt.Sprintf(`(describe-method flow-task-flavor %s out)`, method)).Eval(scope, nil)
		tt.Equal(t, true, strings.Contains(out.String(), method))
		out.Reset()
	}
}

// func TestTaskInitDoc(t *testing.T) {
// 	scope := slip.NewScope()
// 	var out strings.Builder
// 	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})
// 	_ = slip.ReadString(`(describe-method flow-task-flavor :init out)`).Eval(scope, nil)
// 	fmt.Printf("*** docs: \n%s\n", out.String())
// }
