// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

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

func (sc *strCaller) value(s *slip.Scope, bi slip.Object) (val string) {
	val = sc.str
	if sc.caller != nil {
		val = mustBeString(sc.caller.Call(s, slip.List{bi}, 0), "value")
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

func mustBeString(arg slip.Object, name string) (str string) {
	switch ta := arg.(type) {
	case slip.String:
		str = string(ta)
	case slip.Symbol:
		str = string(ta)
	default:
		slip.PanicType(name, arg, "string", "symbol")
	}
	return
}
