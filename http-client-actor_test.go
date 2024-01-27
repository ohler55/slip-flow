// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ohler55/ojg/tt"
	"github.com/ohler55/slip"
)

func TestHTTPClientActorDocs(t *testing.T) {
	scope := slip.NewScope()
	var out strings.Builder
	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})

	for _, method := range []string{
		":init",
		":start",
		":perform",
	} {
		_ = slip.ReadString(fmt.Sprintf(`(describe-method flow-http-client-actor %s out)`, method)).Eval(scope, nil)
		fmt.Printf("*** %s\n", out.String())
		tt.Equal(t, true, strings.Contains(out.String(), method))
		out.Reset()
	}
}
