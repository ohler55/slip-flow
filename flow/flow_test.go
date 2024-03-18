// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ohler55/ojg/tt"
	"github.com/ohler55/slip"
)

func TestFlowDocs(t *testing.T) {
	scope := slip.NewScope()
	var out strings.Builder
	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})

	for _, method := range []string{
		":add-task",
		":entry",
		":find-task",
		":height",
		":init",
		":link",
		":metrics",
		":name",
		":remove-task",
		":reset-metrics",
		":running",
		":set-entry",
		":set-level",
		":shutdown",
		":start",
		":submit",
		":tasks",
		":validate",
		":width",
		":write",
	} {
		_ = slip.ReadString(fmt.Sprintf(`(describe-method flow %s out)`, method)).Eval(scope, nil)
		tt.Equal(t, true, strings.Contains(out.String(), method))
		out.Reset()
	}
}

// func TestFlowInitDoc(t *testing.T) {
// 	scope := slip.NewScope()
// 	var out strings.Builder
// 	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})
// 	_ = slip.ReadString(`(describe-method flow :init out)`).Eval(scope, nil)
// 	fmt.Printf("*** docs: \n%s\n", out.String())
// 	_ = slip.ReadString(`(describe-flavor flow)`).Eval(scope, nil)
// }
