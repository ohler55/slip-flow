// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"fmt"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/bag"
	"github.com/ohler55/slip/pkg/flavors"
)

var (
	tailFileActorFlavor *flavors.Flavor
)

func init() {
	Pkg.Initialize(nil)
	tailFileActorFlavor = flavors.DefFlavor("flow-tail-file-actor",
		map[string]slip.Object{},
		nil,
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`A flow-tail-file-actor reads a file and parses the content according to the
specified format. A format of nil indicates a best guess will be made. If the
file contains multiple values such as a CSV file or JSON file each will be
placed in a separate _box_ and delivered to the task linked by the "ok" link.


If the _only-new_ option is true then the actor will only read new additions to the file
and create a new _box_ for each element as it is added. Paritial or incomplete elements of a JSON
or CSV will not trigger an error but instead will try again when the element is complete.`),
			},
			slip.List{
				slip.Symbol(":init-keywords"),
				slip.Symbol(":filename"),
				slip.Symbol(":format"),
				slip.Symbol(":destination"),
				slip.Symbol(":only-new"),
			},
		},
		&Pkg,
	)
	tailFileActorFlavor.DefMethod(":init", "", tailFileInitCaller{})
	tailFileActorFlavor.DefMethod(":start", "", tailFileActorStartCaller{})
	tailFileActorFlavor.DefMethod(":perform", "", tailFileActorPerformCaller{})
	tailFileActorFlavor.DefMethod(":init-key-values", "", tailFileActorInitKeyValuesCaller{})
}

type tailFileCtx struct {
	fileCtx
	filename string
	onlyNew  bool
}

type tailFileInitCaller struct{}

func (caller tailFileInitCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	self := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		args = args[0].(slip.List)
	}
	var rfc tailFileCtx
	for pos := 0; pos < len(args)-1; pos += 2 {
		sym, _ := args[pos].(slip.Symbol)
		switch string(sym) {
		case ":filename":
			if ss, ok := args[pos+1].(slip.String); ok {
				rfc.filename = string(ss)
			} else {
				slip.PanicType(":filename", args[pos+1], "string")
			}
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
		case ":only-new":
			rfc.onlyNew = args[pos+1] != nil
		}
	}
	self.Any = &rfc

	return nil
}

func (caller tailFileInitCaller) Docs() string {
	return `__:init__ &key _target_
   _:filename_ [string] of the file to read.
   _:format_ [symbol] the expected format: _nil_|_:text_|_:json_|_:csv_|_:xml_
   _:destination_ [string] the location in the _box_ to place the result.
   _:only-new_ [boolean] if true only new additions to the file will trigger a send.


Sets the initial value when _make-instance_ is called.
`
}

type tailFileActorStartCaller struct{}

func (caller tailFileActorStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	rfc := obj.Any.(*tailFileCtx)
	rfc.task = args[0].(*flavors.Instance).Any.(*task)

	// TBD if tail start a loop
	// add only-new options

	return nil
}

func (caller tailFileActorStartCaller) Docs() string {
	return `__:start__ _task_
   _:task_ [instance] the task that contains the actor.


Sets the context for the actor.
`
}

type tailFileActorPerformCaller struct{}

func (caller tailFileActorPerformCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	rfc := obj.Any.(*tailFileCtx)

	fmt.Printf("*** ctx: %v\n", rfc)

	// TBD

	return slip.List{nil, nil}
}

func (caller tailFileActorPerformCaller) Docs() string {
	return `__:perform__ _box_
   _:box_ [instance] the data to submit to the target flow.


Submits a box to the target flow.
`
}

type tailFileActorInitKeyValuesCaller struct{}

func (caller tailFileActorInitKeyValuesCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	rfc := obj.Any.(*tailFileCtx)
	var onlyNew slip.Object
	if rfc.onlyNew {
		onlyNew = slip.True
	}
	var dest slip.Object
	if rfc.dest != nil {
		dest = slip.String(rfc.dest.String())
	}
	return slip.List{
		slip.Symbol(":filename"), slip.String(rfc.filename),
		slip.Symbol(":format"), rfc.format,
		slip.Symbol(":destination"), dest,
		slip.Symbol(":only-new"), onlyNew,
	}
}

func (caller tailFileActorInitKeyValuesCaller) Docs() string {
	return `__:init-key-values__ => (:filename "log.json")


Returns the keywords and values needed to recreate the instance as a property list.
`
}
