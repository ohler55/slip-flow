// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

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

func defLogErrorActor() {
	logErrorActorFlavor = flavors.DefFlavor("flow-log-error-actor",
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
				slip.String(`A flow-log-error-actor is an actor that logs an error and exits the flow if
no links are attached. If the _box_ has a watcher are no attached links
then then _box_ received is placed on the watcher channel.
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
		&Pkg,
	)
	logErrorActorFlavor.DefMethod(":init", "", logErrorInitCaller{})
	logErrorActorFlavor.DefMethod(":start", "", logErrorActorStartCaller{})
	logErrorActorFlavor.DefMethod(":perform", "", logErrorActorPerformCaller{})
	logErrorActorFlavor.DefMethod(":init-key-values", "", logErrorActorInitKeyValuesCaller{})
}

type logErrorCtx struct {
	task *task
	pw   *pretty.Writer
	full bool
	prty bool
}

type logErrorInitCaller struct{}

func (caller logErrorInitCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	self := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		args = args[0].(slip.List)
	}
	var lec logErrorCtx
	lec.pw, lec.full, lec.prty, _ = parseBoxWriteOptions(s, args, false, depth)
	self.Any = &lec

	return nil
}

func (caller logErrorInitCaller) FuncDocs() *slip.FuncDoc {
	return &slip.FuncDoc{
		Name: ":init",
		Text: `Sets the initial value when _make-instance_ is called.`,
		Args: []*slip.DocArg{
			{Name: "&key"},
			{
				Name: ":pretty",
				Type: "boolean",
				Text: `Value to use in place of the _*print-pretty*_ value.
If _t_ then the JSON or SEN output is indented according to the other keyword options.`,
			},
			{
				Name: ":depth",
				Type: "fixnum",
				Text: `The maximum number of nested elements on a line in the output.
A value of zero outputs a tight single line output.`,
				Default: slip.Fixnum(4),
			},
			{
				Name: ":right-margin",
				Type: "fixnum",
				Text: "The value to use in place of the _*print-right-margin*_ value.",
			},
			{
				Name: ":indent",
				Type: "fixnum",
				Text: "The number of spaces to indent JSON or SEN output if :pretty is not non-nil.",
			},
			{
				Name: ":time-format",
				Type: "string",
				Text: "The value to use in place of the _*flow-box-time-format*_ value.",
			},
			{
				Name: ":time-wrap",
				Type: "string",
				Text: "The value to use in place of the _*flow-box-time-wrap*_ value.",
			},
			{
				Name: ":json",
				Type: "boolean",
				Text: "If true the output is JSON formatted otherwise output is SEN format.",
			},
			{
				Name: ":color",
				Type: "boolean",
				Text: "If true the output is colorized.",
			},
			{
				Name: ":full",
				Type: "boolean",
				Text: "If true the output includes the box track and the content is nested on level down.",
			},
		},
	}
}

type logErrorActorStartCaller struct{}

func (caller logErrorActorStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any.(*logErrorCtx).task = args[0].(*flavors.Instance).Any.(*task)

	return nil
}

func (caller logErrorActorStartCaller) FuncDocs() *slip.FuncDoc {
	return &slip.FuncDoc{
		Name: ":start",
		Text: `Sets the context for the actor.`,
		Args: []*slip.DocArg{
			{
				Name: "task",
				Type: "task",
				Text: "The task that contains the actor.",
			},
		},
	}
}

type logErrorActorPerformCaller struct{}

func (caller logErrorActorPerformCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
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
		notifyBox(s, bi, notifiers, depth)
	}
	tsk.flow.exit(args[0])

	return slip.List{nil, nil}
}

func (caller logErrorActorPerformCaller) FuncDocs() *slip.FuncDoc {
	return &slip.FuncDoc{
		Name: ":perform",
		Text: `Log the box error message or the content and then place the _box_ on the watcher channels.`,
		Args: []*slip.DocArg{
			{
				Name: "box",
				Type: "box",
				Text: "The data to log and then place on any watcher channels.",
			},
		},
	}
}

type logErrorActorInitKeyValuesCaller struct{}

func (caller logErrorActorInitKeyValuesCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	lec := obj.Any.(*logErrorCtx)

	kvs := slip.List{
		slip.Symbol(":notifiers"), obj.Get("notifiers"),
		slip.Symbol(":depth"), slip.Fixnum(lec.pw.MaxDepth),
		slip.Symbol(":right-margin"), slip.Fixnum(lec.pw.Width),
		slip.Symbol(":indent"), slip.Fixnum(lec.pw.Indent),
		slip.Symbol(":time-format"), slip.String(lec.pw.TimeFormat),
		slip.Symbol(":time-wrap"), slip.String(lec.pw.TimeWrap),
	}
	if lec.pw.SEN {
		kvs = append(kvs, slip.Symbol(":json"), nil)
	} else {
		kvs = append(kvs, slip.Symbol(":json"), slip.True)
	}
	if lec.pw.Color {
		kvs = append(kvs, slip.Symbol(":color"), slip.True)
	} else {
		kvs = append(kvs, slip.Symbol(":color"), nil)
	}
	if lec.prty {
		kvs = append(kvs, slip.Symbol(":pretty"), slip.True)
	} else {
		kvs = append(kvs, slip.Symbol(":pretty"), nil)
	}
	if lec.full {
		kvs = append(kvs, slip.Symbol(":full"), slip.True)
	} else {
		kvs = append(kvs, slip.Symbol(":full"), nil)
	}
	return kvs
}

func (caller logErrorActorInitKeyValuesCaller) FuncDocs() *slip.FuncDoc {
	return &slip.FuncDoc{
		Name: ":init-key-values",
		Text: `Returns the keywords and values needed to recreate the instance
as a property list. (e.g., (:full t))`,
		Return: "list",
	}
}
