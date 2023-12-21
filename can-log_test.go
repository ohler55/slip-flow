// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/ohler55/ojg/tt"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/sliptest"
)

func TestCanLogErrorLevel(t *testing.T) {
	scope := slip.NewScope()
	var log bytes.Buffer
	scope.Let(slip.Symbol("log-out"), &slip.OutputStream{Writer: &log})
	_ = slip.ReadString("(setq logger (make-instance 'logger-flavor))").Eval(scope, nil)
	_ = slip.ReadString("(send logger :set-out log-out)").Eval(scope, nil)
	_ = slip.ReadString(`(unless (find-flavor 'can-log-tester)
                          (defflavor can-log-tester () (can-log-flavor)))`).Eval(scope, nil)

	_ = slip.ReadString(`(let ((can (make-instance 'can-log-tester :log-level 0 :logger logger)))
                  (send can :error "error ~D" 0)
                  (send can :warn "warn ~D" 1)
                  (send can :info "info ~D" 2)
                  (send can :debug "debug ~D" 3)

                  (send can :set-level -1)
                  (send can :error "----- -1")
                  (send can :error "error ~D" 0)
                  (send can :warn "warn ~D" 1)
                  (send can :info "info ~D" 2)
                  (send can :debug "debug ~D" 3)

                  (send can :set-level 4)
                  (send can :error "----- 4")
                  (send can :error "error ~D" 0)
                  (send can :warn "warn ~D" 1)
                  (send can :info "info ~D" 2)
                  (send can :debug "debug ~D" 3)

                  (send can :set-level :error)
                  (send can :error "----- :error")
                  (send can :error "error ~D" 0)
                  (send can :warn "warn ~D" 1)
                  (send can :info "info ~D" 2)
                  (send can :debug "debug ~D" 3)

                  (send can :set-level :warn)
                  (send can :error "----- :warn")
                  (send can :error "error ~D" 0)
                  (send can :warn "warn ~D" 1)
                  (send can :info "info ~D" 2)
                  (send can :debug "debug ~D" 3)

                  (send can :set-level :info)
                  (send can :error "----- :info")
                  (send can :error "error ~D" 0)
                  (send can :warn "warn ~D" 1)
                  (send can :info "info ~D" 2)
                  (send can :debug "debug ~D" 3)

                  (send can :set-level :debug)
                  (send can :error "----- :debug")
                  (send can :error "error ~D" 0)
                  (send can :warn "warn ~D" 1)
                  (send can :info "info ~D" 2)
                  (send can :debug "debug ~D" 3)

                  (send logger :shutdown))`).Eval(scope, nil)
	tt.Equal(t, `E error 0
E ----- -1
E error 0
E ----- 4
E error 0
W warn 1
I info 2
D debug 3
E ----- :error
E error 0
E ----- :warn
E error 0
W warn 1
I info 2
E ----- :info
E error 0
W warn 1
I info 2
E ----- :debug
E error 0
W warn 1
I info 2
D debug 3
`, log.String())

	(&sliptest.Function{
		Scope:     scope,
		Source:    `(let ((can (make-instance 'can-log-tester :logger logger))) (send can :set-level t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestCanLogDocs(t *testing.T) {
	scope := slip.NewScope()
	var out strings.Builder
	scope.Let(slip.Symbol("out"), &slip.OutputStream{Writer: &out})

	for _, method := range []string{
		":error",
		":warn",
		":info",
		":debug",
		":set-level",
	} {
		_ = slip.ReadString(fmt.Sprintf(`(describe-method can-log-flavor %s out)`, method)).Eval(scope, nil)
		tt.Equal(t, true, strings.Contains(out.String(), method))
		out.Reset()
	}
}
