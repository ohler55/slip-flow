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
		switch tv := sc.caller.Call(s, slip.List{bi}, 0).(type) {
		case slip.String:
			val = string(tv)
		case slip.Symbol:
			val = string(tv)
		default:
			slip.PanicType("value", tv, "string", "symbol")
		}
	}
	return
}
