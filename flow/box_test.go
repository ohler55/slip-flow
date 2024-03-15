// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ohler55/ojg/tt"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip-flow/flow"
)

func TestMakeBox(t *testing.T) {
	box, _ := flow.MakeBox(slip.Fixnum(123))
	tt.Equal(t, "/#<flow-box [0-9a-f]+>/", box.String())
}

func TestBoxDocs(t *testing.T) {
	scope := slip.NewScope()
	var out strings.Builder
	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})

	for _, method := range []string{
		":bag",
		":copy",
		":freeze",
		":frozen",
		":get",
		":has",
		":history",
		":init",
		":merge",
		":modify",
		":native",
		":notify",
		":parse",
		":read",
		":remove",
		":scan",
		":set",
		":thaw",
		":track",
		":tracking-id",
		":unwatch",
		":walk",
		":watch",
		":write",
	} {
		_ = slip.ReadString(fmt.Sprintf(`(describe-method flow-box %s out)`, method)).Eval(scope, nil)
		tt.Equal(t, true, strings.Contains(out.String(), method))
		out.Reset()
	}
}

// func TestBoxInitDoc(t *testing.T) {
// 	scope := slip.NewScope()
// 	var out strings.Builder
// 	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})
// 	_ = slip.ReadString(`(describe-method flow-box :write out)`).Eval(scope, nil)
// 	fmt.Printf("***\n%s\n", out.String())
// 	_ = slip.ReadString(`(describe-flavor flow-box)`).Eval(scope, nil)
// }
