// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/cl"
)

type listCaller struct {
	list   slip.List
	caller slip.Caller
}

func (lc *listCaller) extract(s *slip.Scope, arg slip.Object) {
	switch ta := arg.(type) {
	case slip.List:
		lc.list = ta
	default:
		lc.caller = cl.ResolveToCaller(s, arg, 0)
	}
}

func (lc *listCaller) value(s *slip.Scope, bi slip.Object, depth int) (val slip.List) {
	val = lc.list
	if lc.caller != nil {
		v := lc.caller.Call(s, slip.List{bi}, 0)
		var ok bool
		if val, ok = v.(slip.List); !ok {
			slip.TypePanic(s, depth, "list", v, "list")
		}
	}
	return
}

func (lc *listCaller) raw() (rv slip.Object) {
	var ok bool
	if rv, ok = lc.caller.(*slip.Lambda); !ok {
		rv = lc.list
	}
	return
}
