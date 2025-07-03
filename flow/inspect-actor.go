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

func defInspectActor() {
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
	inspectActorFlavor.DefMethod(":init-key-values", "", inspectActorInitKeyValuesCaller{})
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

func (caller inspectInitCaller) FuncDocs() *slip.FuncDoc {
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
			{
				Name: ":output",
				Type: "symbol|nil",
				Text: `If nil the box is written to _*standard-output*_ otherwise the _output_ must be
_:error_, _:warn_, _:info_, or _:debug_ matching the logger methods and filtered accordingly.`,
			},
		},
	}
}

type inspectActorStartCaller struct{}

func (caller inspectActorStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any.(*inspectCtx).task = args[0].(*flavors.Instance).Any.(*task)

	return nil
}

func (caller inspectActorStartCaller) FuncDocs() *slip.FuncDoc {
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

func (caller inspectActorPerformCaller) FuncDocs() *slip.FuncDoc {
	return &slip.FuncDoc{
		Name: ":perform",
		Text: `Write the box to either _*standard-output*_ or to the logger.`,
		Args: []*slip.DocArg{
			{
				Name: "box",
				Type: "box",
				Text: "The data to log.",
			},
		},
	}
}

type inspectActorLinksCaller struct{}

func (caller inspectActorLinksCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	return slip.List{slip.String("ok")}
}

func (caller inspectActorLinksCaller) FuncDocs() *slip.FuncDoc {
	return &slip.FuncDoc{
		Name:   ":links",
		Text:   `Returns a list of transitions that can be followed.`,
		Return: "list",
	}
}

type inspectActorInitKeyValuesCaller struct{}

func (caller inspectActorInitKeyValuesCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	ic := obj.Any.(*inspectCtx)

	kvs := slip.List{
		slip.Symbol(":depth"), slip.Fixnum(ic.pw.MaxDepth),
		slip.Symbol(":right-margin"), slip.Fixnum(ic.pw.Width),
		slip.Symbol(":indent"), slip.Fixnum(ic.pw.Indent),
		slip.Symbol(":time-format"), slip.String(ic.pw.TimeFormat),
		slip.Symbol(":time-wrap"), slip.String(ic.pw.TimeWrap),
		slip.Symbol(":output"), ic.output,
	}
	if ic.pw.SEN {
		kvs = append(kvs, slip.Symbol(":json"), nil)
	} else {
		kvs = append(kvs, slip.Symbol(":json"), slip.True)
	}
	if ic.pw.Color {
		kvs = append(kvs, slip.Symbol(":color"), slip.True)
	} else {
		kvs = append(kvs, slip.Symbol(":color"), nil)
	}
	if ic.prty {
		kvs = append(kvs, slip.Symbol(":pretty"), slip.True)
	} else {
		kvs = append(kvs, slip.Symbol(":pretty"), nil)
	}
	if ic.full {
		kvs = append(kvs, slip.Symbol(":full"), slip.True)
	} else {
		kvs = append(kvs, slip.Symbol(":full"), nil)
	}
	return kvs
}

func (caller inspectActorInitKeyValuesCaller) FuncDocs() *slip.FuncDoc {
	return &slip.FuncDoc{
		Name: ":init-key-values",
		Text: `Returns the keywords and values needed to recreate the instance
as a property list. (e.g., (:full t))`,
		Return: "list",
	}
}
