// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

var (
	splitActorFlavor *flavors.Flavor
)

func defSplitActor() {
	splitActorFlavor = flavors.DefFlavor("flow-split-actor",
		map[string]slip.Object{},
		nil,
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`A flow-split-actor sends a _box_ on multiple links in parallel.
A _flow-merge-actor_ can be used to merge the branch of the slit back together.`),
			},
		},
		&Pkg,
	)
	splitActorFlavor.DefMethod(":start", "", splitActorStartCaller{})
	splitActorFlavor.DefMethod(":perform", "", splitActorPerformCaller{})
}

type splitActorStartCaller struct{}

func (caller splitActorStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any = args[0].(*flavors.Instance).Any

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

	tsk := obj.Any.(*task)
	for name := range tsk.links {
		if name != "error" {
			tsk.transition(s, name, bi)
		}
	}
	return slip.List{nil, nil}
}

func (caller splitActorPerformCaller) Docs() string {
	return `__:perform__ _box_
   _:box_ [instance] the data to send on the configured links.


Send a _box_ on one or more links.
`
}
