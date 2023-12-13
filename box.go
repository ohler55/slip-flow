// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main

import (
	"io"
	"strings"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/bag"
	"github.com/ohler55/slip/pkg/cl"
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
				slip.String(`A container for data passed between instances of the
_flow-task-flavor_ in a flow. The content of the box can be frozen which forces a
if an attempt is made to modify the content. Typically when transitioning from one
task to another a shallow copy of the box is made and the new box as well as the
original is frozen so that the original box content will not be modified by
modifications to the new box.
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
	return `__:read__ _stream_ &optional _path_ => _self_
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

type boxRemoveCaller struct{}

func (caller boxRemoveCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	if len(args) == 1 {
		removeBox(obj, args[0])
	} else {
		flavors.PanicMethodArgChoice(obj, ":remove", len(args), "1")
	}
	return obj
}

func (caller boxRemoveCaller) Docs() string {
	return `__:remove__ _path_ => _object_
  _path_ to the location in the box to remove. The path must follow the JSONPath format.


Returns the object itself.
`
}

type boxModifyCaller struct{}

func (caller boxModifyCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	modifyBox(s, obj, args, depth+1)
	return obj
}

func (caller boxModifyCaller) Docs() string {
	return `__:modify__ _function_ &optional _path_ &key _:as-bag_ => _object_
  _function_ to modify the value at _path_
  _path_ to the location in the box to modify. The path must follow the JSONPath format.
  _:as-bag_ if true the _function_ expects a _bag_ otherwise it expects a lisp object.


Returns the object itself.
`
}

type boxNativeCaller struct{}

func (caller boxNativeCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		flavors.PanicMethodArgChoice(obj, ":native", len(args), "0")
	}
	return slip.SimpleObject(obj.Any)
}

func (caller boxNativeCaller) Docs() string {
	return `__:native__ => _object_

Returns the box contents as a native LISP form.
`
}

type boxWriteCaller struct{}

func (caller boxWriteCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	return writeBox(obj, args)
}

func (caller boxWriteCaller) Docs() string {
	return `__:write__ &optional _stream_ &key _pretty_ _depth_ _right-margin_ _time-format_ _time-wrap_ _json_ _color_
  _stream_ an output-stream. Default: nil (return a string).
  _:pretty_ The value to use in place of the _*print-pretty*_ value. If _t_ then the
JSON or SEN output is indented according to the other keyword options.
  _:depth_ The maximum number of nested elements on a line in the output. A value
of zero outputs a tight single line output. Default: 4.
  _:right-margin_ The value to use in place of the _*print-right-margin*_ value.
  _:time-format_ The value to use in place of the _*flow-box-time-format*_ value.
  _:time-wrap_ The value to use in place of the _*flow-box-time-wrap*_ value.
  _:json_ If true the output is JSON formatted otherwise output is SEN format.
  _:color_ If true the output is colorized.


Writes the instance to _*standard-output*_, a provided output _stream_,
or to a string that is returned. If _stream_ is _t_ then output is to _*standard-output*_. If
_stream_ is _nil_ then output is a returned string. Any other _stream_ value must be an output
stream which is where output is written to. Output can be either JSON or SEN format as defined
in the OjG package.
`
}

type boxWalkCaller struct{}

func (caller boxWalkCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	walkBox(s, obj, args, depth)
	return nil
}

func (caller boxWalkCaller) Docs() string {
	return `__:walk__ _function_ &optional _path_ _as-lisp_
  _function_ The function to apply to each node in the instance matching the _path_.
  _path_ The path to the location in the box to walk. The path must follow the JSONPath format. Default: ".."
  _as-lisp_ If not nil then the value to the _function_ is a LISP value otherwise a new _bag_.


Walks the values at the location described by _path_.
`
}

type boxBagCaller struct{}

func (caller boxBagCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	bg := bag.Flavor().MakeInstance().(*flavors.Instance)
	bx := obj.Any.(*box)
	if bx.frozen {
		bx.content = alt.Dup(bx.content)
		bx.frozen = false
	}
	bg.Any = bx.content

	return bg
}

func (caller boxBagCaller) Docs() string {
	return `__:bag__ => _instance_


Returns an instance of the _bag-flavor_ with the contents of the box.
`
}

type boxFreezeCaller struct{}

func (caller boxFreezeCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any.(*box).frozen = true

	return nil
}

func (caller boxFreezeCaller) Docs() string {
	return `__:freeze__


Makes the box immutable.
`
}

type boxThawCaller struct{}

func (caller boxThawCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any.(*box).frozen = false

	return nil
}

func (caller boxThawCaller) Docs() string {
	return `__:thaw__


Makes the box mutable.
`
}

type boxFrozenCaller struct{}

func (caller boxFrozenCaller) Call(s *slip.Scope, args slip.List, depth int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)
	if obj.Any.(*box).frozen {
		value = slip.True
	}
	return
}

func (caller boxFrozenCaller) Docs() string {
	return `__:frozen__ => _boolean_


Returns true if the box is immutable.
`
}

type boxTrackingIDCaller struct{}

