// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

var (
	splitActorFlavor *flavors.Flavor
)

func init() {
	splitActorFlavor = flavors.DefFlavor("flow-split-actor",
		map[string]slip.Object{},
		nil,
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`A flow-split-actor sends a _box_ on multiple links in parallel.
A _flow-merge-actor_ can be used to merge the branch of the slit back together.`),
			},
			slip.List{
				slip.Symbol(":init-keywords"),
				slip.Symbol(":links"),
			},
		},
	)
	splitActorFlavor.DefMethod(":init", "", splitInitCaller{})
	splitActorFlavor.DefMethod(":start", "", splitActorStartCaller{})
	splitActorFlavor.DefMethod(":perform", "", splitActorPerformCaller{})
}

type splitCtx struct {
	task  *task
	links []string
}

type splitInitCaller struct{}

func (caller splitInitCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	self := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		args = args[0].(slip.List)
	}
	var sc splitCtx
	for pos := 0; pos < len(args)-1; pos += 2 {
		if string(args[pos].(slip.Symbol)) == ":links" {
			names, ok := args[pos+1].(slip.List)
			if !ok {
				slip.PanicType("links", args[pos+1], "list")
			}
			for _, n := range names {
				switch tn := n.(type) {
				case slip.String:
					sc.links = append(sc.links, string(tn))
				case slip.Symbol:
					sc.links = append(sc.links, string(tn))
				default:
					slip.PanicType("links element", tn, "string", "symbol")
				}
			}
		}
	}
	self.Any = &sc

	return nil
}

func (caller splitInitCaller) Docs() string {
	return `__:init__ &key _links_
   _:links_ [list] a list of link names as either strings or symbols.


Sets the initial value when _make-instance_ is called.
`
}

type splitActorStartCaller struct{}

func (caller splitActorStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any.(*splitCtx).task = args[0].(*flavors.Instance).Any.(*task)

	return nil
}

func (caller splitActorStartCaller) Docs() string {
	return `__:start__ _task_
   _:task_ [instance] the task that contains the actor.


Sets the context for the actor.
`
}

type splitActorPerformCaller struct{}

func (caller splitActorPerformCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	bi := args[0].(*flavors.Instance)

	sc := obj.Any.(*splitCtx)
	tsk := sc.task
	for _, name := range sc.links {
		tsk.transition(s, name, bi)
	}
	return slip.List{nil, nil}
}

func (caller splitActorPerformCaller) Docs() string {
	return `__:perform__ _box_
   _:box_ [instance] the data to send on the configured links.


Send a _box_ on one or more links.
`
}
