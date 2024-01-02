// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main

import (
	"io"

	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/bag"
	"github.com/ohler55/slip/pkg/flavors"
	"github.com/ohler55/slip/pkg/gi"
)

var (
	boxFlavor *flavors.Flavor
)

func init() {
	boxFlavor = flavors.DefFlavor("flow-box-flavor",
		map[string]slip.Object{},
		nil,
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`A container for data passed between instances of the _flow-task-flavor_ in a
flow. The content of the box can be frozen which forces a if an attempt is
made to modify the content. Typically when transitioning from one task to
another a shallow copy of the box is made and the new box as well as the
original is frozen so that the original box content will not be modified by
modifications to the new box.


Like the _gi:bag-flavor_ the data in a _box_ is a tree composed of primitives,
_lists_, and _hash-tables_ like collections. The primitives type are:
  - _t_
  - _nil_
  - _:false_ which maps to a boolean false for a JSON mapping.
  - _integer_
  - _double-float_
  - _string_
  - _time_


Data in a _box_ can be accessed and modified using methods that use a JSONPath
to identify one or more values.

`),
			},
			slip.List{
				slip.Symbol(":init-keywords"),
				slip.Symbol(":tracking-id"),
				slip.Symbol(":track"),
				slip.Symbol(":set"),
				slip.Symbol(":parse"),
				slip.Symbol(":read"),
			},
		},
	)
	boxFlavor.Final = true
	boxFlavor.GoMakeOnly = true
	boxFlavor.DefMethod(":init", "", boxInitCaller{})
	boxFlavor.DefMethod(":set", "", boxSetCaller{})
	boxFlavor.DefMethod(":parse", "", boxParseCaller{})
	boxFlavor.DefMethod(":read", "", boxReadCaller{})
	boxFlavor.DefMethod(":get", "", boxGetCaller{})
	boxFlavor.DefMethod(":has", "", boxHasCaller{})
	boxFlavor.DefMethod(":remove", "", boxRemoveCaller{})
	boxFlavor.DefMethod(":modify", "", boxModifyCaller{})
	boxFlavor.DefMethod(":native", "", boxNativeCaller{})
	boxFlavor.DefMethod(":write", "", boxWriteCaller{})
	boxFlavor.DefMethod(":walk", "", boxWalkCaller{})
	boxFlavor.DefMethod(":bag", "", boxBagCaller{})
	boxFlavor.DefMethod(":freeze", "", boxFreezeCaller{})
	boxFlavor.DefMethod(":thaw", "", boxThawCaller{})
	boxFlavor.DefMethod(":frozen", "", boxFrozenCaller{})
	boxFlavor.DefMethod(":tracking-id", "", boxTrackingIDCaller{})
	boxFlavor.DefMethod(":track", "", boxTrackCaller{})
	boxFlavor.DefMethod(":history", "", boxHistoryCaller{})
	boxFlavor.DefMethod(":scan", "", boxScanCaller{})
	boxFlavor.DefMethod(":copy", "", boxCopyCaller{})
}

type box struct {
	track   track
	content any
	frozen  bool
}

type boxInitCaller struct{}

func (caller boxInitCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		args = args[0].(slip.List)
	}
	var bx box
	for i := 0; i < len(args)-1; i += 2 {
		switch args[i] {
		case slip.Symbol(":tracking-id"):
			bx.track.id = args[i+1]
		case slip.Symbol(":track"):
			if inst, ok := args[i+1].(*flavors.Instance); ok && inst.Flavor == trackFlavor {
				bx.track = *inst.Any.(*track)
			} else {
				slip.PanicType("box :init :track", args[i+1], "flow-track-flavor instance")
			}
		case slip.Symbol(":set"):
			if inst, ok := args[i+1].(*flavors.Instance); ok {
				if inst.Flavor != bag.Flavor() {
					slip.PanicType("box :init :set", args[i+1], "bag-flavor instance")
				}
				bx.content = inst.Any
			} else {
				bx.content = bag.ObjectToBag(args[i+1])
			}
		case slip.Symbol(":parse"):
			so, ok := args[i+1].(slip.String)
			if !ok {
				slip.PanicType("box :init :parse", args[i+1], "string")
			}
			bx.content = sen.MustParse([]byte(so))
			if options.Converter != nil {
				bx.content = options.Converter.Convert(bx.content)
			}
		case slip.Symbol(":read"):
			r, ok := args[i+1].(io.Reader)
			if !ok {
				slip.PanicType("box :init :read", args[i+1], "input-stream")
			}
			bx.content = sen.MustParseReader(r)
			if options.Converter != nil {
				bx.content = options.Converter.Convert(bx.content)
			}
		default:
			slip.PanicType("box :init", args[i], ":tracking-id", ":track", ":set")
		}
	}
	if bx.track.id == nil {
		bx.track.id = gi.NewUUID()
	}
	obj.Any = &bx
	return nil
}

func (caller boxInitCaller) Docs() string {
	return `__:init__ &key _set_ _tracking-id_ _track_ _parse_ _read_
   _:tracking-id_ sets the tracking id of the box to the provided value which can be a string, fixnum, or gi:uuid.
   _:track_ sets the tracking id and events of the box to the provided values.
   _:set_ the contents with the LISP or _bag-flavor_ instance.
   _:parse_ a JSON or SEN string to form the content of the box.
   _:read_ from an _input-stream_ and parses read JSON or SEN to form the content.


Sets the initial value when _make-instance_ is called.


See also: __make-flow-box__
`
}

// MakeBox is only public for testing purposes.
func MakeBox(id slip.Object) (self *flavors.Instance, bx *box) {
	self = boxFlavor.MakeInstance().(*flavors.Instance)
	bx = &box{track: track{id: id}}
	self.Any = bx

	return
}

func boxDup(bi *flavors.Instance) (self *flavors.Instance, bx *box) {
	b := bi.Any.(*box)
	self = boxFlavor.MakeInstance().(*flavors.Instance)

	bx = &box{
		track:   track{id: b.track.id, history: make([]*event, len(b.track.history))},
		content: b.content,
		frozen:  true,
	}
	for i, ev := range b.track.history {
		bx.track.history[i] = &event{when: ev.when, flow: ev.flow, task: ev.task}
	}
	self.Any = bx

	return
}
