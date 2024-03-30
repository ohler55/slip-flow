// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow_test

import (
	"testing"

	"github.com/ohler55/slip/sliptest"
)

func TestWriteFileActorBadDestination(t *testing.T) {
	(&sliptest.Function{
		Source: `(make-instance 'flow-write-file-actor
                                :filename "testdata/write-me.text"
                                :content "Sample content"
                                :permissions 0660)`,
		Expect: "/#<flow-write-file-actor [0-9a-f]+>/",
	}).Test(t)
}
