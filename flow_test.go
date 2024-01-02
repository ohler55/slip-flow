// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ohler55/slip"
)

func TestFlowInitDoc(t *testing.T) {
	scope := slip.NewScope()
	var out strings.Builder
	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})
	_ = slip.ReadString(`(describe-method flow-flavor :init out)`).Eval(scope, nil)
	fmt.Printf("*** docs: \n%s\n", out.String())
	_ = slip.ReadString(`(describe-flavor flow-flavor)`).Eval(scope, nil)
}
