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
	testHTTPClientActorOk(t, `(make-instance 'flow-http-client-actor
                                             :method 'get
                                             :timeout 1
                                             :header '((Accept . "text/html"))
                                             :trailer '()
                                             :url test-url)`, true)
}

func TestHTTPClientActorAlt(t *testing.T) {
	testHTTPClientActorOk(t, `(make-instance 'flow-http-client-actor
                                             :method 'get
                                             :timeout '1s
                                             :url test-url)`, true)
}

func TestHTTPClientActorAlt2(t *testing.T) {
	testHTTPClientActorOk(t, `(make-instance 'flow-http-client-actor
                                             :method "get"
                                             :timeout "1s"
                                             :url test-url)`, true)
}

func TestHTTPClientActorAlt3(t *testing.T) {
	testHTTPClientActorOk(t, `(make-instance 'flow-http-client-actor
                                             :method "get"
                                             :timeout (lambda (b) "1s")
                                             :url test-url)`, true)
}

func TestHTTPClientActorAlt4(t *testing.T) {
	testHTTPClientActorOk(t, `(make-instance 'flow-http-client-actor
                                             :method "get"
                                             :timeout (lambda (b) '1s)
                                             :url test-url)`, true)
}

func TestHTTPClientActorAlt5(t *testing.T) {
	testHTTPClientActorOk(t, `(make-instance 'flow-http-client-actor
                                             :method "get"
                                             :timeout (lambda (b) nil)
                                             :url test-url)`, true)
}

func TestHTTPClientActorAlt6(t *testing.T) {
	testHTTPClientActorOk(t, `(make-instance 'flow-http-client-actor
                                             :method "get"
                                             :timeout nil
                                             :url test-url)`, true)
}

func TestHTTPClientActorFuncs(t *testing.T) {
	testHTTPClientActorOk(t, `(make-instance 'flow-http-client-actor
                                             :method (lambda (b) 'get)
                                             :timeout (lambda (b) (send b :get "a"))
                                             :header (lambda (b) '((Accept . "text/html")))
                                             :url (lambda (b) test-url)
                                             :reply-handler (lambda (r b)
                                                             (send b :set (send r :content) "response.body")))`, false)
}

func TestHTTPClientActorPostFuncStream(t *testing.T) {
	testHTTPClientActorOk(t, `(make-instance 'flow-http-client-actor
                                             :method "post"
                                             :timeout 1
                                             :body (lambda (b) (make-string-input-stream "Hello\n"))
                                             :url test-url)`, true)
}

func TestHTTPClientActorPostFuncString(t *testing.T) {
	testHTTPClientActorOk(t, `(make-instance 'flow-http-client-actor
                                             :method "post"
                                             :timeout 1
                                             :body (lambda (b) "Hello\n")
                                             :url test-url)`, true)
}

func TestHTTPClientActorPostString(t *testing.T) {
	testHTTPClientActorOk(t, `(make-instance 'flow-http-client-actor
                                             :method "post"
                                             :timeout 1
                                             :body "Hello\n"
                                             :url test-url)`, true)
}

func testHTTPClientActorOk(t *testing.T, actor string, checkAll bool) {
	exitChan := make(gi.Channel, 5)
	scope := slip.NewScope()
	scope.Let("exit-channel", exitChan)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		_, _ = w.Write([]byte("Hello\n"))
	}))
	defer server.Close()

	scope.Let("test-url", slip.String(server.URL))
	scope.Let("actor", slip.ReadString(actor).Eval(scope, nil))

	(&sliptest.Function{
		Scope: scope,
		Source: `
(let ((flow (make-flow :name 'flo :exit-channel exit-channel)))
  (flow-add-task flow
                 :name "start"
                 :actor actor)
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

	if !checkAll {
		return
	}
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

func TestHTTPClientActorBadMethod(t *testing.T) {
	(&sliptest.Function{
		Source:    `(make-instance 'flow-http-client-actor :method t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
	(&sliptest.Function{
		Source:    `(send (make-instance 'flow-http-client-actor :method (lambda (b) t)) :perform (make-flow-box))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
	(&sliptest.Function{
		Source:    `(send (make-instance 'flow-http-client-actor :method 'quux) :perform (make-flow-box))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
	(&sliptest.Function{
		Source:    `(send (make-instance 'flow-http-client-actor :method " ") :perform (make-flow-box))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}

func TestHTTPClientActorBadTimeout(t *testing.T) {
	(&sliptest.Function{
		Source:    `(make-instance 'flow-http-client-actor :timeout t)`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
	(&sliptest.Function{
		Source:    `(make-instance 'flow-http-client-actor :timeout "xyz")`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
	(&sliptest.Function{
		Source:    `(make-instance 'flow-http-client-actor :timeout 'xyz)`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
	(&sliptest.Function{
		Source:    `(send (make-instance 'flow-http-client-actor :timeout (lambda (b) t)) :perform (make-flow-box))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
	(&sliptest.Function{
		Source:    `(send (make-instance 'flow-http-client-actor :timeout (lambda (b) "xyz")) :perform (make-flow-box))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
	(&sliptest.Function{
		Source:    `(send (make-instance 'flow-http-client-actor :timeout (lambda (b) 'xyz)) :perform (make-flow-box))`,
		PanicType: slip.Symbol("error"),
	}).Test(t)
}

func TestHTTPClientActorBadHeader(t *testing.T) {
	(&sliptest.Function{
		Source:    `(make-instance 'flow-http-client-actor :header '(t))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
	(&sliptest.Function{
		Source:    `(send (make-instance 'flow-http-client-actor :header (lambda (b) t)) :perform (make-flow-box))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
	(&sliptest.Function{
		Source:    `(send (make-instance 'flow-http-client-actor :header (lambda (b) '(t))) :perform (make-flow-box))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}

func TestHTTPClientActorBadBody(t *testing.T) {
	(&sliptest.Function{
		Source:    `(send (make-instance 'flow-http-client-actor :body (lambda (b) t)) :perform (make-flow-box))`,
		PanicType: slip.Symbol("type-error"),
	}).Test(t)
}
