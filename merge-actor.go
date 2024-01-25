// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

import (
	"sync"
	"time"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

var (
	mergeActorFlavor *flavors.Flavor
)

func init() {
	mergeActorFlavor = flavors.DefFlavor("flow-merge-actor",
		map[string]slip.Object{},
		nil,
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`A flow-merge-actor merged boxes from multiple branches.  A merge is usually
used to merge branches created by a _flow-split-actor_. Both the tracks and
content of the boxes with the same track-id are merged. A track merge forms a
union of all events in the track history keeping only unique events and then
sorts the history by time. Content is merged by setting any element not in the
first box received with values from subsequent boxes. Element are not added to
arrays.
`),
			},
			slip.List{
				slip.Symbol(":init-keywords"),
				slip.Symbol(":number"),
				slip.Symbol(":timeout"),
			},
		},
	)
	mergeActorFlavor.DefMethod(":init", "", mergeInitCaller{})
	mergeActorFlavor.DefMethod(":start", "", mergeActorStartCaller{})
	mergeActorFlavor.DefMethod(":perform", "", mergeActorPerformCaller{})
	mergeActorFlavor.DefMethod(":shutdown", "", mergeActorShutdownCaller{})
}

type boxCnt struct {
	cnt     int
	updated time.Time
	box     *box
}

type mergeCtx struct {
	task    *task
	timeout time.Duration
	number  int
	mu      sync.Mutex
	pending map[string]*boxCnt
	stop    chan struct{}
	done    chan struct{}
}

func (mc *mergeCtx) timeoutLoop(s *slip.Scope) {
	tick := time.NewTicker(mc.timeout / 2)
	defer tick.Stop()
	for {
		select {
		case <-mc.stop:
			mc.done <- struct{}{}
			return
		case <-tick.C:
			mc.mu.Lock()
			for k, bc := range mc.pending {
				if mc.timeout < time.Since(bc.updated) {
					delete(mc.pending, k)
					bi := boxFlavor.MakeInstance().(*flavors.Instance)
					bi.Any = bc
					mc.task.handlePanic(s, bi)
				}
			}
			mc.mu.Unlock()
		}
	}
}

// return a box if completed or nil otherwise
func (mc *mergeCtx) addBox(bx *box) (full *box) {
	id := bx.track.idString()
	mc.mu.Lock()
	bc := mc.pending[id]
	if bc == nil {
		bc = &boxCnt{cnt: 0, updated: time.Now(), box: bx}
		mc.pending[id] = bc
	} else {
		bc.box.merge(bx)
	}
	bc.cnt++
	if mc.number <= bc.cnt {
		delete(mc.pending, id)
		full = bc.box
	}
	mc.mu.Unlock()
	return
}

type mergeInitCaller struct{}

func (caller mergeInitCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	self := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		args = args[0].(slip.List)
	}
	mc := mergeCtx{number: 1, timeout: time.Second * 10}
	for pos := 0; pos < len(args)-1; pos += 2 {
		sym, _ := args[pos].(slip.Symbol)
		switch string(sym) {
		case ":timeout":
			if num, ok := args[pos+1].(slip.Fixnum); ok && 0 < num {
				mc.timeout = time.Duration(num) * time.Second
			} else {
				slip.PanicType(":timeout", args[pos+1], "positive fixnum")
			}
		case ":number":
			if num, ok := args[pos+1].(slip.Fixnum); ok && 0 < num {
				mc.number = int(num)
			} else {
				slip.PanicType(":number", args[pos+1], "positive fixnum")
			}
		default:
			slip.PanicType("keywords", args[pos], ":timeout", ":number")
		}
	}
	self.Any = &mc

	return nil
}

func (caller mergeInitCaller) Docs() string {
	return `__:init__ &key _timeout_ _number_
   _:timeout_ [fixnum] seconds before timing out waiting for _number_ of boxes.
   _:number_ [fixnum] number of boxes expected before transitioning.


Waits for _number_ of boxes with the same track-id before transitioning on the "ok" link.
`
}

type mergeActorStartCaller struct{}

func (caller mergeActorStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	mc := obj.Any.(*mergeCtx)
	mc.task = args[0].(*flavors.Instance).Any.(*task)
	mc.pending = map[string]*boxCnt{}
	mc.stop = make(chan struct{}, 1)
	mc.done = make(chan struct{}, 1)
	go mc.timeoutLoop(s)

	return nil
}

func (caller mergeActorStartCaller) Docs() string {
	return `__:start__ _task_
   _:task_ [instance] the task that contains the actor.


Sets the context for the actor.
`
}

type mergeActorPerformCaller struct{}

func (caller mergeActorPerformCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	bi := args[0].(*flavors.Instance)

	mc := obj.Any.(*mergeCtx)
	if full := mc.addBox(bi.Any.(*box)); full != nil {
		bi = boxFlavor.MakeInstance().(*flavors.Instance)
		bi.Any = full

		return slip.List{slip.Symbol("ok"), bi}
	}
	return slip.List{nil, nil}
}

func (caller mergeActorPerformCaller) Docs() string {
	return `__:perform__ _box_
   _:box_ [instance] box to merge with other boxes received.


Merges all the boxes received with matching track IDs. When the expected
number of boxes is received and merged a transition is made on the "ok"
link. If a timeout occurs it is handled like any other error.
`
}

type mergeActorShutdownCaller struct{}

func (caller mergeActorShutdownCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	mc := obj.Any.(*mergeCtx)

	mc.stop <- struct{}{}
	<-mc.done

	return nil
}

func (caller mergeActorShutdownCaller) Docs() string {
	return `__:shutdown__


Shuts down the actor by exiting the timeout checking loop.
`
}
