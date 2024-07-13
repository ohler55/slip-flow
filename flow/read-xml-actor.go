// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"os"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

var (
	readXMLActorFlavor *flavors.Flavor
)

func defReadXmlActor() {
	readXMLActorFlavor = flavors.DefFlavor("flow-read-xml-actor",
		map[string]slip.Object{},
		nil,
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`A flow-read-xml-actor reads a XML or SEN file and parses the content. If the
file contains multiple XML documents each will be placed in a separate _box_
and delivered to the task linked by the "ok" link.`),
			},
			slip.List{
				slip.Symbol(":init-keywords"),
				slip.Symbol(":filename"),
				slip.Symbol(":destination"),
				slip.Symbol(":count"),
				slip.Symbol(":strict"),
				slip.Symbol(":html"),
				slip.Symbol(":trim"),
			},
		},
		&Pkg,
	)
	readXMLActorFlavor.DefMethod(":init", "", readXMLInitCaller{})
	readXMLActorFlavor.DefMethod(":start", "", readXMLActorStartCaller{})
	readXMLActorFlavor.DefMethod(":perform", "", readXMLActorPerformCaller{})
	readXMLActorFlavor.DefMethod(":init-key-values", "", readXMLActorInitKeyValuesCaller{})
}

type readXMLCtx struct {
	fileCtx
	trim   bool
	html   bool
	strict bool
}

type readXMLInitCaller struct{}

func (caller readXMLInitCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	self := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		args = args[0].(slip.List)
	}
	var rxc readXMLCtx
	rxc.parseArgs(s, args)
	for pos := 0; pos < len(args)-1; pos += 2 {
		sym, _ := args[pos].(slip.Symbol)
		switch string(sym) {
		case ":strict":
			rxc.strict = args[pos+1] != nil
		case ":trim":
			rxc.trim = args[pos+1] != nil
		case ":html":
			rxc.html = args[pos+1] != nil
		}
	}
	self.Any = &rxc

	return nil
}

func (caller readXMLInitCaller) Docs() string {
	return `__:init__ &key _filename_ _destination_ _count_ _strict_ _trim_ _html_
   _:filename_ [string|symbol|function] of the file to read.
   _:destination_ [string] the location in the _box_ to place the result.
   _:count_ [string] the location in the _box_ to place the count. If _nil_ then no count is set.
   _:strict_ [boolean] if true, the default, strict XML parsing is used.
   _:trim_ [boolean] if true white space is removed and empty fields are removed.
   _:html_ [boolean] if true parsing will handle HTML. The default is _nil.


Sets the initial value when _make-instance_ is called.
`
}

type readXMLActorStartCaller struct{}

func (caller readXMLActorStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any.(*readXMLCtx).task = args[0].(*flavors.Instance).Any.(*task)

	return nil
}

func (caller readXMLActorStartCaller) Docs() string {
	return `__:start__ _task_
   _:task_ [instance] the task that contains the actor.


Sets the context for the actor.
`
}

type readXMLActorPerformCaller struct{}

func (caller readXMLActorPerformCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	rxc := obj.Any.(*readXMLCtx)
	bi := args[0].(*flavors.Instance)

	filename := rxc.filename.value(s, bi)
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer func() { _ = f.Close() }()

	return rxc.readXML(s, f, bi)
}

func (caller readXMLActorPerformCaller) Docs() string {
	return `__:perform__ _box_
   _:box_ [instance] the data to submit to the target flow.


Submits a box to the target flow.
`
}

type readXMLActorInitKeyValuesCaller struct{}

func (caller readXMLActorInitKeyValuesCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	rxc := obj.Any.(*readXMLCtx)
	var (
		dest  slip.Object
		count slip.Object
	)
	if rxc.dest != nil {
		dest = slip.String(jp.Expr(rxc.dest).String())
	}
	if rxc.count != nil {
		count = slip.String(jp.Expr(rxc.count).String())
	}
	kvs := slip.List{
		slip.Symbol(":filename"), rxc.filename.raw(),
		slip.Symbol(":destination"), dest,
		slip.Symbol(":count"), count,
	}
	if rxc.strict {
		kvs = append(kvs, slip.Symbol(":strict"), slip.True)
	} else {
		kvs = append(kvs, slip.Symbol(":strict"), nil)
	}
	if rxc.trim {
		kvs = append(kvs, slip.Symbol(":trim"), slip.True)
	} else {
		kvs = append(kvs, slip.Symbol(":trim"), nil)
	}
	if rxc.html {
		kvs = append(kvs, slip.Symbol(":html"), slip.True)
	} else {
		kvs = append(kvs, slip.Symbol(":html"), nil)
	}
	return kvs
}

func (caller readXMLActorInitKeyValuesCaller) Docs() string {
	return `__:init-key-values__ => (:target "sub-flow")


Returns the keywords and values needed to recreate the instance as a property list.
`
}

func (rxc *readXMLCtx) readXML(s *slip.Scope, r io.Reader, bi *flavors.Instance) slip.List {
	dec := xml.NewDecoder(r)
	dec.Strict = rxc.strict
	if rxc.html {
		dec.Strict = false
		dec.AutoClose = xml.HTMLAutoClose
		dec.Entity = xml.HTMLEntity
	}
	var (
		stack   [][]any
		element []any
		count   int64
	)
out:
	for {
		token, err := dec.Token()
		switch tt := token.(type) {
		case nil:
			if errors.Is(err, io.EOF) {
				break out
			}
			panic(err)
		case xml.StartElement:
			if element != nil {
				stack = append(stack, element)
			}
			attrs := map[string]any{}
			for _, a := range tt.Attr {
				attrs[a.Name.Local] = a.Value
			}
			element = []any{tt.Name.Local, attrs}
		case xml.EndElement:
			stack[len(stack)-1] = append(stack[len(stack)-1], element)
			element = stack[len(stack)-1]
			stack[len(stack)-1] = nil
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				var bx *box
				bi, bx = boxDup(bi)
				setBx(bx, element, jp.Expr(rxc.dest))
				if rxc.count != nil {
					setBx(bx, count, jp.Expr(rxc.count))
				}
				count++
				rxc.task.transition(s, "ok", bi)
				element = nil
			}
		case xml.CharData:
			if rxc.trim {
				tt = bytes.TrimSpace(tt)
			}
			if 0 < len(tt) {
				element = append(element, string(tt))
			}
		case xml.Comment:
			element = append(element, []any{":comment", string(tt)})
		case xml.Directive:
			element = append(element, []any{":directive", string(tt)})
		case xml.ProcInst:
			element = append(element, []any{":processing-instruction", tt.Target, string(tt.Inst)})
		}
	}
	return slip.List{nil, nil}
}
