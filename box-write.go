// Copyright (c) 2023, Peter Ohler, All rights reserved.

package main

import (
	"io"

	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/pretty"
	"github.com/ohler55/ojg/sen"
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := BoxWrite{Function: slip.Function{Name: "flow-box-write", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-box-write",
			Args: []*slip.DocArg{
				{
					Name: "box",
					Type: "flow-box",
					Text: "to write.",
				},
				{Name: "&optional"},
				{
					Name: "stream",
					Type: "output-stream",
					Text: "stream to write to.",
				},
				{Name: "&key"},
				{
					Name: "pretty",
					Type: "boolean",
					Text: `value to use in place of the _*print-pretty*_ value.
If _t_ then the JSON or SEN output is indented according to the other keyword options.`,
				},
				{
					Name: "depth",
					Type: "fixnum",
					Text: `maximum number of nested elements on a line in the output.
A value of zero outputs a tight single line output. Default: 4.`,
				},
				{
					Name: "right-margin",
					Type: "fixnum",
					Text: "value to use in place of the _*print-right-margin*_ value.",
				},
				{
					Name: "indent",
					Type: "fixnum",
					Text: "is the number of spaces to indent JSON or SEN output if :pretty is not non-nil.",
				},
				{
					Name: "time-format",
					Type: "string",
					Text: "value to use in place of the _*flow-box-time-format*_ value.",
				},
				{
					Name: "time-wrap",
					Type: "string",
					Text: "value to use in place of the _*flow-box-time-wrap*_ value.",
				},
				{
					Name: "json",
					Type: "boolean",
					Text: "if true the output is JSON formatted otherwise output is SEN format.",
				},
				{
					Name: "color",
					Type: "boolean",
					Text: "if true the output is colorized.",
				},
				{
					Name: "full",
					Type: "boolean",
					Text: "if true the output includes the box track and the content is nested on level down.",
				},
			},
			Return: "box",
			Text: `__flow-box-write__ writes the instance to _*standard-output*_, a provided output _stream_,
or to a string that is returned. If _stream_ is _t_ then output is to _*standard-output*_. If
_stream_ is _nil_ then output is a returned string. Any other _stream_ value must be an output
stream which is where output is written to. Output can be either JSON or SEN format as defined
in the OjG package.`,
			Examples: []string{
				`(setq box (make-instance 'flow-box-flavor :parse "{a:[1 2 3]}"))`,
				`(flow-box-write box nil) => "{a: [1 2 3]}"`,
			},
		}, &Pkg)
}

// BoxWrite represents the flow-box-write function.
type BoxWrite struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *BoxWrite) Call(s *slip.Scope, args slip.List, depth int) (result slip.Object) {
	self, ok := args[0].(*flavors.Instance)
	if !ok {
		slip.PanicType("box", args[0], "box")
	}
	return self.Receive(s, ":write", args[1:], depth)
}

type boxWriteCaller struct{}

func (caller boxWriteCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	return writeBox(s, obj, args)
}

func (caller boxWriteCaller) Docs() string {
	return methodDocFromFunc(":write", "flow-box-write", "flow-box-flavor", "box")
}

func writeBox(s *slip.Scope, obj *flavors.Instance, args slip.List) (result slip.Object) {
	var out io.Writer
	dp := slip.DefaultPrinter()
	pw := pretty.Writer{
		Options:  options,
		Width:    int(dp.RightMargin),
		MaxDepth: 4,
		SEN:      true,
	}
	var full bool
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
				out = s.Get("*standard-output*").(io.Writer)
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
			case ":indent":
				num, ok := args[pos+1].(slip.Fixnum)
				if !ok {
					slip.PanicType(":indent", args[pos+1], "fixnum")
				}
				pw.Indent = int(num)
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
			case ":full":
				full = args[pos+1] != nil

			default:
				slip.PanicType("keyword", sym, ":pretty", ":depth", ":right-margin", "indent",
					":time-format", ":time-wrap", ":json", ":color", "full")
			}
		}
	}
	bx := obj.Any.(*box)
	content := bx.content
	if full {
		history := make([]any, len(bx.track.history))
		for i, ev := range bx.track.history {
			history[i] = map[string]any{
				"when": ev.when,
				"flow": ev.flow,
				"task": ev.task,
			}
		}
		content = map[string]any{
			"track": map[string]any{
				"id":      slip.Simplify(bx.track.id),
				"history": history,
			},
			"content": content,
		}
	}
	var b []byte
	switch {
	case prty && 1 < pw.MaxDepth:
		b = pw.Encode(content)
	case pw.SEN:
		b = sen.Bytes(content, &pw.Options)
	default:
		b = []byte(oj.JSON(content, &pw.Options))
	}
	if out == nil {
		return slip.String(b)
	}
	if _, err := out.Write(b); err != nil {
		panic(err)
	}
	return nil
}
