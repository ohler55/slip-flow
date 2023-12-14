// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ohler55/ojg/tt"
	"github.com/ohler55/slip"
)

func TestBoxDocs(t *testing.T) {
	scope := slip.NewScope()
	var out strings.Builder
	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})

	for _, method := range []string{
		":init",
		":set",
		":parse",
		":read",
		":get",
		":has",
		":remove",
		":modify",
		":native",
		":write",
		":walk",
		":bag",
		":freeze",
		":thaw",
		":frozen",
		":tracking-id",
		":track",
		":history",
		":scan",
		":copy",
	} {
		_ = slip.ReadString(fmt.Sprintf(`(describe-method flow-box-flavor %s out)`, method)).Eval(scope, nil)
		tt.Equal(t, true, strings.Contains(out.String(), method))
		out.Reset()
	}
}
