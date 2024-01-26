// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

import (
	"fmt"
	"time"

	"github.com/ohler55/slip"
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
				slip.Symbol(":number"),
				slip.Symbol(":timeout"),
			},
		},
	)
	httpClientActorFlavor.DefMethod(":init", "", httpClientInitCaller{})
	httpClientActorFlavor.DefMethod(":start", "", httpClientActorStartCaller{})
	httpClientActorFlavor.DefMethod(":perform", "", httpClientActorPerformCaller{})
}

type httpClientCtx struct {
	task    *task
	timeout time.Duration
	// TBD
}

type httpClientInitCaller struct{}

func (caller httpClientInitCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	self := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		args = args[0].(slip.List)
	}
	mc := httpClientCtx{timeout: time.Second * 10}
	for pos := 0; pos < len(args)-1; pos += 2 {
		sym, _ := args[pos].(slip.Symbol)
		switch string(sym) {
		case ":timeout":
			if num, ok := args[pos+1].(slip.Fixnum); ok && 0 < num {
				mc.timeout = time.Duration(num) * time.Second
			} else {
				slip.PanicType(":timeout", args[pos+1], "positive fixnum")
			}
		}
	}
	self.Any = &mc

	return nil
}

func (caller httpClientInitCaller) Docs() string {
	return `__:init__ &key _timeout_ _number_
   _:timeout_ [fixnum] seconds before timing out waiting for _number_ of boxes.
   _:number_ [fixnum] number of boxes expected before transitioning.


TBD
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
	// TBD
	fmt.Printf("*** %s %v\n", bi, mc)

	return slip.List{nil, nil}
}

func (caller httpClientActorPerformCaller) Docs() string {
	return `__:perform__ _box_
   _:box_ [instance] box TBD


TBD
`
}
