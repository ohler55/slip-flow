// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

import (
	"fmt"
	"net/http"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/cl"
	"github.com/ohler55/slip/pkg/flavors"
)

var (
	httpClientActorFlavor *flavors.Flavor
)

func init() {
	httpClientActorFlavor = flavors.DefFlavor("flow-http-client-actor",
		map[string]slip.Object{},
		nil,
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`A flow-http-client-actor TBD
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
	)
	httpClientActorFlavor.DefMethod(":init", "", httpClientInitCaller{})
	httpClientActorFlavor.DefMethod(":start", "", httpClientActorStartCaller{})
	httpClientActorFlavor.DefMethod(":perform", "", httpClientActorPerformCaller{})
}

type httpClientCtx struct {
	task    *task
	method  strCaller
	url     strCaller
	header  slip.Object // TBD
	trailer slip.Object // TBD
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

	// TBD check early or wait until :perform?
	//  will need similar for both cases

	for pos := 0; pos < len(args)-1; pos += 2 {
		sym, _ := args[pos].(slip.Symbol)
		switch string(sym) {
		case ":method":
			hcc.method.extract(s, args[pos+1])
		case ":url":
			hcc.method.extract(s, args[pos+1])
		case ":header":
			hcc.header = args[pos+1]
			// TBD verify either assoc or function
		case ":trailer":
			hcc.trailer = args[pos+1]
			// TBD verify either assoc or function
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
   _method_ [string|symbol|function] of the request.
   _url_ [string|function] for the query including the host, port, and path.
   _header_ [assoc|function] headers for the request.
   _trailer_ [assoc|function] trailers for the request.
   _body_ [string|output-stream] for of the request for PUT and POST requests as well as other that have content.
   _:timeout_ [fixnum|function] seconds before timing out waiting for a reply from the HTTP request.
   _reply-handler_ [function] to call with the response from a request and the box received. If _nil_ then
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
	url := hcc.method.value(s, bi)

	// TBD hcc.getString(hcc.method, hcc.methodCaller)
	// TBD maybe struct for string-caller, same for headers, timeout, body

	// TBD get params
	//  example: method
	//   if string or symbol then set
	//   else resolve to caller and call then verify again

	fmt.Printf("*** method: %s url: %s\n", method, url)

	return slip.List{nil, nil}
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

var validMethods = map[string]bool{
	http.MethodGet:     true,
	http.MethodPost:    true,
	http.MethodPut:     true,
	http.MethodDelete:  true,
	http.MethodHead:    true,
	http.MethodPatch:   true,
	http.MethodConnect: true,
	http.MethodOptions: true,
	http.MethodTrace:   true,
}
