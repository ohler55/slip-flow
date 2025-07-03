// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

var (
	globActorFlavor *flavors.Flavor
)

func defGlobActor() {
	globActorFlavor = flavors.DefFlavor("flow-glob-actor",
		map[string]slip.Object{},
		nil,
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`Places a list of the files matching the _pattern_ at the _destination_ in the
_box_.  If the _with-info_ option is true then each file entry will be a map
with the information includes with the keys of name, size, mode,
modified-time, and is-dir.`),
			},
			slip.List{
				slip.Symbol(":init-keywords"),
				slip.Symbol(":pattern"),
				slip.Symbol(":destination"),
				slip.Symbol(":with-info"),
			},
		},
		&Pkg,
	)
	globActorFlavor.DefMethod(":init", "", globInitCaller{})
	globActorFlavor.DefMethod(":start", "", globActorStartCaller{})
	globActorFlavor.DefMethod(":perform", "", globActorPerformCaller{})
	globActorFlavor.DefMethod(":init-key-values", "", globActorInitKeyValuesCaller{})
}

type globCtx struct {
	fileCtx
	info bool
}

type globInitCaller struct{}

func (caller globInitCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	self := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		args = args[0].(slip.List)
	}
	var rdc globCtx
	rdc.parseArgs(s, args)
	for pos := 0; pos < len(args)-1; pos += 2 {
		if slip.Symbol(":with-info") == args[pos] {
			rdc.info = args[pos+1] != nil
		}
	}
	self.Any = &rdc

	return nil
}

func (caller globInitCaller) FuncDocs() *slip.FuncDoc {
	return &slip.FuncDoc{
		Name: ":init",
		Text: `Sets the initial value when _make-instance_ is called.`,
		Args: []*slip.DocArg{
			{Name: "&key"},
			{
				Name: ":pattern",
				Type: "string|symbol|function",
				Text: "Glob pattern to find the files of.",
			},
			{
				Name: ":destination",
				Type: "string",
				Text: "Location in the _box_ to place the result.",
			},
			{
				Name: ":with-info",
				Type: "boolean",
				Text: "If true each entry is a map of information about the file.",
			},
		},
	}
}

type globActorStartCaller struct{}

func (caller globActorStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any.(*globCtx).task = args[0].(*flavors.Instance).Any.(*task)

	return nil
}

func (caller globActorStartCaller) FuncDocs() *slip.FuncDoc {
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

type globActorPerformCaller struct{}

func (caller globActorPerformCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	rdc := obj.Any.(*globCtx)
	bi := args[0].(*flavors.Instance)

	pattern := rdc.filename.value(s, bi)
	paths, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}
	list := make([]any, len(paths))
	if rdc.info {
		var fi fs.FileInfo
		for i, p := range paths {
			m := map[string]any{"name": p}
			if fi, err = os.Stat(p); err == nil {
				m["size"] = fi.Size()
				m["mode"] = fi.Mode().String()
				m["modified-time"] = fi.ModTime().UTC()
				m["is-dir"] = fi.IsDir()
			}
			list[i] = m
		}
	} else {
		for i, p := range paths {
			list[i] = p
		}
	}
	setBx(bi.Any.(*box), list, jp.Expr(rdc.dest))
	return slip.List{slip.String("ok"), bi}
}

func (caller globActorPerformCaller) FuncDocs() *slip.FuncDoc {
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

type globActorInitKeyValuesCaller struct{}

func (caller globActorInitKeyValuesCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	rdc := obj.Any.(*globCtx)
	var (
		dest slip.Object
	)
	if rdc.dest != nil {
		dest = slip.String(jp.Expr(rdc.dest).String())
	}
	kvs := slip.List{
		slip.Symbol(":pattern"), rdc.filename.raw(),
		slip.Symbol(":destination"), dest,
	}
	if rdc.info {
		kvs = append(kvs, slip.Symbol(":with-info"), slip.True)
	} else {
		kvs = append(kvs, slip.Symbol(":with-info"), nil)
	}
	return kvs
}

func (caller globActorInitKeyValuesCaller) FuncDocs() *slip.FuncDoc {
	return &slip.FuncDoc{
		Name: ":init-key-values",
		Text: `Returns the keywords and values needed to recreate the instance
as a property list. (e.g., (:pattern "dir/*"))`,
		Return: "list",
	}
}
