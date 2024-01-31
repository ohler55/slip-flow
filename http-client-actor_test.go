// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ohler55/ojg/tt"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/gi"
	"github.com/ohler55/slip/sliptest"
)

func TestHTTPClientActorBasic(t *testing.T) {
	exitChan := make(gi.Channel, 5)
	scope := slip.NewScope()
	scope.Let("exit-channel", exitChan)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		_, _ = w.Write([]byte("Hello\n"))
	}))
	defer server.Close()

	scope.Let("test-url", slip.String(server.URL))

	(&sliptest.Function{
		Scope: scope,
		Source: `
(let ((flow (make-flow :name 'flo :exit-channel exit-channel)))
  (flow-add-task flow
                 :name "start"
                 :actor (make-instance 'flow-http-client-actor
                                       :method "get"
                                       :timeout 1
                                       :url test-url))
  (flow-add-task flow
                 :name "done"
                 :actor (make-instance 'flow-exit-actor))

  (flow-link flow "200" 'start "done")
  (flow-set-entry flow 'start)
  (send flow :set-level 'warn)
  (flow-submit flow (make-flow-box :parse "{a:1}"))
  (flow-shutdown flow)
)`,
		Expect: "nil",
	}).Test(t)

	out := <-exitChan
	scope.Let("http-client-out", out)

	value := slip.ReadString(`(send http-client-out :get "response.body")`).Eval(scope, nil)
	tt.Equal(t, `"Hello
"`, slip.ObjectString(value))

	value = slip.ReadString(`(send http-client-out :get "response.status")`).Eval(scope, nil)
	tt.Equal(t, "200", slip.ObjectString(value))

	value = slip.ReadString(`(send http-client-out :get "response.contentLength")`).Eval(scope, nil)
	tt.Equal(t, "6", slip.ObjectString(value))

	value = slip.ReadString(`(send http-client-out :get "['response']['header']['Content-Type'][0]")`).Eval(scope, nil)
	tt.Equal(t, `"text/plain; charset=utf-8"`, slip.ObjectString(value))
}

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
		// fmt.Printf("*** %s\n", out.String())
		tt.Equal(t, true, strings.Contains(out.String(), method))
		out.Reset()
	}
}
