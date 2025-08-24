// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"net/http"
	"sort"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/cl"
)

type headerCaller struct {
	header http.Header
	caller slip.Caller
}

func (hc *headerCaller) extract(s *slip.Scope, arg slip.Object, depth int) {
	if list, ok := arg.(slip.List); ok {
		hc.header = hc.extractFromList(s, list, depth)
	} else {
		hc.caller = cl.ResolveToCaller(s, arg, 0)
	}
}

func (hc *headerCaller) value(s *slip.Scope, bi slip.Object, depth int) (val http.Header) {
	val = hc.header
	if hc.caller != nil {
		result := hc.caller.Call(s, slip.List{bi}, 0)
		if list, ok := result.(slip.List); ok {
			val = hc.extractFromList(s, list, depth)
		} else {
			slip.TypePanic(s, depth, "value", result, "assoc")
		}
	}
	return
}

func (hc *headerCaller) extractFromList(s *slip.Scope, list slip.List, depth int) http.Header {
	header := http.Header{}
	for _, v := range list {
		cons, ok := v.(slip.List)
		if !ok {
			slip.TypePanic(s, depth, "header element", v, "cons")
		}
		header.Add(mustBeString(s, cons.Car(), "header key", depth), mustBeString(s, cons.Cdr(), "header value", depth))
	}
	return header
}

func (hc *headerCaller) raw() (rv slip.Object) {
	if lam, ok := hc.caller.(*slip.Lambda); ok {
		rv = lam
	} else {
		var assoc slip.List
		keys := make([]string, 0, len(hc.header))
		for k := range hc.header {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			kv := slip.List{slip.String(k)}
			for _, v := range hc.header[k] {
				kv = append(kv, slip.String(v))
			}
			assoc = append(assoc, kv)
		}
		if 0 < len(assoc) {
			rv = assoc
		}
	}
	return
}
