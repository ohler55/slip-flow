// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

func init() {
	slip.Define(
		func(args slip.List) slip.Object {
			f := TaskTransition{Function: slip.Function{Name: "flow-task-transition", Args: args}}
			f.Self = &f
			return &f
		},
		&slip.FuncDoc{
			Name: "flow-task-transition",
			Args: []*slip.DocArg{
				{
					Name: "task",
					Type: "flow-task",
					Text: "to transition a box from.",
				},
				{
					Name: "box",
					Type: "instance",
					Text: "to transition.",
				},
				{
					Name: "link",
					Type: "string",
					Text: "name of link to transition on.",
				},
			},
			Return: "nil",
			Text:   `__flow-task-transition__ sends a _box_ to on _link_.`,
			Examples: []string{
				`(setq flow (make-instance 'flow :link "flo"))`,
				`(setq tick (flow-add-task flow :name 'tick :actor (lambda (b) (list 'ok b))))`,
				`(setq tock (flow-add-task flow :name 'tock :actor (lambda (b) (list 'ok b))))`,
				`(flow-link flow "ok" 'tick 'tock) => nil`,
				`;; The transition method or function is usually called by an actor wrapped by a task.`,
				`(flow-task-transition tick (make-flow-box :parse "[1 2]") 'ok) => nil`,
			},
		}, &Pkg)
}

// TaskTransition represents the flow-task-transition function.
type TaskTransition struct {
	slip.Function
}

// Call the function with the arguments provided.
func (f *TaskTransition) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	slip.CheckArgCount(s, depth, f, args, 3, 3)
	self, ok := args[0].(*flavors.Instance)
	if !ok || self.Type != taskFlavor {
		slip.TypePanic(s, depth, "task", args[0], "task")
	}
	var (
		bi       *flavors.Instance
		linkName string
	)
	if bi, ok = args[1].(*flavors.Instance); !ok || boxFlavor != bi.Type {
		slip.TypePanic(s, depth, "box", args[1], "box")
	}
	switch ta := args[2].(type) {
	case nil:
		// an empty string
	case slip.String:
		linkName = string(ta)
	case slip.Symbol:
		linkName = string(ta)
	default:
		slip.TypePanic(s, depth, "link", ta, "nil", "string", "symbol")
	}
	self.Any.(*task).transition(s, linkName, bi)

	return nil
}

type taskTransitionCaller struct{}

func (caller taskTransitionCaller) Call(s *slip.Scope, args slip.List, depth int) (value slip.Object) {
	obj := s.Get("self").(*flavors.Instance)

	bi, ok := args[0].(*flavors.Instance)
	if !ok || boxFlavor != bi.Type {
		slip.TypePanic(s, depth, "box", args[0], "box")
	}
	var linkName string
	switch ta := args[1].(type) {
	case nil:
		// an empty string
	case slip.String:
		linkName = string(ta)
	case slip.Symbol:
		linkName = string(ta)
	default:
		slip.TypePanic(s, depth, "link", ta, "nil", "string", "symbol")
	}
	obj.Any.(*task).transition(s, linkName, bi)

	return nil
}

func (caller taskTransitionCaller) FuncDocs() *slip.FuncDoc {
	return methodDocsFromFunc(":transition", "flow-task-transition", "flow-task", "task")
}
