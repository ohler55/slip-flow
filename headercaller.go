// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

import (
	"net/http"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/cl"
)

type headerCaller struct {
	header http.Header
	caller slip.Caller
}

func (hc *headerCaller) extract(s *slip.Scope, arg slip.Object) {
	if list, ok := arg.(slip.List); ok {
		hc.header = hc.extractFromList(list)
	} else {
		hc.caller = cl.ResolveToCaller(s, arg, 0)
	}
}

func (hc *headerCaller) value(s *slip.Scope, bi slip.Object) (val http.Header) {
	val = hc.header
	if hc.caller != nil {
		result := hc.caller.Call(s, slip.List{bi}, 0)
		if list, ok := result.(slip.List); ok {
			val = hc.extractFromList(list)
		} else {
			slip.PanicType("value", result, "assoc")
		}
	}
	return
}

func (hc *headerCaller) extractFromList(list slip.List) http.Header {
	header := http.Header{}
	for _, v := range list {
		cons, ok := v.(slip.List)
		if !ok {
			slip.PanicType("header element", v, "cons")
		}
		header.Add(mustBeString(cons.Car()), mustBeString(cons.Cdr()))
	}
	return header
}

func mustBeString(arg slip.Object) (str string) {
	switch ta := arg.(type) {
	case slip.String:
		str = string(ta)
	case slip.Symbol:
		str = string(ta)
	default:
		slip.PanicType("string", arg, "string", "symbol")
	}
	return
}
