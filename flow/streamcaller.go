// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"io"
	"strings"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/cl"
)

type streamCaller struct {
	str    string
	caller slip.Caller
}

func (sc *streamCaller) extract(s *slip.Scope, arg slip.Object) {
	switch ta := arg.(type) {
	case slip.String:
		sc.str = string(ta)
	default:
		sc.caller = cl.ResolveToCaller(s, arg, 0)
	}
}

func (sc *streamCaller) value(s *slip.Scope, bi slip.Object, depth int) (val io.Reader) {
	if sc.caller != nil {
		switch tv := sc.caller.Call(s, slip.List{bi}, 0).(type) {
		case slip.String:
			val = &slip.InputStream{Reader: strings.NewReader(string(tv))}
		case io.Reader:
			val = tv
		default:
			slip.TypePanic(s, depth, "value", tv, "string", "symbol")
		}
	} else {
		val = &slip.InputStream{Reader: strings.NewReader(sc.str)}
	}
	return
}

func (sc *streamCaller) raw() (rv slip.Object) {
	var ok bool
	if rv, ok = sc.caller.(*slip.Lambda); !ok {
		rv = slip.String(sc.str)
	}
	return
}
