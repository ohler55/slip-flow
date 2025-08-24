// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"errors"
	"io"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/bag"
	"github.com/ohler55/slip/pkg/flavors"
)

type fileCtx struct {
	task     *task
	filename strCaller
	dest     bag.Path
	count    bag.Path
}

func (fc *fileCtx) parseArgs(s *slip.Scope, args slip.List, depth int) {
	for pos := 0; pos < len(args)-1; pos += 2 {
		sym, _ := args[pos].(slip.Symbol)
		switch string(sym) {
		case ":filename", ":pattern":
			fc.filename.extract(s, args[pos+1])
		case ":destination":
			switch ta := args[pos+1].(type) {
			case slip.String:
				fc.dest = bag.Path(jp.MustParse([]byte(ta)))
			case slip.Symbol:
				fc.dest = bag.Path(jp.MustParse([]byte(ta)))
			default:
				slip.TypePanic(s, depth, ":destination", args[pos+1], "string", "symbol")
			}
		case ":count":
			switch ta := args[pos+1].(type) {
			case slip.String:
				fc.count = bag.Path(jp.MustParse([]byte(ta)))
			case slip.Symbol:
				fc.count = bag.Path(jp.MustParse([]byte(ta)))
			default:
				slip.TypePanic(s, depth, ":count", args[pos+1], "string", "symbol")
			}
		}
	}
}

func (fc *fileCtx) readText(s *slip.Scope, r io.Reader, bi *flavors.Instance, depth int) slip.List {
	var content []byte
	buf := make([]byte, 4096)
	for {
		cnt, err := r.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}
		content = append(content, buf[:cnt]...)
	}
	setBox(s, bi, slip.String(content), fc.dest, depth)

	return slip.List{slip.String("ok"), bi}
}

func (fc *fileCtx) readJSON(s *slip.Scope, r io.Reader, bi *flavors.Instance) slip.List {
	var count int64
	_ = sen.MustParseReader(r, func(v any) {
		var bx *box
		bi, bx = boxDup(bi)
		setBx(bx, v, jp.Expr(fc.dest))
		if fc.count != nil {
			setBx(bx, count, jp.Expr(fc.count))
		}
		count++
		fc.task.transition(s, "ok", bi)
	})
	return slip.List{nil, nil}
}
