// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

import (
	"time"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/cl"
)

type durCaller struct {
	dur    time.Duration
	caller slip.Caller
}

func (dc *durCaller) extract(s *slip.Scope, arg slip.Object) {
	switch ta := arg.(type) {
	case nil:
		dc.dur = 0
	case slip.String:
		var err error
		if dc.dur, err = time.ParseDuration(string(ta)); err != nil {
			panic(err)
		}
	case slip.Symbol:
		var err error
		if dc.dur, err = time.ParseDuration(string(ta)); err != nil {
			panic(err)
		}
	case slip.Fixnum:
		dc.dur = time.Duration(ta) * time.Second
	default:
		dc.caller = cl.ResolveToCaller(s, arg, 0)
	}
}

func (dc *durCaller) value(s *slip.Scope, bi slip.Object) (val time.Duration) {
	val = dc.dur
	if dc.caller != nil {
		switch tv := dc.caller.Call(s, slip.List{bi}, 0).(type) {
		case nil:
			dc.dur = 0
		case slip.String:
			var err error
			if dc.dur, err = time.ParseDuration(string(tv)); err != nil {
				panic(err)
			}
		case slip.Symbol:
			var err error
			if dc.dur, err = time.ParseDuration(string(tv)); err != nil {
				panic(err)
			}
		case slip.Fixnum:
			dc.dur = time.Duration(tv) * time.Second
		default:
			slip.PanicType("value", tv, "fixnum", "string", "symbol")
		}
	}
	return
}
