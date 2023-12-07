// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main

import (
	"time"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

var (
	trackFlavor *flavors.Flavor
)

func init() {
	trackFlavor = flavors.DefFlavor("flow-track-flavor",
		map[string]slip.Object{},
		nil,
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`A container for tracking information as an instance of the
_flow-box-flavor_ is traverses a flow.`),
			},
		},
	)
	trackFlavor.Final = true
	trackFlavor.GoMakeOnly = true
	trackFlavor.DefMethod(":id", "", trackIDCaller{})
	trackFlavor.DefMethod(":events", "", trackEventsCaller{})
}

type event struct {
	when time.Time
	flow string
	task string
}

type track struct {
	id     slip.Object
	events []*event
}

type trackIDCaller struct{}

func (caller trackIDCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	self := s.Get("self").(*flavors.Instance)
	return self.Any.(*track).id
}

func (caller trackIDCaller) Docs() string {
	return `__:id__ => _string_|_fixnum_


Returns the id of the track.
`
}

type trackEventsCaller struct{}

func (caller trackEventsCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	self := s.Get("self").(*flavors.Instance)
	return self.Any.(*track).eventsList()
}

func (caller trackEventsCaller) Docs() string {
	return `__:events__ => _list_


Returns the events of the track as a list of triples where each triple is a list of
the time, the task name, and the flow name.
`
}

// Scan adds an event to the track.
func (t *track) Scan(flow, task string) {
	t.events = append(t.events, &event{flow: flow, task: task, when: time.Now().UTC()})
}

func (t *track) eventsList() (events slip.List) {
	if 0 < len(t.events) {
		events = make(slip.List, len(t.events))
		for i, ev := range t.events {
			events[i] = slip.List{slip.Time(ev.when), slip.String(ev.task), slip.String(ev.flow)}
		}
	}
	return
}

// MakeTrack is only public for testing purposes.
func MakeTrack(id slip.Object) (self *flavors.Instance, t *track) {
	self = trackFlavor.MakeInstance().(*flavors.Instance)
	t = &track{id: id}
	self.Any = t

	return
}
