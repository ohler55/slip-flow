// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"io"

	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

var (
	inspectActorFlavor *flavors.Flavor
)

func init() {
	Pkg.Initialize(nil)
	inspectActorFlavor = flavors.DefFlavor("flow-inspect-actor",
		map[string]slip.Object{},
		nil,
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`A flow-inspect-actor is an actor that logs or prints the content of a box.`),
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
				slip.Symbol(":output"),
			},
		},
		&Pkg,
	)
	inspectActorFlavor.DefMethod(":init", "", inspectInitCaller{})
	inspectActorFlavor.DefMethod(":start", "", inspectActorStartCaller{})
	inspectActorFlavor.DefMethod(":perform", "", inspectActorPerformCaller{})
	inspectActorFlavor.DefMethod(":links", "", inspectActorLinksCaller{})
}

type inspectCtx struct {
	task   *task
	pw     *pretty.Writer
	full   bool
	prty   bool
	output slip.Object
}

type inspectInitCaller struct{}

func (caller inspectInitCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	self := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		args = args[0].(slip.List)
	}
	var ic inspectCtx
	ic.pw, ic.full, ic.prty, ic.output = parseBoxWriteOptions(args, true)
	self.Any = &ic

	return nil
}

func (caller inspectInitCaller) Docs() string {
	return `__:init__ &key _pretty_ _depth_ _right-margin_ _indent_ _time-format_ _time-wrap_ _json_ _color_ _full_
   _:pretty_ [boolean] value to use in place of the _*print-pretty*_ value.
If _t_ then the JSON or SEN output is indented according to the other keyword options.
   _:depth_ [fixnum] maximum number of nested elements on a line in the output.
A value of zero outputs a tight single line output. Default: 4.
   _:right-margin_ [fixnum] value to use in place of the _*print-right-margin*_ value.
   _:indent_ [fixnum] is the number of spaces to indent JSON or SEN output if :pretty is not non-nil.
   _:time-format_ [string] value to use in place of the _*flow-box-time-format*_ value.
   _:time-wrap_ [string] value to use in place of the _*flow-box-time-wrap*_ value.
   _:json_ [boolean] if true the output is JSON formatted otherwise output is SEN format.
   _:color_ [boolean] if true the output is colorized.
   _:full_ [boolean] if true the output includes the box track and the content is nested on level down.
   _:output_ [nil|symbol] if nil the box is written to _*standard-output*_ otherwise the _output_ must be
_:error_, _:warn_, _:info_, or _:debug_ matching the logger methods and filtered accordingly.


Sets the initial value when _make-instance_ is called.
`
}

type inspectActorStartCaller struct{}

func (caller inspectActorStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any.(*inspectCtx).task = args[0].(*flavors.Instance).Any.(*task)

	return nil
}

func (caller inspectActorStartCaller) Docs() string {
	return `__:start__ _task_
   _:task_ [instance] the task that contains the actor.


Sets the context for the actor.
`
}

type inspectActorPerformCaller struct{}

func (caller inspectActorPerformCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	bx := args[0].(*flavors.Instance).Any.(*box)
	ic := obj.Any.(*inspectCtx)
	tsk := ic.task
	b := bx.toString(ic.pw, ic.full, ic.prty)
	switch level := ic.output.(type) {
	case nil:
		_, _ = s.Get("*standard-output*").(io.Writer).Write(b)
	case slip.Symbol:
		tsk.self.Receive(s, string(level), slip.List{slip.String(b)}, 0)
	default:
		slip.PanicType("output", level, "nil", "symbol")
	}
	return slip.List{slip.String("ok"), args[0]}
}

func (caller inspectActorPerformCaller) Docs() string {
	return `__:perform__ _box_
   _:box_ [instance] the data to log.


Write the box to either _*standard-output*_ or to the logger.
`
}

type inspectActorLinksCaller struct{}

func (caller inspectActorLinksCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	return slip.List{slip.String("ok")}
}

func (caller inspectActorLinksCaller) Docs() string {
	return `__:links__ => _list_


Returns a list of transitions that can be followed.
`
}
