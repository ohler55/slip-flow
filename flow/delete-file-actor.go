// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"os"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

var (
	deleteFileActorFlavor *flavors.Flavor
)

func defDeleteFileActor() {
	deleteFileActorFlavor = flavors.DefFlavor("flow-delete-file-actor",
		map[string]slip.Object{},
		nil,
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`Delete a file or files.`),
			},
			slip.List{
				slip.Symbol(":init-keywords"),
				slip.Symbol(":filename"),
			},
		},
		&Pkg,
	)
	deleteFileActorFlavor.DefMethod(":init", "", deleteFileInitCaller{})
	deleteFileActorFlavor.DefMethod(":start", "", deleteFileActorStartCaller{})
	deleteFileActorFlavor.DefMethod(":perform", "", deleteFileActorPerformCaller{})
	deleteFileActorFlavor.DefMethod(":init-key-values", "", deleteFileActorInitKeyValuesCaller{})
}

type deleteFileCtx struct {
	fileCtx
}

type deleteFileInitCaller struct{}

func (caller deleteFileInitCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	self := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		args = args[0].(slip.List)
	}
	var dfc deleteFileCtx
	dfc.parseArgs(s, args)
	self.Any = &dfc

	return nil
}

func (caller deleteFileInitCaller) FuncDocs() *slip.FuncDoc {
	return &slip.FuncDoc{
		Name: ":init",
		Text: `Sets the initial value when _make-instance_ is called.`,
		Args: []*slip.DocArg{
			{Name: "&key"},
			{
				Name: ":filename",
				Type: "string",
				Text: "Filename of the file to delete.",
			},
		},
	}
}

type deleteFileActorStartCaller struct{}

func (caller deleteFileActorStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any.(*deleteFileCtx).task = args[0].(*flavors.Instance).Any.(*task)

	return nil
}

func (caller deleteFileActorStartCaller) FuncDocs() *slip.FuncDoc {
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

type deleteFileActorPerformCaller struct{}

func (caller deleteFileActorPerformCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	dfc := obj.Any.(*deleteFileCtx)
	bi := args[0].(*flavors.Instance)

	_ = os.RemoveAll(dfc.filename.value(s, bi))

	return slip.List{slip.String("ok"), bi}
}

func (caller deleteFileActorPerformCaller) FuncDocs() *slip.FuncDoc {
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

type deleteFileActorInitKeyValuesCaller struct{}

func (caller deleteFileActorInitKeyValuesCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	dfc := obj.Any.(*deleteFileCtx)
	return slip.List{
		slip.Symbol(":filename"), dfc.filename.raw(),
	}
}

func (caller deleteFileActorInitKeyValuesCaller) FuncDocs() *slip.FuncDoc {
	return &slip.FuncDoc{
		Name: ":init-key-values",
		Text: `Returns the keywords and values needed to recreate the instance as a property
list. (e.g., (:filename "file.txt"))`,
		Return: "list",
	}
}
