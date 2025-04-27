// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ohler55/ojg/tt"
	"github.com/ohler55/slip"
)

func TestGroupDocs(t *testing.T) {
	scope := slip.NewScope()
	var out strings.Builder
	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})

	for _, method := range []string{
		":init",
		":add",
		":remove",
		":find",
		":flows",
		":start",
		":shutdown",
		":set-level",
	} {
		_ = slip.ReadString(fmt.Sprintf(`(describe-method flow-group %s out)`, method), scope).Eval(scope, nil)
		tt.Equal(t, true, strings.Contains(out.String(), method))
		out.Reset()
	}
}
