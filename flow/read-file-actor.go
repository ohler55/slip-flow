// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"os"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

var (
	readFileActorFlavor *flavors.Flavor
)

func defReadFileActor() {
	readFileActorFlavor = flavors.DefFlavor("flow-read-file-actor",
		map[string]slip.Object{},
		nil,
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`A flow-read-file-actor reads a file and creates a single string element from
the content.  The content string is placed in the _box_ at _destination_
before delivering to the task linked by the "ok" link.`),
			},
			slip.List{
				slip.Symbol(":init-keywords"),
				slip.Symbol(":filename"),
				slip.Symbol(":destination"),
			},
		},
		&Pkg,
	)
	readFileActorFlavor.DefMethod(":init", "", readFileInitCaller{})
	readFileActorFlavor.DefMethod(":start", "", readFileActorStartCaller{})
	readFileActorFlavor.DefMethod(":perform", "", readFileActorPerformCaller{})
	readFileActorFlavor.DefMethod(":init-key-values", "", readFileActorInitKeyValuesCaller{})
}

type readFileCtx struct {
	fileCtx
}

type readFileInitCaller struct{}

func (caller readFileInitCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	self := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		args = args[0].(slip.List)
	}
	var rfc readFileCtx
	rfc.parseArgs(s, args)

	self.Any = &rfc

	return nil
}

func (caller readFileInitCaller) FuncDocs() *slip.FuncDoc {
	return &slip.FuncDoc{
		Name: ":init",
		Text: `Sets the initial value when _make-instance_ is called.`,
		Args: []*slip.DocArg{
			{Name: "&key"},
			{
				Name: ":filename",
				Type: "string|symbol|function",
				Text: "Filename of the file to read.",
			},
			{
				Name: ":destination",
				Type: "string",
				Text: "Location in the _box_ to place the result.",
			},
		},
	}
}

type readFileActorStartCaller struct{}

func (caller readFileActorStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any.(*readFileCtx).task = args[0].(*flavors.Instance).Any.(*task)

	return nil
}

func (caller readFileActorStartCaller) FuncDocs() *slip.FuncDoc {
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

type readFileActorPerformCaller struct{}

func (caller readFileActorPerformCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	rfc := obj.Any.(*readFileCtx)
	bi := args[0].(*flavors.Instance)

	filename := rfc.filename.value(s, bi)
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer func() { _ = f.Close() }()

	return rfc.readText(f, bi)
}

func (caller readFileActorPerformCaller) FuncDocs() *slip.FuncDoc {
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

type readFileActorInitKeyValuesCaller struct{}

func (caller readFileActorInitKeyValuesCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	rfc := obj.Any.(*readFileCtx)
	var dest slip.Object
	if rfc.dest != nil {
		dest = slip.String(jp.Expr(rfc.dest).String())
	}
	return slip.List{
		slip.Symbol(":filename"), rfc.filename.raw(),
		slip.Symbol(":destination"), dest,
	}
}

func (caller readFileActorInitKeyValuesCaller) FuncDocs() *slip.FuncDoc {
	return &slip.FuncDoc{
		Name: ":init-key-values",
		Text: `Returns the keywords and values needed to recreate the instance
as a property list. (e.g., (:filename "file.txt"))`,
		Return: "list",
	}
}
