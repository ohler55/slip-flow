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
	// Pkg.Initialize(nil)
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

func (caller readFileInitCaller) Docs() string {
	return `__:init__ &key _filename_ _destination_
   _:filename_ [string|symbol|function] of the file to read.
   _:destination_ [string] the location in the _box_ to place the result.


Sets the initial value when _make-instance_ is called.
`
}

type readFileActorStartCaller struct{}

func (caller readFileActorStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any.(*readFileCtx).task = args[0].(*flavors.Instance).Any.(*task)

	return nil
}

func (caller readFileActorStartCaller) Docs() string {
	return `__:start__ _task_
   _:task_ [instance] the task that contains the actor.


Sets the context for the actor.
`
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

func (caller readFileActorPerformCaller) Docs() string {
	return `__:perform__ _box_
   _:box_ [instance] the data to submit to the target flow.


Submits a box to the target flow.
`
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

func (caller readFileActorInitKeyValuesCaller) Docs() string {
	return `__:init-key-values__ => (:filename "file.txt")


Returns the keywords and values needed to recreate the instance as a property list.
`
}
