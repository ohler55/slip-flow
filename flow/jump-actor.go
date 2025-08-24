// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

var (
	jumpActorFlavor *flavors.Flavor
)

func defJumpActor() {
	jumpActorFlavor = flavors.DefFlavor("flow-jump-actor",
		map[string]slip.Object{},
		nil,
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`A flow-jump-actor jump out of the current flow and sends the box to
the entry of another flow in the current group.`),
			},
			slip.List{
				slip.Symbol(":init-keywords"),
				slip.Symbol(":target"),
			},
		},
		&Pkg,
	)
	jumpActorFlavor.DefMethod(":init", "", jumpInitCaller{})
	jumpActorFlavor.DefMethod(":start", "", jumpActorStartCaller{})
	jumpActorFlavor.DefMethod(":perform", "", jumpActorPerformCaller{})
	jumpActorFlavor.DefMethod(":init-key-values", "", jumpActorInitKeyValuesCaller{})
}

type jumpCtx struct {
	task   *task
	target string
}

type jumpInitCaller struct{}

func (caller jumpInitCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	self := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		args = args[0].(slip.List)
	}
	var jc jumpCtx
	for pos := 0; pos < len(args)-1; pos += 2 {
		sym := args[pos].(slip.Symbol)
		if string(sym) == ":target" {
			switch ta := args[pos+1].(type) {
			case slip.String:
				jc.target = string(ta)
			case slip.Symbol:
				jc.target = string(ta)
			default:
				slip.TypePanic(s, depth, ":target", args[pos+1], "string", "symbol")
			}
		}
	}
	self.Any = &jc

	return nil
}

func (caller jumpInitCaller) FuncDocs() *slip.FuncDoc {
	return &slip.FuncDoc{
		Name: ":init",
		Text: `Sets the initial value when _make-instance_ is called.`,
		Args: []*slip.DocArg{
			{Name: "&key"},
			{
				Name: ":target",
				Type: "string|symbol",
				Text: "The name of the target flow to jump to.",
			},
		},
	}
}

type jumpActorStartCaller struct{}

func (caller jumpActorStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any.(*jumpCtx).task = args[0].(*flavors.Instance).Any.(*task)

	return nil
}

func (caller jumpActorStartCaller) FuncDocs() *slip.FuncDoc {
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

type jumpActorPerformCaller struct{}

func (caller jumpActorPerformCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	jc := obj.Any.(*jumpCtx)
	if jc.task.flow == nil || jc.task.flow.group == nil {
		slip.ErrorPanic(s, depth, "task is not in a flow that is in a group")
	}
	f := jc.task.flow.group.flows[jc.target]
	if f == nil {
		slip.ErrorPanic(s, depth, "flow %s is not in the same group that %s is in", jc.target, jc.task.flow.name)
	}
	f.submit(s, args[0], nil, depth)

	return slip.List{nil, nil}
}

func (caller jumpActorPerformCaller) FuncDocs() *slip.FuncDoc {
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

type jumpActorInitKeyValuesCaller struct{}

func (caller jumpActorInitKeyValuesCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	jc := obj.Any.(*jumpCtx)

	return slip.List{slip.Symbol(":target"), slip.String(jc.target)}
}

func (caller jumpActorInitKeyValuesCaller) FuncDocs() *slip.FuncDoc {
	return &slip.FuncDoc{
		Name: ":init-key-values",
		Text: `Returns the keywords and values needed to recreate the instance
as a property list. (e.g., (:target "sub-flow"))`,
		Return: "list",
	}
}
