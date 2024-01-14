// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

import (
	"fmt"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

var (
	logErrorActorFlavor *flavors.Flavor
)

func init() {
	logErrorActorFlavor = flavors.DefFlavor("flow-log-error-actor",
		map[string]slip.Object{},
		nil,
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`A flow-log-error-actor is an actor that logs an error and exits the flow if
no links are attached. If the _flow_ _log-error-channel_ has been and set there are no attached links
then then _box_ received is placed on the _log-error-channel_.
`),
			},
		},
	)
	logErrorActorFlavor.DefMethod(":start", "", logErrorActorStartCaller{})
	logErrorActorFlavor.DefMethod(":perform", "", logErrorActorPerformCaller{})
}

type logErrorActorStartCaller struct{}

func (caller logErrorActorStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any = args[0].(*flavors.Instance).Any // task
	return nil
}

func (caller logErrorActorStartCaller) Docs() string {
	return `__:start__ _task_
   _:task_ [instance] the task that contains the actor.


Sets the context for the actor.
`
}

type logErrorActorPerformCaller struct{}

func (caller logErrorActorPerformCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	bx := args[0].(*flavors.Instance).Any.(*box)
	tsk := obj.Any.(*task)
	if msg, _ := jp.C("error").First(bx.content).(string); 0 < len(msg) {
		ev := bx.track.history[len(bx.track.history)-2]
		tsk.self.Receive(s,
			":error",
			slip.List{slip.String(fmt.Sprintf("%s:%s %s - %s", ev.flow, ev.task, bx.track.id, msg))}, 0)
	} else {
		history := make([]any, len(bx.track.history))
		for i, ev := range bx.track.history {
			history[i] = map[string]any{
				"when": ev.when,
				"flow": ev.flow,
				"task": ev.task,
			}
		}
		content := map[string]any{
			"track": map[string]any{
				"id":      slip.Simplify(bx.track.id),
				"history": history,
			},
			"content": bx.content,
		}
		tsk.self.Receive(s, ":error", slip.List{slip.String(sen.Bytes(content))}, 0)
	}
	for linkName := range tsk.links {
		return slip.List{slip.String(linkName), args[0]}
	}
	tsk.flow.exit(args[0])

	return slip.List{nil, nil}
}

func (caller logErrorActorPerformCaller) Docs() string {
	return `__:perform__ _box_
   _:box_ [instance] the data to log and then place on the flow exit-channel.


Log the box error message or the content and then place the _box_ on the flow exit-channel
is the log-error-channel is not nil.
`
}
