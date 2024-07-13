// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"encoding/csv"
	"errors"
	"io"
	"os"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

var (
	readCSVActorFlavor *flavors.Flavor
)

func defReadCsvActor() {
	// Pkg.Initialize(nil)
	readCSVActorFlavor = flavors.DefFlavor("flow-read-csv-actor",
		map[string]slip.Object{},
		nil,
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`A flow-read-csv-actor reads a CSV or SEN file and parses the content. If the
file contains multiple CSV documents each will be placed in a separate _box_
and delivered to the task linked by the "ok" link.`),
			},
			slip.List{
				slip.Symbol(":init-keywords"),
				slip.Symbol(":filename"),
				slip.Symbol(":destination"),
				slip.Symbol(":count"),
				slip.Symbol(":separator"),
				slip.Symbol(":comment"),
				slip.Symbol(":trim"),
				slip.Symbol(":as-map"),
			},
		},
		&Pkg,
	)
	readCSVActorFlavor.DefMethod(":init", "", readCSVInitCaller{})
	readCSVActorFlavor.DefMethod(":start", "", readCSVActorStartCaller{})
	readCSVActorFlavor.DefMethod(":perform", "", readCSVActorPerformCaller{})
	readCSVActorFlavor.DefMethod(":init-key-values", "", readCSVActorInitKeyValuesCaller{})
}

type readCSVCtx struct {
	fileCtx
	sep     rune
	comment rune
	trim    bool
	asMap   bool
}

type readCSVInitCaller struct{}

func (caller readCSVInitCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	self := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		args = args[0].(slip.List)
	}
	var rcc readCSVCtx
	rcc.sep = ','
	rcc.parseArgs(s, args)
	for pos := 0; pos < len(args)-1; pos += 2 {
		sym, _ := args[pos].(slip.Symbol)
		switch string(sym) {
		case ":separator":
			if c, ok := args[pos+1].(slip.Character); ok {
				rcc.sep = rune(c)
			} else {
				slip.PanicType(":separator", args[pos+1], "character")
			}
		case ":comment":
			if c, ok := args[pos+1].(slip.Character); ok {
				rcc.comment = rune(c)
			} else {
				slip.PanicType(":comment", args[pos+1], "character")
			}
		case ":trim":
			rcc.trim = args[pos+1] != nil
		case ":as-map":
			rcc.asMap = args[pos+1] != nil
		}
	}
	self.Any = &rcc

	return nil
}

func (caller readCSVInitCaller) Docs() string {
	return `__:init__ &key _filename_ _destination_ _count_ _separator_ _comment_ _trim_ _as-map_
   _:filename_ [string|symbol|function] of the file to read.
   _:destination_ [string] the location in the _box_ to place the result.
   _:count_ [string] the location in the _box_ to place the count. If _nil_ then no count is set.
   _:separator_ [character] the field separator character, Default is comma.
   _:comment_ [character] the line comment character.
   _:trim_ [boolean] if true leading spaces are trimmed from each field.
   _:as-map_ [boolean] if true use the first line as the header and add each row as a map instead of a list.


Sets the initial value when _make-instance_ is called.
`
}

type readCSVActorStartCaller struct{}

func (caller readCSVActorStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any.(*readCSVCtx).task = args[0].(*flavors.Instance).Any.(*task)

	return nil
}

func (caller readCSVActorStartCaller) Docs() string {
	return `__:start__ _task_
   _:task_ [instance] the task that contains the actor.


Sets the context for the actor.
`
}

type readCSVActorPerformCaller struct{}

func (caller readCSVActorPerformCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	rcc := obj.Any.(*readCSVCtx)
	bi := args[0].(*flavors.Instance)

	filename := rcc.filename.value(s, bi)
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer func() { _ = f.Close() }()

	return rcc.readCSV(s, f, bi)
}

func (caller readCSVActorPerformCaller) Docs() string {
	return `__:perform__ _box_
   _:box_ [instance] the data to submit to the target flow.


Submits a box to the target flow.
`
}

type readCSVActorInitKeyValuesCaller struct{}

func (caller readCSVActorInitKeyValuesCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	rcc := obj.Any.(*readCSVCtx)
	var (
		dest  slip.Object
		count slip.Object
	)
	if rcc.dest != nil {
		dest = slip.String(jp.Expr(rcc.dest).String())
	}
	if rcc.count != nil {
		count = slip.String(jp.Expr(rcc.count).String())
	}
	kvs := slip.List{
		slip.Symbol(":filename"), rcc.filename.raw(),
		slip.Symbol(":destination"), dest,
		slip.Symbol(":count"), count,
		slip.Symbol(":separator"), slip.Character(rcc.sep),
	}
	if rcc.comment != 0 {
		kvs = append(kvs, slip.Symbol(":comment"), slip.Character(rcc.comment))
	}
	if rcc.trim {
		kvs = append(kvs, slip.Symbol(":trim"), slip.True)
	} else {
		kvs = append(kvs, slip.Symbol(":trim"), nil)
	}
	if rcc.asMap {
		kvs = append(kvs, slip.Symbol(":as-map"), slip.True)
	} else {
		kvs = append(kvs, slip.Symbol(":as-map"), nil)
	}
	return kvs
}

func (caller readCSVActorInitKeyValuesCaller) Docs() string {
	return `__:init-key-values__ => (:target "sub-flow")


Returns the keywords and values needed to recreate the instance as a property list.
`
}

func (rcc *readCSVCtx) readCSV(s *slip.Scope, r io.Reader, bi *flavors.Instance) slip.List {
	cr := csv.NewReader(r)
	cr.ReuseRecord = true
	cr.LazyQuotes = true
	if rcc.sep != 0 {
		cr.Comma = rcc.sep
	}
	cr.Comment = rcc.comment
	cr.TrimLeadingSpace = rcc.trim
	var (
		count  int64
		header []string
	)
	for {
		row, err := cr.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}
		var v any
		if rcc.asMap {
			if count == 0 {
				header = make([]string, len(row))
				copy(header, row)
				count++
				continue
			} else {
				m := map[string]any{}
				for i, f := range row {
					m[header[i]] = f
				}
				v = m
			}
		} else {
			list := make([]any, len(row))
			for i, f := range row {
				list[i] = f
			}
			v = list
		}
		var bx *box
		bi, bx = boxDup(bi)
		setBx(bx, v, jp.Expr(rcc.dest))
		if rcc.count != nil {
			setBx(bx, count, jp.Expr(rcc.count))
		}
		count++
		rcc.task.transition(s, "ok", bi)
	}
	return slip.List{nil, nil}
}
