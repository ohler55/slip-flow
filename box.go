// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main

import (
	"io"

	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/bag"
	"github.com/ohler55/slip/pkg/flavors"
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
				slip.String(`A container for data passed between instances of the
_flow-task-flavor_ in a flow.`),
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

	// boxFlavor.DefMethod(":remove", "", boxRemoveCaller{})
	// boxFlavor.DefMethod(":modify", "", boxModifyCaller{})
	// boxFlavor.DefMethod(":native", "", boxNativeCaller{})
	// boxFlavor.DefMethod(":write", "", boxWriteCaller{})
	// boxFlavor.DefMethod(":walk", "", boxWalkCaller{})
	// boxFlavor.DefMethod(":bag", "", boxBagCaller{})
	// boxFlavor.DefMethod(":native", "", boxNativeCaller{})
	// boxFlavor.DefMethod(":tracking-id", "", boxTrackingIDCaller{})
	// boxFlavor.DefMethod(":freeze", "", boxFreezeCaller{})
	// boxFlavor.DefMethod(":track", "", boxTrackCaller{}) // instance
	// boxFlavor.DefMethod(":events", "", boxEventsCaller{})
	// boxFlavor.DefMethod(":scan", "", boxScanCaller{}) // adds and entry to the track
	// TBD
}

type box struct {
	track   track
	content any
	// err     slip.Object
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
				bx.track = *inst.Any.(*track)
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
	obj.Any = &bx
	return nil
}

func (caller boxInitCaller) Docs() string {
	return `__:init__ &key _set_ _tracking-id_ _track_
   _:tracking-id_ sets the tracking id of the box to the provided value.
   _:track_ sets the tracking id and events of the box to the provided values.
   _:set_ the contents with the LISP or _bag-flavor_ instance.
   _:parse_ a JSON or SEN string to form the content of the box.
   _:read_ from an _input-stream_ and parses read JSON or SEN to form the content.


Sets the initial value when _make-instance_ is called.
`
}

type boxSetCaller struct{}

func (caller boxSetCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	switch len(args) {
	case 1:
		setBox(obj, args[0], nil)
	case 2:
		setBox(obj, args[0], args[1])
	default:
		flavors.PanicMethodArgChoice(obj, ":set", len(args), "1 or 2")
	}
	return obj
}

func (caller boxSetCaller) Docs() string {
	return `__:set__ _value_ &optional _path_ => _self_
  _value_ The value to set in the instance according to the path.
  _path_ The path to the location in the box to set the _value_.
The path must follow the JSONPath format.

Sets a _value_ at the location described by _path_.
If no _path_ is provided the entire contents of the box is replaced.
`
}

type boxParseCaller struct{}

func (caller boxParseCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	switch len(args) {
	case 1:
		parseBox(obj, args[0], nil)
	case 2:
		parseBox(obj, args[0], args[1])
	default:
		flavors.PanicMethodArgChoice(obj, ":parse", len(args), "1 or 2")
	}
	return obj
}

func (caller boxParseCaller) Docs() string {
	return `__:parse__ _string_ &optional _path_ => _self_
  _string_ The string to parse and set in the instance according to the _path_.
  _path_ The path to the location in the box to set the parsed value.
The path must follow the JSONPath format.


Parses _string_ and sets the parsed value at the location described by _path_.
If no _path_ is provided the entire contents of the box is replaced.
`
}

type boxReadCaller struct{}

func (caller boxReadCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	switch len(args) {
	case 1:
		readBox(obj, args[0], nil)
	case 2:
		readBox(obj, args[0], args[1])
	default:
		flavors.PanicMethodArgChoice(obj, ":read", len(args), "1 or 2")
	}
	return obj
}

func (caller boxReadCaller) Docs() string {
	return `__:read__ _string_ &optional _path_ => _self_
  _stream_ The _input-stream_ to read and set in the instance according to the _path_.
  _path_ The path to the location in the box to set the readd value.
The path must follow the JSONPath format.


Read from _stream_ and sets the parsed value at the location described by _path_.
If no _path_ is provided the entire contents of the box is replaced.
`
}

type boxGetCaller struct{}

