// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"os"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

var (
	readJSONActorFlavor *flavors.Flavor
)

func init() {
	Pkg.Initialize(nil)
	readJSONActorFlavor = flavors.DefFlavor("flow-read-json-actor",
		map[string]slip.Object{},
		nil,
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`A flow-read-json-actor reads a JSON or SEN file and parses the content. If the
file contains multiple JSON documents each will be placed in a separate _box_
and delivered to the task linked by the "ok" link.`),
			},
			slip.List{
				slip.Symbol(":init-keywords"),
				slip.Symbol(":filename"),
				slip.Symbol(":destination"),
				slip.Symbol(":count"),
			},
		},
		&Pkg,
	)
	readJSONActorFlavor.DefMethod(":init", "", readJSONInitCaller{})
	readJSONActorFlavor.DefMethod(":start", "", readJSONActorStartCaller{})
	readJSONActorFlavor.DefMethod(":perform", "", readJSONActorPerformCaller{})
	readJSONActorFlavor.DefMethod(":init-key-values", "", readJSONActorInitKeyValuesCaller{})
}

type readJSONCtx struct {
	fileCtx
}

type readJSONInitCaller struct{}

func (caller readJSONInitCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	self := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		args = args[0].(slip.List)
	}
	var rjc readJSONCtx
	rjc.parseArgs(s, args)
	self.Any = &rjc

	return nil
}

func (caller readJSONInitCaller) Docs() string {
	return `__:init__ &key _filename_ _destination_ _count_
   _:filename_ [string|symbol|function] of the file to read.
   _:destination_ [string] the location in the _box_ to place the result.
   _:count_ [string] the location in the _box_ to place the count. If _nil_ then no count is set.


Sets the initial value when _make-instance_ is called.
`
}

type readJSONActorStartCaller struct{}

func (caller readJSONActorStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any.(*readJSONCtx).task = args[0].(*flavors.Instance).Any.(*task)

	return nil
}

func (caller readJSONActorStartCaller) Docs() string {
	return `__:start__ _task_
   _:task_ [instance] the task that contains the actor.


Sets the context for the actor.
`
}

type readJSONActorPerformCaller struct{}

func (caller readJSONActorPerformCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	rjc := obj.Any.(*readJSONCtx)
	bi := args[0].(*flavors.Instance)

	filename := rjc.filename.value(s, bi)
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer func() { _ = f.Close() }()

	return rjc.readJSON(s, f, bi)
}

func (caller readJSONActorPerformCaller) Docs() string {
	return `__:perform__ _box_
   _:box_ [instance] the data to submit to the target flow.


Submits a box to the target flow.
`
}

type readJSONActorInitKeyValuesCaller struct{}

func (caller readJSONActorInitKeyValuesCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	rjc := obj.Any.(*readJSONCtx)
	var (
		dest  slip.Object
		count slip.Object
	)
	if rjc.dest != nil {
		dest = slip.String(jp.Expr(rjc.dest).String())
	}
	if rjc.count != nil {
		count = slip.String(jp.Expr(rjc.count).String())
	}
	return slip.List{
		slip.Symbol(":filename"), rjc.filename.raw(),
		slip.Symbol(":destination"), dest,
		slip.Symbol(":count"), count,
	}
}

func (caller readJSONActorInitKeyValuesCaller) Docs() string {
	return `__:init-key-values__ => (:target "sub-flow")


Returns the keywords and values needed to recreate the instance as a property list.
`
}
