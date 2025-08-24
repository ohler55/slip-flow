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

func (caller readXMLInitCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	self := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		args = args[0].(slip.List)
	}
	var rxc readXMLCtx
	rxc.parseArgs(s, args, depth)
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

func (caller readXMLInitCaller) FuncDocs() *slip.FuncDoc {
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
			{
				Name: ":count",
				Type: "string",
				Text: "The location in the _box_ to place the count. If _nil_ then no count is set.",
			},
			{
				Name: ":strict",
				Type: "boolean",
				Text: "If true, the default, strict XML parsing is used.",
			},
			{
				Name: ":trim",
				Type: "boolean",
				Text: "If true white space is removed and empty fields are removed.",
			},
			{
				Name: ":html",
				Type: "boolean",
				Text: "If true parsing will handle HTML. The default is _nil.",
			},
		},
	}
}

type readXMLActorStartCaller struct{}

func (caller readXMLActorStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any.(*readXMLCtx).task = args[0].(*flavors.Instance).Any.(*task)

	return nil
}

func (caller readXMLActorStartCaller) FuncDocs() *slip.FuncDoc {
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

type readXMLActorPerformCaller struct{}

func (caller readXMLActorPerformCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	rxc := obj.Any.(*readXMLCtx)
	bi := args[0].(*flavors.Instance)

	filename := rxc.filename.value(s, bi, depth)
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer func() { _ = f.Close() }()

	return rxc.readXML(s, f, bi)
}

func (caller readXMLActorPerformCaller) FuncDocs() *slip.FuncDoc {
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

func (caller readXMLActorInitKeyValuesCaller) FuncDocs() *slip.FuncDoc {
	return &slip.FuncDoc{
		Name: ":init-key-values",
		Text: `Returns the keywords and values needed to recreate the instance
as a property list. (e.g., (:target "sub-flow"))`,
		Return: "list",
	}
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
