// Copyright (c) 2024, Peter Ohler, All rights reserved.

package main

import (
	"fmt"
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
				slip.Symbol(":links"),
			},
		},
	)
	mergeActorFlavor.DefMethod(":init", "", mergeInitCaller{})
	mergeActorFlavor.DefMethod(":start", "", mergeActorStartCaller{})
	mergeActorFlavor.DefMethod(":perform", "", mergeActorPerformCaller{})
	mergeActorFlavor.DefMethod(":shutdown", "", mergeActorShutdownCaller{})
}

// TBD struct for count and current box

type mergeCtx struct {
	task    *task
	timeout time.Duration
	number  int
	//  mutex, map, timeout loop
}

type mergeInitCaller struct{}

func (caller mergeInitCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	self := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		args = args[0].(slip.List)
	}
	var mc mergeCtx
	for pos := 0; pos < len(args)-1; pos += 2 {
		sym, _ := args[pos].(slip.Symbol)
		switch string(sym) {
		case ":timeout":
			if num, ok := args[pos+1].(slip.Fixnum); ok {
				mc.timeout = time.Duration(num) * time.Second
			} else {
				slip.PanicType(":timeout", args[pos+1], "fixnum")
			}
		case ":number":
			if num, ok := args[pos+1].(slip.Fixnum); ok && 0 < num {
				mc.number = int(num)
			} else {
				slip.PanicType(":number", args[pos+1], "fixnum")
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
	obj.Any.(*mergeCtx).task = args[0].(*flavors.Instance).Any.(*task)

	// TBD start by creating map and starting timeout loop
	//  use select for exit channel and timer

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

	// TBD

	fmt.Printf("*** mc: %v %s\n", mc, bi)

	return slip.List{nil, nil}
}

func (caller mergeActorPerformCaller) Docs() string {
	return `__:perform__ _box_
   _:box_ [instance] TBD


TBD
`
}

type mergeActorShutdownCaller struct{}

func (caller mergeActorShutdownCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	mc := obj.Any.(*mergeCtx)

	fmt.Printf("*** mc: %v\n", mc)

	return nil
}

func (caller mergeActorShutdownCaller) Docs() string {
	return `__:shutdown__


Shutsdown the actor by exiting the timeout checking loop.
`
}
