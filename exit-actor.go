// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

var (
	exitActorFlavor *flavors.Flavor
)

func init() {
	exitActorFlavor = flavors.DefFlavor("flow-exit-actor",
		map[string]slip.Object{},
		nil,
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`An exit-actor is an actor that terminates or exits a flow. If the _flow_
_exit-channel_ has been set then then _box_ received is placed on the
_exit-channel_.
`),
			},
		},
	)
	exitActorFlavor.DefMethod(":start", "", exitActorStartCaller{})
	exitActorFlavor.DefMethod(":perform", "", exitActorPerformCaller{})
}

type exitActorStartCaller struct{}

func (caller exitActorStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any = args[0].(*flavors.Instance).Any.(*task).flow
	return nil
}

func (caller exitActorStartCaller) Docs() string {
	return `__:start__ _task_
   _:task_ [instance] the task that contains the actor.


Sets the context for the actor.
`
}

type exitActorPerformCaller struct{}

func (caller exitActorPerformCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)

	obj.Any.(*flow).exit(args[0])

	return slip.List{nil, nil}
}

func (caller exitActorPerformCaller) Docs() string {
	return `__:perform__ _box_
   _:box_ [instance] the data to place on the flow exit-channel.


Place the _box_ on the flow exit-channel is the exit-channel is not nil.
`
}
