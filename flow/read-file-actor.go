// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"os"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/bag"
	"github.com/ohler55/slip/pkg/flavors"
)

var (
	readFileActorFlavor *flavors.Flavor
)

func init() {
	Pkg.Initialize(nil)
	readFileActorFlavor = flavors.DefFlavor("flow-read-file-actor",
		map[string]slip.Object{},
		nil,
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`A flow-read-file-actor reads a file and parses the content according to the
specified format. A format of nil indicates a best guess will be made. If the
file contains multiple values such as a CSV file or JSON file each will be
placed in a separate _box_ and delivered to the task linked by the "ok" link.`),
			},
			slip.List{
				slip.Symbol(":init-keywords"),
				slip.Symbol(":filename"),
				slip.Symbol(":format"),
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
	filename strCaller
}

type readFileInitCaller struct{}

func (caller readFileInitCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	self := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		args = args[0].(slip.List)
	}
	var rfc readFileCtx
	for pos := 0; pos < len(args)-1; pos += 2 {
		sym, _ := args[pos].(slip.Symbol)
		switch string(sym) {
		case ":filename":
			rfc.filename.extract(s, args[pos+1])
		case ":format":
			switch ta := args[pos+1].(type) {
			case nil:
			case slip.Symbol:
				if ta == slip.Symbol(":text") ||
					ta == slip.Symbol(":json") ||
					ta == slip.Symbol(":csv") ||
					ta == slip.Symbol(":xml") {
					rfc.format = ta
				} else {
					slip.PanicType(":format", args[pos+1], "nil", ":text", ":json", ":csv", ":xml")
				}
			}
		case ":destination":
			switch ta := args[pos+1].(type) {
			case slip.String:
				rfc.dest = bag.Path(jp.MustParse([]byte(ta)))
			case slip.Symbol:
				rfc.dest = bag.Path(jp.MustParse([]byte(ta)))
			default:
				slip.PanicType(":destination", args[pos+1], "string", "symbol")
			}
		}
	}
	self.Any = &rfc

	return nil
}

func (caller readFileInitCaller) Docs() string {
	return `__:init__ &key _target_
   _:filename_ [string|symbol|function] of the file to read.
   _:format_ [symbol] the expected format: _nil_|_:text_|_:json_|_:csv_|_:xml_
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
	result := slip.List{nil, nil}
	switch rfc.format {
	case nil:
		result = rfc.readAuto(f, bi)
	case slip.Symbol(":text"):
		result = rfc.readText(f, bi)
	case slip.Symbol(":json"):
		result = rfc.readJSON(s, f, bi)
	case slip.Symbol(":csv"):
		result = rfc.readCSV(f, bi)
	case slip.Symbol(":xml"):
		result = rfc.readXML(f, bi)
	}
	return result
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
		dest = slip.String(rfc.dest.String())
	}
	return slip.List{
		slip.Symbol(":filename"), rfc.filename.raw(),
		slip.Symbol(":format"), rfc.format,
		slip.Symbol(":destination"), dest,
	}
}

func (caller readFileActorInitKeyValuesCaller) Docs() string {
	return `__:init-key-values__ => (:target "sub-flow")


Returns the keywords and values needed to recreate the instance as a property list.
`
}