func (caller boxTrackingIDCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)

	return obj.Any.(*box).track.id
}

func (caller boxTrackingIDCaller) Docs() string {
	return `__:tracking-id__ => _string_|_fixnum_


Returns the tracking identifier of the box.
`
}

type boxTrackCaller struct{}

func (caller boxTrackCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	trk := trackFlavor.MakeInstance().(*flavors.Instance)
	trk.Any = &(obj.Any.(*box).track)

	return trk
}

func (caller boxTrackCaller) Docs() string {
	return `__:track__ => _instance_


Returns the track instance of the box.
`
}

type boxHistoryCaller struct{}

func (caller boxHistoryCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)

	return obj.Any.(*box).track.historyList()
}

func (caller boxHistoryCaller) Docs() string {
	return `__:history__ => _list_


Returns the history of the box track as a list of triples where each triple is a list of
the time, the task name, and the flow name.
`
}

type boxScanCaller struct{}

func (caller boxScanCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	var (
		flowName string
		taskName string
	)
	if ss, ok := args[0].(slip.String); ok {
		flowName = string(ss)
	} else {
		slip.PanicType("flow-name", args[0], "string")
	}
	if ss, ok := args[1].(slip.String); ok {
		taskName = string(ss)
	} else {
		slip.PanicType("task-name", args[1], "string")
	}
	obj.Any.(*box).track.Scan(flowName, taskName)

	return nil
}

func (caller boxScanCaller) Docs() string {
	return `__:scan__ _flow-name_ _task-name_


Add a scan consisting of the current time, task name, and flow name.
`
}

type boxCopyCaller struct{}

func (caller boxCopyCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	orig := obj.Any.(*box)
	inst := boxFlavor.MakeInstance().(*flavors.Instance)
	bx := &box{track: track{id: orig.track.id}, frozen: true}
	orig.frozen = true
	bx.track.history = make([]*event, len(orig.track.history))
	for i, ev := range orig.track.history {
		nev := *ev
		bx.track.history[i] = &nev
	}
	inst.Any = bx

	return inst
}

func (caller boxCopyCaller) Docs() string {
	return `__:copy__ => _box_


Makes the copy of the box with shared content. Both the box and the copy are frozen.
`
}

