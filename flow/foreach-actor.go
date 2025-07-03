// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/bag"
	"github.com/ohler55/slip/pkg/flavors"
)

var (
	foreachActorFlavor *flavors.Flavor
)

func defForeachActor() {
	foreachActorFlavor = flavors.DefFlavor("flow-foreach-actor",
		map[string]slip.Object{},
		nil,
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`A flow-foreach-actor iterates over a list and generates a new _box_ that is
sent on the "ok" link for each element of the list.`),
			},
			slip.List{
				slip.Symbol(":init-keywords"),
				slip.Symbol(":list"),
				slip.Symbol(":destination"),
			},
		},
		&Pkg,
	)
	foreachActorFlavor.DefMethod(":init", "", foreachInitCaller{})
	foreachActorFlavor.DefMethod(":start", "", foreachActorStartCaller{})
	foreachActorFlavor.DefMethod(":perform", "", foreachActorPerformCaller{})
	foreachActorFlavor.DefMethod(":init-key-values", "", foreachActorInitKeyValuesCaller{})
}

type foreachCtx struct {
	task *task
	list listCaller
	dest bag.Path
}

type foreachInitCaller struct{}

func (caller foreachInitCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	self := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		args = args[0].(slip.List)
	}
	var fec foreachCtx
	for pos := 0; pos < len(args)-1; pos += 2 {
		sym, _ := args[pos].(slip.Symbol)
		switch string(sym) {
		case ":list":
			fec.list.extract(s, args[pos+1])
		case ":destination":
			switch ta := args[pos+1].(type) {
			case slip.String:
				fec.dest = bag.Path(jp.MustParse([]byte(ta)))
			case slip.Symbol:
				fec.dest = bag.Path(jp.MustParse([]byte(ta)))
			default:
				slip.PanicType(":destination", args[pos+1], "string", "symbol")
			}
		}
	}
	self.Any = &fec

	return nil
}

func (caller foreachInitCaller) FuncDocs() *slip.FuncDoc {
	return &slip.FuncDoc{
		Name: ":init",
		Text: `Sets the initial value when _make-instance_ is called.`,
		Args: []*slip.DocArg{
			{Name: "&key"},
			{
				Name: ":list",
				Type: "list|function",
				Text: "List to iterate over.",
			},
			{
				Name: ":destination",
				Type: "string",
				Text: "Location in the _box_ to place the each list member.",
			},
		},
	}
}

type foreachActorStartCaller struct{}

func (caller foreachActorStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any.(*foreachCtx).task = args[0].(*flavors.Instance).Any.(*task)

	return nil
}

func (caller foreachActorStartCaller) FuncDocs() *slip.FuncDoc {
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

type foreachActorPerformCaller struct{}

func (caller foreachActorPerformCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	fec := obj.Any.(*foreachCtx)
	bi := args[0].(*flavors.Instance)

	list := fec.list.value(s, bi)
	for _, v := range list {
		var bx *box
		bi, bx = boxDup(bi)
		setBx(bx, v, jp.Expr(fec.dest))
		fec.task.transition(s, "ok", bi)
	}
	return slip.List{nil, nil}
}

func (caller foreachActorPerformCaller) FuncDocs() *slip.FuncDoc {
	return &slip.FuncDoc{
		Name: ":perform",
		Text: `Submits a box to the target flow.`,
		Args: []*slip.DocArg{
			{
				Name: "box",
				Type: "box",
				Text: "The data to submit to the target flow.",
			},
		},
	}
}

type foreachActorInitKeyValuesCaller struct{}

func (caller foreachActorInitKeyValuesCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	fec := obj.Any.(*foreachCtx)
	var dest slip.Object
	if fec.dest != nil {
		dest = slip.String(jp.Expr(fec.dest).String())
	}
	return slip.List{
		slip.Symbol(":list"), fec.list.raw(),
		slip.Symbol(":destination"), dest,
	}
}

func (caller foreachActorInitKeyValuesCaller) FuncDocs() *slip.FuncDoc {
	return &slip.FuncDoc{
		Name: ":init-key-values",
		Text: `Returns the keywords and values needed to recreate the instance
as a property list. (e.g., (:list (1 3 5)))`,
		Return: "list",
	}
}
