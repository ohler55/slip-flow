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
	task   *task
	format slip.Object
	dest   bag.Path
}

func (fc *fileCtx) readText(r io.Reader, bi *flavors.Instance) slip.List {
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
	if bi == nil {
		bi, _ = MakeBox(nil)
	}
	setBox(bi, slip.String(content), fc.dest)

	return slip.List{slip.String("ok"), bi}
}

func (fc *fileCtx) readJSON(s *slip.Scope, r io.Reader, bi *flavors.Instance) slip.List {
	_ = sen.MustParseReader(r, func(v any) {
		var bx *box
		bi, bx = boxDup(bi)
		if fc.dest == nil {
			bx.content = v
		} else {
			jp.Expr(fc.dest).MustSet(bx.content, v)
		}
		fc.task.transition(s, "ok", bi)
	})
	return slip.List{nil, nil}
}

func (fc *fileCtx) readCSV(r io.Reader, bi *flavors.Instance) slip.List {

	// TBD
	return slip.List{nil, nil}
}

func (fc *fileCtx) readXML(r io.Reader, bi *flavors.Instance) slip.List {

	// TBD
	return slip.List{nil, nil}
}

func (fc *fileCtx) readAuto(r io.Reader, bi *flavors.Instance) slip.List {

	// TBD
	return slip.List{nil, nil}
}
