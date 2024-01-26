// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

var (
	taskActorFlavor *flavors.Flavor
)

func init() {
	taskActorFlavor = flavors.DefFlavor("flow-task-actor",
		map[string]slip.Object{"task": nil},
		nil,
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`A flow-task-actor sets the _task_ variable with the provided task argument.`),
			},
			slip.Symbol(":gettable-instance-variables"),
			slip.List{
				slip.Symbol(":init-keywords"),
				slip.Symbol(":links"),
			},
		},
	)
	taskActorFlavor.DefMethod(":start", "", taskActorStartCaller{})
}

type taskActorStartCaller struct{}

func (caller taskActorStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Set("task", args[0].(*flavors.Instance).Any.(*task).self)

	return nil
}

func (caller taskActorStartCaller) Docs() string {
	return `__:start__ _task_
   _:task_ [instance] the task that contains the actor.


Sets the task variable for the actor.
`
}
