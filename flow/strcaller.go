// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/cl"
)

type strCaller struct {
	str    string
	caller slip.Caller
}

func (sc *strCaller) extract(s *slip.Scope, arg slip.Object) {
	switch ta := arg.(type) {
	case slip.String:
		sc.str = string(ta)
	case slip.Symbol:
		sc.str = string(ta)
	default:
		sc.caller = cl.ResolveToCaller(s, arg, 0)
	}
}

func (sc *strCaller) value(s *slip.Scope, bi slip.Object, depth int) (val string) {
	val = sc.str
	if sc.caller != nil {
		val = mustBeString(s, sc.caller.Call(s, slip.List{bi}, 0), "value", depth)
	}
	return
}

func (sc *strCaller) raw() (rv slip.Object) {
	var ok bool
	if rv, ok = sc.caller.(*slip.Lambda); !ok {
		rv = slip.String(sc.str)
	}
	return
}

func mustBeString(s *slip.Scope, arg slip.Object, name string, depth int) (str string) {
	switch ta := arg.(type) {
	case slip.String:
		str = string(ta)
	case slip.Symbol:
		str = string(ta)
	default:
		slip.TypePanic(s, depth, name, arg, "string", "symbol")
	}
	return
}
