// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := BoxNotify{Function: slip.Function{Name: "flow-box-notify", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-box-notify",
			Args: []*slip.DocArg{
				{
					Name: "box",
					Type: "flow-box",
					Text: "to push onto a channel.",
				},
				{Name: "&rest"},
				{
					Name: "names",
					Type: "string",
					Text: `of the watchers to send the box on.`,
				},
			},
			Return: "fixnum",
			Text: `__flow-box-notify__ pushes the box onto the channels named by _names_
if no _names_ then to all. The number of channels pushed to is returned`,
			Examples: []string{
				`(setq box (make-instance 'flow-box-flavor :parse "{a:7}")) => #<flow-box-flavor 12345>`,
				`(setq chan (make-channel 5) => #<channel 12345>`,
				`(flow-box-watch box "done" chan) => nil`,
				`(flow-box-notify box "done") => 1`,
			},
		}, &Pkg)
}

// BoxNotify represents the flow-box-notify function.
type BoxNotify struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *BoxNotify) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	self, ok := args[0].(*flavors.Instance)
	if !ok {
		slip.PanicType("box", args[0], "box")
	}
	_ = self.Receive(s, ":notify", args[1:], depth)

	return nil
}

type boxNotifyCaller struct{}

func (caller boxNotifyCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)

	return notifyBox(obj, args)
}

func (caller boxNotifyCaller) Docs() string {
	return methodDocFromFunc(":notify", "flow-box-notify", "flow-box-flavor", "box")
}

func notifyBox(obj *flavors.Instance, args slip.List) slip.Object {
	bx := obj.Any.(*box)
	var cnt int
	if 0 < len(args) {
		var key string
		for _, arg := range args {
			switch tn := arg.(type) {
			case slip.String:
				key = string(tn)
			case slip.Symbol:
				key = string(tn)
			default:
				slip.PanicType("names", tn, "string", "symbol")
			}
			if c := bx.watchers[key]; c != nil {
				c <- obj
				cnt++
			}
		}
	} else {
		for _, c := range bx.watchers {
			c <- obj
			cnt++
		}
	}
	return slip.Fixnum(cnt)
}
