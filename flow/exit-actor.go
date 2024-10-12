// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

var (
	exitActorFlavor *flavors.Flavor
)

func defExitActor() {
	exitActorFlavor = flavors.DefFlavor("flow-exit-actor",
		map[string]slip.Object{ // instance variables
			"notifiers": nil, // list of strings or symbols
		},
		nil,
		slip.List{
			slip.Symbol(":gettable-instance-variables"),
			slip.Symbol(":settable-instance-variables"),
			slip.Symbol(":inittable-instance-variables"),
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`An exit-actor is an actor that terminates or exits a flow. If the _box_
has a watcher the _box_ received is placed on the watcher channels that match the _notifiers_ or to
all watchers if no _notifiers_ have been identified.
`),
			},
		},
		&Pkg,
	)
	exitActorFlavor.DefMethod(":start", "", exitActorStartCaller{})
	exitActorFlavor.DefMethod(":perform", "", exitActorPerformCaller{})
	exitActorFlavor.DefMethod(":init-key-values", "", exitActorInitKeyValuesCaller{})
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

	var notifiers slip.List
	switch tn := obj.Get("notifiers").(type) {
	case nil:
		// leave as nil
	case slip.Symbol, slip.String:
		notifiers = slip.List{tn}
	case slip.List:
		notifiers = tn
	}
	if bi, ok := args[0].(*flavors.Instance); ok && bi.Type == boxFlavor {
		notifyBox(bi, notifiers)
	}
	obj.Any.(*flow).exit(args[0])

	return slip.List{nil, nil}
}

func (caller exitActorPerformCaller) Docs() string {
	return `__:perform__ _box_
   _:box_ [instance] the data to place on the watcher channels.


Place the _box_ on the box watcher channels.
`
}

type exitActorInitKeyValuesCaller struct{}

func (caller exitActorInitKeyValuesCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)

	return slip.List{slip.Symbol(":notifiers"), obj.Get("notifiers")}
}

func (caller exitActorInitKeyValuesCaller) Docs() string {
	return `__:init-key-values__ => (:notifiers ("chan1"))


Returns the keywords and values needed to recreate the instance as a property list.
`
}