func (caller boxGetCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)
	switch len(args) {
	case 0:
		value = getBox(obj, nil, false)
	case 1:
		value = getBox(obj, args[0], false)
	case 2:
		value = getBox(obj, args[0], args[1] != nil)
	default:
		flavors.PanicMethodArgCount(obj, ":get", len(args), 0, 2)
	}
	return
}

func (caller boxGetCaller) Docs() string {
	return `__:get__ &optional _path_ _as-bag_ => _object_|_bag_
  _path_ to the location in the box to get the _value_ from. The path must follow the JSONPath format.
  _as-bag_ if not nil then the returned value is a _bag_ otherwise a new LISP value is returned.


Gets a _value_ at the location described by _path_.
If no _path_ is provided the entire contents of the box is returned.
.`
}

type boxHasCaller struct{}

func (caller boxHasCaller) Call(s *slip.Scope, args slip.List, _ int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)
	if len(args) == 1 {
		value = hasBox(obj, args[0])
	} else {
		flavors.PanicMethodArgChoice(obj, ":has", len(args), "1")
	}
	return
}

func (caller boxHasCaller) Docs() string {
	return `__:has__ _path_ => _boolean_
  _path_ to the location in the box to get the value from. The path must follow the JSONPath format.


Returns true if a value at the location described by _path_ exists.
`
}

// MakeBox is only public for testing purposes.
func MakeBox(id slip.Object) (self *flavors.Instance, bx *box) {
	self = boxFlavor.MakeInstance().(*flavors.Instance)
	bx = &box{track: track{id: id}}
	self.Any = bx

	return
}

// TBD move the following functions to individual function files (flow-box-set ...)

func parseBox(obj *flavors.Instance, value, path slip.Object) {
	var x jp.Expr
	switch p := path.(type) {
	case nil:
	case slip.String:
		x = jp.MustParseString(string(p))
	case bag.Path:
		x = jp.Expr(p)
	default:
		slip.PanicType("path", p, "string")
	}
	ss, ok := value.(slip.String)
	if !ok {
		slip.PanicType("string", value, "string")
	}
	v := sen.MustParse([]byte(ss))
	if options.Converter != nil {
		v = options.Converter.Convert(v)
	}
	if x == nil {
		obj.Any.(*box).content = v
	} else {
		x.MustSet(obj.Any.(*box).content, v)
	}
}

func readBox(obj *flavors.Instance, value, path slip.Object) {
	var x jp.Expr
	switch p := path.(type) {
	case nil:
	case slip.String:
		x = jp.MustParseString(string(p))
	case bag.Path:
		x = jp.Expr(p)
	default:
		slip.PanicType("path", p, "string")
	}
	r, ok := value.(io.Reader)
	if !ok {
		slip.PanicType("stream", value, "input-stream")
	}
	v := sen.MustParseReader(r)
	if options.Converter != nil {
		v = options.Converter.Convert(v)
	}
	if x == nil {
		obj.Any.(*box).content = v
	} else {
		x.MustSet(obj.Any.(*box).content, v)
	}
}

func getBox(obj *flavors.Instance, path slip.Object, asBag bool) slip.Object {
	var x jp.Expr
	switch p := path.(type) {
	case nil:
	case slip.String:
		x = jp.MustParseString(string(p))
	case bag.Path:
		x = jp.Expr(p)
	default:
		slip.PanicType("path", p, "string", "bag-path")
	}
	var value any
	if x == nil {
		value = obj.Any
	} else {
		value = x.First(obj.Any)
	}
	if value == nil {
		return nil
	}
	if asBag {
		obj = bag.Flavor().MakeInstance().(*flavors.Instance)
		obj.Any = value

		return obj
	}
	return slip.SimpleObject(value)
}

func hasBox(obj *flavors.Instance, path slip.Object) slip.Object {
	var x jp.Expr
	switch p := path.(type) {
	case nil:
	case slip.String:
		x = jp.MustParseString(string(p))
	case bag.Path:
		x = jp.Expr(p)
	default:
		slip.PanicType("path", p, "string")
	}
	if x == nil || x.Has(obj.Any.(*box).content) {
		return slip.True
	}
	return nil
}