//////////////////

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
	bx := obj.Any.(*box)
	if bx.frozen {
		bx.content = alt.Dup(bx.content)
		bx.frozen = false
	}
	if x == nil {
		bx.content = v
	} else {
		x.MustSet(bx.content, v)
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
	bx := obj.Any.(*box)
	if bx.frozen {
		bx.content = alt.Dup(bx.content)
		bx.frozen = false
	}
	if x == nil {
		bx.content = v
	} else {
		x.MustSet(bx.content, v)
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
	bx := obj.Any.(*box)
	var value any
	if x == nil {
		value = bx.content
	} else {
		value = x.First(bx.content)
	}
	if value == nil {
		return nil
	}
	if asBag {
		obj = bag.Flavor().MakeInstance().(*flavors.Instance)
		if bx.frozen {
			value = alt.Dup(value)
		}
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

func removeBox(obj *flavors.Instance, path slip.Object) {
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
	bx := obj.Any.(*box)
	if x == nil {
		if bx.frozen {
			bx.frozen = false
		}
		bx.content = nil
	} else {
		if bx.frozen {
			bx.content = alt.Dup(bx.content)
			bx.frozen = false
		}
		bx.content = x.MustRemove(bx.content)
	}
}

func modifyBox(s *slip.Scope, obj *flavors.Instance, args slip.List, depth int) {
	caller := cl.ResolveToCaller(s, args[0], depth)
	var (
		x     jp.Expr
		asBag bool
	)
	if 1 < len(args) {
		switch p := args[1].(type) {
		case nil:
		case slip.String:
			x = jp.MustParseString(string(p))
		case bag.Path:
			x = jp.Expr(p)
		default:
			slip.PanicType("path", p, "string")
		}
		if 2 < len(args) {
			for pos := 2; pos < len(args); pos += 2 {
				sym, ok := args[pos].(slip.Symbol)
				if !ok {
					slip.PanicType("keyword", args[pos], "keyword")
				}
				if len(args)-1 <= pos {
					slip.NewPanic("keyword %s is missing a value", sym)
				}
				if strings.EqualFold(string(sym), ":as-bag") {
					asBag = args[pos+1] != nil
				} else {
					slip.PanicType("keyword", sym, ":as-bag")
				}
			}
		}
	}
	bx := obj.Any.(*box)
	if bx.frozen {
		bx.content = alt.Dup(bx.content)
		bx.frozen = false
	}
	if x == nil {
		bx.content = modifyValue(s, bx.content, caller, asBag, depth)
	} else {
		obj.Any = x.MustModify(bx.content, func(element any) (altered any, changed bool) {
			return modifyValue(s, element, caller, asBag, depth), true
		})
	}
}

func modifyValue(s *slip.Scope, value any, caller slip.Caller, asBag bool, depth int) any {
	bagFlavor := bag.Flavor()
	if asBag {
		bg := bagFlavor.MakeInstance().(*flavors.Instance)
		bg.Any = value
		obj := caller.Call(s, slip.List{bg}, depth)
		if bg, _ := obj.(*flavors.Instance); bg != nil && bg.Flavor == bagFlavor {
			return bg.Any
		}
		return slip.Simplify(obj)
	}
	obj := slip.SimpleObject(value)
	obj = caller.Call(s, slip.List{obj}, depth)
	if bg, _ := obj.(*flavors.Instance); bg != nil && bg.Flavor == bagFlavor {
		return bg.Any
	}
	return slip.Simplify(obj)
}

func writeBox(obj *flavors.Instance, args slip.List) (result slip.Object) {
	var out io.Writer
	dp := slip.DefaultPrinter()
	pw := pretty.Writer{
		Options:  options,
		Width:    int(dp.RightMargin),
		MaxDepth: 4,
		SEN:      true,
	}
	pw.Indent = 2
	prty := dp.Pretty
	if 0 < len(args) {
		pos := 0
		switch ta := args[0].(type) {
		case nil:
			// leave as nil for output to string
		case io.Writer:
			out = ta
			pos++
		case slip.Symbol:
			// probably a key or an error
		default:
			if ta == slip.True {
				out = slip.StandardOutput.(io.Writer)
			} else {
				slip.PanicType("stream", ta, "nil", "t", "output-stream")
			}
			pos++
		}
		for ; pos < len(args)-1; pos += 2 {
			sym := args[pos].(slip.Symbol)
			switch string(sym) {
			case ":pretty":
				prty = args[pos+1] != nil
			case ":depth":
				num, ok := args[pos+1].(slip.Fixnum)
				if !ok {
					slip.PanicType(":depth", args[pos+1], "fixnum")
				}
				pw.MaxDepth = int(num)
				if pw.MaxDepth <= 0 {
					pw.Indent = 0
				}
			case ":right-margin":
				num, ok := args[pos+1].(slip.Fixnum)
				if !ok {
					slip.PanicType(":right-margin", args[pos+1], "fixnum")
				}
				pw.Width = int(num)
			case ":time-format":
				switch ta := args[pos+1].(type) {
				case nil:
					pw.TimeFormat = ""
				case slip.String:
					pw.TimeFormat = string(ta)
				default:
					slip.PanicType(":time-format", args[pos+1], "string")
				}
			case ":time-wrap":
				switch ta := args[pos+1].(type) {
				case nil:
					pw.TimeWrap = ""
				case slip.String:
					pw.TimeWrap = string(ta)
				default:
					slip.PanicType(":time-wrap", args[pos+1], "string")
				}
			case ":json":
				pw.SEN = args[pos+1] == nil
			case ":color":
				pw.Color = args[pos+1] != nil

			default:
				slip.PanicType("keyword", sym, ":pretty", ":depth", ":right-margin",
					":time-format", ":time-wrap", ":json", ":color")
			}
		}
	}
	var b []byte
	switch {
	case prty && 1 < pw.MaxDepth:
		b = pw.Encode(obj.Any)
	case pw.SEN:
		b = sen.Bytes(obj.Any, &pw.Options)
	default:
		b = []byte(oj.JSON(obj.Any, &pw.Options))
	}
	if out == nil {
		return slip.String(b)
	}
	if _, err := out.Write(b); err != nil {
		panic(err)
	}
	return nil
}

func walkBox(s *slip.Scope, obj *flavors.Instance, args slip.List, depth int) {
	fn := args[0]
	path := jp.D()
	var asBag bool
	if 1 < len(args) {
		switch p := args[1].(type) {
		case nil:
		case slip.String:
			path = jp.MustParseString(string(p))
		case bag.Path:
			path = jp.Expr(p)
		default:
			slip.PanicType("path", p, "string", "bag-path")
		}
		asBag = 2 < len(args) && args[2] != nil
	}
	d2 := depth + 1
CallFunc:
	switch tf := fn.(type) {
	case *slip.Lambda:
		if asBag {
			for _, v := range path.Get(obj.Any) {
				arg := bag.Flavor().MakeInstance().(*flavors.Instance)
				arg.Any = v
				_ = tf.Call(s, slip.List{arg}, d2)
			}
		} else {
			for _, v := range path.Get(obj.Any) {
				arg := slip.SimpleObject(v)
				_ = tf.Call(s, slip.List{arg}, d2)
			}
		}
	case *slip.FuncInfo:
		if asBag {
			for _, v := range path.Get(obj.Any) {
				arg := bag.Flavor().MakeInstance().(*flavors.Instance)
				arg.Any = v
				_ = tf.Apply(s, slip.List{arg}, d2)
			}
		} else {
			for _, v := range path.Get(obj.Any) {
				arg := slip.SimpleObject(v)
				tf.Apply(s, slip.List{arg}, d2)
			}
		}
	case slip.Symbol:
		fn = slip.FindFunc(string(tf))
		goto CallFunc
	case slip.List:
		fn = s.Eval(tf, d2)
		goto CallFunc
	default:
		slip.PanicType("function", tf, "function")
	}
}
