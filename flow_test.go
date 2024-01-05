// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

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
		":init",
		":name",
		":start",
		":shutdown",
		":running",
		":add-task",
		// ":link",
		// ":unlink",
		// ":remove-task",
		// ":tasks",
		// ":find-task",
		// ":entry",
		// ":set-entry",
		// ":submit",
		":exit-channel",
	} {
		_ = slip.ReadString(fmt.Sprintf(`(describe-method flow-flavor %s out)`, method)).Eval(scope, nil)
		tt.Equal(t, true, strings.Contains(out.String(), method))
		out.Reset()
	}
}

// func TestFlowInitDoc(t *testing.T) {
// 	scope := slip.NewScope()
// 	var out strings.Builder
// 	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})
// 	_ = slip.ReadString(`(describe-method flow-flavor :init out)`).Eval(scope, nil)
// 	fmt.Printf("*** docs: \n%s\n", out.String())
// 	_ = slip.ReadString(`(describe-flavor flow-flavor)`).Eval(scope, nil)
// }
