// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/cl"
	"github.com/ohler55/slip/pkg/flavors"
	"github.com/ohler55/slip/pkg/net"
)

var (
	httpClientActorFlavor *flavors.Flavor
)

func init() {
	Pkg.Initialize(nil)
	httpClientActorFlavor = flavors.DefFlavor("flow-http-client-actor",
		map[string]slip.Object{},
		nil,
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`An HTTP client actor that can be used to make HTTP requests and
then add the results to a box for transition to the next task along a link with the same name
as the HTTP status of the response.
`),
			},
			slip.List{
				slip.Symbol(":init-keywords"),
				slip.Symbol(":method"),
				slip.Symbol(":url"),
				slip.Symbol(":header"),
				slip.Symbol(":trailer"),
				slip.Symbol(":body"),
				slip.Symbol(":timeout"),
				slip.Symbol(":reply-handler"),
			},
		},
		&Pkg,
	)
	httpClientActorFlavor.DefMethod(":init", "", httpClientInitCaller{})
	httpClientActorFlavor.DefMethod(":start", "", httpClientActorStartCaller{})
	httpClientActorFlavor.DefMethod(":perform", "", httpClientActorPerformCaller{})
	httpClientActorFlavor.DefMethod(":init-key-values", "", httpClientActorInitKeyValuesCaller{})
}

type httpClientCtx struct {
	task    *task
	method  strCaller
	url     strCaller
	header  headerCaller
	trailer headerCaller
	body    streamCaller
	timeout durCaller
	handler slip.Caller
}

type httpClientInitCaller struct{}

func (caller httpClientInitCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	self := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		args = args[0].(slip.List)
	}
	hcc := httpClientCtx{}
	for pos := 0; pos < len(args)-1; pos += 2 {
		sym, _ := args[pos].(slip.Symbol)
		switch string(sym) {
		case ":method":
			hcc.method.extract(s, args[pos+1])
		case ":url":
			hcc.url.extract(s, args[pos+1])
		case ":header":
			hcc.header.extract(s, args[pos+1])
		case ":trailer":
			hcc.trailer.extract(s, args[pos+1])
		case ":body":
			hcc.body.extract(s, args[pos+1])
		case ":timeout":
			hcc.timeout.extract(s, args[pos+1])
		case ":reply-handler":
			hcc.handler = cl.ResolveToCaller(s, args[pos+1], depth+1)
		}
	}
	self.Any = &hcc

	return nil
}

func (caller httpClientInitCaller) Docs() string {
	return `__:init__ &key _method_ _url_ _header_ _trailer_ _body_ _timeout_ _reply-handler_
   _:method_ [string|symbol|function] of the request.
   _:url_ [string|function] for the query including the host, port, and path.
   _:header_ [assoc|function] headers for the request.
   _:trailer_ [assoc|function] trailers for the request.
   _:body_ [string|output-stream] for of the request for PUT and POST requests as well as other that have content.
   _:timeout_ [fixnum|function] seconds before timing out waiting for a reply from the HTTP request.
   _:reply-handler_ [function] to call with the response from a request and the box received. If _nil_ then
place the content in a "response" element of the box.


Each argument can be a function that takes a single argument that is the box
of the _:perform_. Is should return the expected type for the argument.
`
}

type httpClientActorStartCaller struct{}

func (caller httpClientActorStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	hcc := obj.Any.(*httpClientCtx)
	hcc.task = args[0].(*flavors.Instance).Any.(*task)

	return nil
}

func (caller httpClientActorStartCaller) Docs() string {
	return `__:start__ _task_
   _:task_ [instance] the task that contains the actor.


Sets the context for the actor.
`
}

type httpClientActorPerformCaller struct{}

