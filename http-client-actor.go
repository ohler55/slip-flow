// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

import (
	"fmt"

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
	method  slip.Object
	url     slip.Object
	header  slip.Object
	trailer slip.Object
	body    slip.Object // string or function that returns a string or stream
	timeout slip.Object
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
			hcc.method = args[pos+1]
			// TBD verify either string, symbol or function
		case ":url":
			hcc.url = args[pos+1]
			// TBD verify either string or function
		case ":header":
			hcc.header = args[pos+1]
			// TBD verify either assoc or function
		case ":trailer":
			hcc.trailer = args[pos+1]
			// TBD verify either assoc or function
		case ":body":
			hcc.body = args[pos+1]
		case ":timeout":
			hcc.timeout = args[pos+1]
			// TBD verify either fixnum or function
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
   _reply-handler_ [function] to call with the response from a request. If _nil_ then
place the content in a "response" element of the box.


Each argument can be a function that takes a single argument that is the box
of the _:perform_. Is should return the expected type for the argument.
`
}

type httpClientActorStartCaller struct{}

func (caller httpClientActorStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	mc := obj.Any.(*httpClientCtx)
	mc.task = args[0].(*flavors.Instance).Any.(*task)

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

	mc := obj.Any.(*httpClientCtx)
	// TBD get params
	//  example: method
	//   if string or symbol then set
	//   else resolve to caller and call then verify again

	fmt.Printf("*** %s %v\n", bi, mc)

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
