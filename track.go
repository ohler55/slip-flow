// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

import (
	"fmt"
	"sort"
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
	trackFlavor.DefMethod(":history", "", trackHistoryCaller{})
	trackFlavor.DefMethod(":merge", "", trackMergeCaller{})
}

type event struct {
	when time.Time
	flow string
	task string
}

type track struct {
	id      slip.Object
	history []*event
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

type trackHistoryCaller struct{}

func (caller trackHistoryCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	self := s.Get("self").(*flavors.Instance)
	return self.Any.(*track).historyList()
}

func (caller trackHistoryCaller) Docs() string {
	return `__:history__ => _list_


Returns the history of the track as a list of triples where each triple is a list of
the time, the task name, and the flow name.
`
}

type trackMergeCaller struct{}

func (caller trackMergeCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	self := s.Get("self").(*flavors.Instance)
	ti, _ := args[0].(*flavors.Instance)
	if ti == nil || ti.Flavor != trackFlavor {
		slip.PanicType("other", args[0], "track")
	}
	self.Any.(*track).merge(ti.Any.(*track))

	return self
}

func (caller trackMergeCaller) Docs() string {
	return `__:merge__ _other_ => _track_


Merge one track into this track and return the updated track.
`
}

func (t *track) idString() (id string) {
	switch ti := t.id.(type) {
	case slip.String:
		id = string(ti)
	case slip.Symbol:
		id = string(ti)
	default:
		id = slip.ObjectString(ti)
	}
	return
}

// Scan adds an event to the track.
func (t *track) Scan(flow, task string) {
	t.history = append(t.history, &event{flow: flow, task: task, when: time.Now().UTC()})
}

func (t *track) historyList() (history slip.List) {
	if 0 < len(t.history) {
		history = make(slip.List, len(t.history))
		for i, ev := range t.history {
			history[i] = slip.List{slip.Time(ev.when), slip.String(ev.task), slip.String(ev.flow)}
		}
	}
	return
}

func (t *track) merge(t2 *track) {
	m := map[string]*event{}
	for _, ev := range t.history {
		key := fmt.Sprintf("%s:%s:%d", ev.flow, ev.task, ev.when.UnixNano())
		m[key] = ev
	}
	for _, ev := range t2.history {
		key := fmt.Sprintf("%s:%s:%d", ev.flow, ev.task, ev.when.UnixNano())
		m[key] = ev
	}
	history := make([]*event, 0, len(m))
	for _, ev := range m {
		history = append(history, ev)
	}
	sort.Slice(history, func(i, j int) bool {
		return history[i].when.Before(history[j].when)
	})
	t.history = history
}

// MakeTrack is only public for testing purposes.
func MakeTrack(id slip.Object) (self *flavors.Instance, t *track) {
	self = trackFlavor.MakeInstance().(*flavors.Instance)
	t = &track{id: id}
	self.Any = t

	return
}
