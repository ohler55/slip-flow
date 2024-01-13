// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

var (
	groupFlavor *flavors.Flavor
)

func init() {
	groupFlavor = flavors.DefFlavor("flow-group-flavor",
		map[string]slip.Object{},
		[]string{
			"can-log-flavor",
		},
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`Groups are the access point for a collection of _flows_.
`),
			},
		},
	)
	groupFlavor.DefMethod(":init", "", groupInitCaller{})
	groupFlavor.DefMethod(":add", "", groupAddCaller{})
	groupFlavor.DefMethod(":find", "", groupFindCaller{})
	groupFlavor.DefMethod(":remove", "", groupRemoveCaller{})
	groupFlavor.DefMethod(":flows", "", groupFlowsCaller{})
	groupFlavor.DefMethod(":start", "", groupStartCaller{})
	groupFlavor.DefMethod(":shutdown", "", groupShutdownCaller{})
	groupFlavor.DefMethod(":set-level", ":after", groupSetLevelCaller{})
}

type group struct {
	self  *flavors.Instance
	flows map[string]*flow
}

func (g *group) add(obj slip.Object) {
	if fi, _ := obj.(*flavors.Instance); fi != nil && fi.Flavor == flowFlavor {
		g.flows[fi.Any.(*flow).name] = fi.Any.(*flow)
		fi.Any.(*flow).group = g
	} else {
		slip.PanicType("flow", obj, "flow")
	}
}

func (g *group) remove(obj slip.Object) {
	var name string
	switch tn := obj.(type) {
	case slip.String:
		name = string(tn)
	case slip.Symbol:
		name = string(tn)
	default:
		slip.PanicType("group :remove flow", tn, "string", "symbol")
	}
	if f := g.flows[name]; f != nil {
		f.group = nil
	}
	delete(g.flows, name)
}

func (g *group) find(obj slip.Object) (fi slip.Object) {
	var name string
	switch tn := obj.(type) {
	case slip.String:
		name = string(tn)
	case slip.Symbol:
		name = string(tn)
	default:
		slip.PanicType("group :find flow", tn, "string", "symbol")
	}
	if f := g.flows[name]; f != nil {
		fi = f.self
	}
	return
}

func (g *group) flowList() slip.List {
	all := make(slip.List, 0, len(g.flows))
	for _, f := range g.flows {
		all = append(all, f.self)
	}
	return all
}

func (g *group) start(s *slip.Scope) {
	logger := g.self.Get("logger")
	if logger == nil {
		logger = slip.ReadString("(make-instance 'logger-flavor)").Eval(s, nil)
		g.self.Set("logger", logger)
	}
	for _, f := range g.flows {
		f.self.Set("logger", logger)
		f.start(s)
	}
}

func (g *group) shutdown(s *slip.Scope) {
	for _, f := range g.flows {
		f.shutdown(s)
	}
}

type groupInitCaller struct{}

func (caller groupInitCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any = &group{self: obj, flows: map[string]*flow{}}

	return nil
}

func (caller groupInitCaller) Docs() string {
	return `__:init__ &key _logger_
   _:logger_ [instance] an instance that has the _:log_ method.


Sets the initial value when _make-instance_ is called.
`
}

type groupSetLevelCaller struct{}

func (caller groupSetLevelCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	g := obj.Any.(*group)
	level := g.self.Get("log-level")
	for _, f := range g.flows {
		_ = f.self.Receive(s, ":set-level", slip.List{level}, 0)
	}
	return nil
}

func (caller groupSetLevelCaller) Docs() string {
	return `__:after :setLevel__


Sets the _log-level_ of all the flows in the group.
`
}
