// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

import (
	"fmt"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/pretty"
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
			slip.List{
				slip.Symbol(":init-keywords"),
				slip.Symbol(":pretty"),
				slip.Symbol(":depth"),
				slip.Symbol(":right-margin"),
				slip.Symbol(":indent"),
				slip.Symbol(":time-format"),
				slip.Symbol(":time-wrap"),
				slip.Symbol(":json"),
				slip.Symbol(":color"),
				slip.Symbol(":full"),
			},
		},
	)
	logErrorActorFlavor.DefMethod(":init", "", logErrorInitCaller{})
	logErrorActorFlavor.DefMethod(":start", "", logErrorActorStartCaller{})
	logErrorActorFlavor.DefMethod(":perform", "", logErrorActorPerformCaller{})
}

type logErrorCtx struct {
	task *task
	pw   *pretty.Writer
	full bool
	prty bool
}

type logErrorInitCaller struct{}

func (caller logErrorInitCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	self := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		args = args[0].(slip.List)
	}
	var lec logErrorCtx
	lec.pw, lec.full, lec.prty, _ = parseBoxWriteOptions(args, false)
	self.Any = &lec

	return nil
}

func (caller logErrorInitCaller) Docs() string {
	return `__:init__ &key _pretty_ _depth_ _right-margin_ _indent_ _time-format_ _time-wrap_ _json_ _color_ _full_
   _:pretty_ [boolean] value to use in place of the _*print-pretty*_ value.
If _t_ then the JSON or SEN output is indented according to the other keyword options.
   _:depth [fixnum] maximum number of nested elements on a line in the output.
A value of zero outputs a tight single line output. Default: 4.
   _:right-margin_ [fixnum] value to use in place of the _*print-right-margin*_ value.
   _:indent_ [fixnum] is the number of spaces to indent JSON or SEN output if :pretty is not non-nil.
   _:time-format_ [string] value to use in place of the _*flow-box-time-format*_ value.
   _:time-wrap_ [string] value to use in place of the _*flow-box-time-wrap*_ value.
   _:json_ [boolean] if true the output is JSON formatted otherwise output is SEN format.
   _:color_ [boolean] if true the output is colorized.
   _:full_ [boolean] if true the output includes the box track and the content is nested on level down.


Sets the initial value when _make-instance_ is called.
`
}

type logErrorActorStartCaller struct{}

func (caller logErrorActorStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any.(*logErrorCtx).task = args[0].(*flavors.Instance).Any.(*task)

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
	lec := obj.Any.(*logErrorCtx)
	tsk := lec.task
	if msg, _ := jp.C("error").First(bx.content).(string); 0 < len(msg) {
		ev := bx.track.history[len(bx.track.history)-2]
		tsk.self.Receive(s,
			":error",
			slip.List{slip.String(fmt.Sprintf("%s:%s %s - %s", ev.flow, ev.task, bx.track.id, msg))}, 0)
	} else {
		b := bx.toString(lec.pw, lec.full, lec.prty)
		tsk.self.Receive(s, ":error", slip.List{slip.String(b)}, 0)
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
