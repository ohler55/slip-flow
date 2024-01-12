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

func TestTrackID(t *testing.T) {
	scope := slip.NewScope()
	ti, _ := flow.MakeTrack(slip.Fixnum(123))
	scope.Let("track", ti)
	(&sliptest.Function{
		Scope:  scope,
		Source: `(send track :id)`,
		Expect: "123",
	}).Test(t)
}

func TestTrackScan(t *testing.T) {
	scope := slip.NewScope()
	ti, track := flow.MakeTrack(slip.String("abc"))
	scope.Let("track", ti)
	track.Scan("flo", "task-1")
	track.Scan("flo", "task-2")
	(&sliptest.Function{
		Scope:  scope,
		Source: `(send track :history)`,
		Validate: func(t *testing.T, v slip.Object) {
			history, ok := v.(slip.List)
			tt.Equal(t, true, ok)
			tt.Equal(t, 2, len(history))
			for i, e := range history {
				ev, ok2 := e.(slip.List)
				tt.Equal(t, true, ok2)
				tt.Equal(t, 3, len(ev))
				tt.SameType(t, slip.Time{}, ev[0])
				tt.Equal(t, slip.String(fmt.Sprintf("task-%d", i+1)), ev[1])
				tt.Equal(t, slip.String("flo"), ev[2])
			}
		},
	}).Test(t)
}

func TestTrackDocs(t *testing.T) {
	scope := slip.NewScope()
	var out strings.Builder
	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})

	for _, method := range []string{
		":id",
		":history",
	} {
		_ = slip.ReadString(fmt.Sprintf(`(describe-method flow-track-flavor %s out)`, method)).Eval(scope, nil)
		tt.Equal(t, true, strings.Contains(out.String(), method))
		out.Reset()
	}
}
