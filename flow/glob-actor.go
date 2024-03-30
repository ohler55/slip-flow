// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"fmt"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

var (
	globActorFlavor *flavors.Flavor
)

func init() {
	Pkg.Initialize(nil)
	globActorFlavor = flavors.DefFlavor("flow-read-directory-actor",
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

func (caller globInitCaller) Docs() string {
	return `__:init__ &key _pattern_ _destination_ _with-info_
   _:pattern_ [string|symbol|function] glob pattern to find the files of.
   _:destination_ [string] the location in the _box_ to place the result.
   _:with-info_ [boolean] if true each entry is a map of information about the file.


Sets the initial value when _make-instance_ is called.
`
}

type globActorStartCaller struct{}

func (caller globActorStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any.(*globCtx).task = args[0].(*flavors.Instance).Any.(*task)

	return nil
}

func (caller globActorStartCaller) Docs() string {
	return `__:start__ _task_
   _:task_ [instance] the task that contains the actor.


Sets the context for the actor.
`
}

type globActorPerformCaller struct{}

func (caller globActorPerformCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	rdc := obj.Any.(*globCtx)
	bi := args[0].(*flavors.Instance)

	pattern := rdc.filename.value(s, bi)

	fmt.Printf("*** pathname: %s\n", pattern)

	// TBD  filepath.Glob(pattern)

	return slip.List{slip.String("ok"), bi}
}

func (caller globActorPerformCaller) Docs() string {
	return `__:perform__ _box_
   _:box_ [instance] the data to submit to the target flow.


Submits a box to the target flow.
`
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

func (caller globActorInitKeyValuesCaller) Docs() string {
	return `__:init-key-values__ => (:pattern "dir/*")


Returns the keywords and values needed to recreate the instance as a property list.
`
}