func (caller httpClientActorPerformCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	bi := args[0].(*flavors.Instance)

	hcc := obj.Any.(*httpClientCtx)
	method := hcc.method.value(s, bi)
	url := hcc.url.value(s, bi)
	body := hcc.body.value(s, bi)
	timeout := hcc.timeout.value(s, bi)
	ctx := context.Background()
	if 0 < timeout {
		var cf context.CancelFunc
		ctx, cf = context.WithTimeout(ctx, timeout)
		defer cf()
	}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		panic(err)
	}
	req.Header = hcc.header.value(s, bi)
	req.Trailer = hcc.trailer.value(s, bi)

	var (
		client http.Client
		resp   *http.Response
	)
	if resp, err = client.Do(req); err != nil {
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()
	if hcc.handler != nil {
		_ = hcc.handler.Call(s, slip.List{net.MakeResponse(resp), bi}, 0)
	} else {
		bx := bi.Any.(*box)
		if bx.frozen {
			bx.content = alt.Dup(bx.content)
			bx.frozen = false
		}
		content := simplifyHTTPResponse(resp)
		jp.C("response").MustSet(bx.content, content)
	}
	return slip.List{slip.String(strconv.Itoa(resp.StatusCode)), bi}
}

func (caller httpClientActorPerformCaller) Docs() string {
	return `__:perform__ _box_
   _:box_ [instance] box to extract request parameter from.


Makes an HTTP request and passes the response to the _reply-handler_ or if no
_reply-handler_ the response is set as the "reponse" element of the
box. Transition is either on a link matching the response status. If there is
no match then the error link is followed.
`
}

type httpClientActorInitKeyValuesCaller struct{}

func (caller httpClientActorInitKeyValuesCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	hcc := obj.Any.(*httpClientCtx)

	var kvs slip.List
	kvs = append(kvs, slip.List{slip.Symbol(":method"), slip.Tail{Value: hcc.method.raw()}})
	kvs = append(kvs, slip.List{slip.Symbol(":url"), slip.Tail{Value: hcc.url.raw()}})
	kvs = append(kvs, slip.List{slip.Symbol(":timeout"), slip.Tail{Value: hcc.timeout.raw()}})
	switch th := hcc.header.raw().(type) {
	case slip.List:
		kvs = append(kvs, append(slip.List{slip.Symbol(":header")}, th...))
	case *slip.Lambda:
		kvs = append(kvs, slip.List{slip.Symbol(":header"), slip.Tail{Value: th}})
	}
	switch th := hcc.trailer.raw().(type) {
	case slip.List:
		kvs = append(kvs, append(slip.List{slip.Symbol(":trailer")}, th...))
	case *slip.Lambda:
		kvs = append(kvs, slip.List{slip.Symbol(":trailer"), slip.Tail{Value: th}})
	}
	if body := hcc.body.raw(); body != slip.String("") {
		kvs = append(kvs, slip.List{slip.Symbol(":body"), slip.Tail{Value: hcc.body.raw()}})
	}
	if lam, ok := hcc.handler.(*slip.Lambda); ok {
		kvs = append(kvs, slip.List{slip.Symbol(":reply-handler"), slip.Tail{Value: lam}})
	}
	return kvs
}

func (caller httpClientActorInitKeyValuesCaller) Docs() string {
	return `__:init-key-values__ => ((:method get) (:timeout 1))


Returns the keywords and values needed to recreate the instance.
`
}

func simplifyHTTPResponse(resp *http.Response) any {
	body, _ := io.ReadAll(resp.Body)
	return map[string]any{
		"status":        int64(resp.StatusCode),
		"proto":         resp.Proto,
		"header":        simplifyHTTPHeader(resp.Header),
		"contentLength": resp.ContentLength,
		"trailer":       simplifyHTTPHeader(resp.Trailer),
		"body":          string(body),
	}
}

func simplifyHTTPHeader(h http.Header) (sh any) {
	if h != nil {
		header := map[string]any{}
		for k, va := range h {
			vlist := make([]any, len(va))
			for i, s := range va {
				vlist[i] = s
			}
			header[k] = vlist
		}
		sh = header
	}
	return
}
