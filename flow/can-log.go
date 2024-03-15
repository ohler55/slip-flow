// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

const (
	levelError int = iota
	levelWarn
	levelInfo
	levelDebug
)

var (
	canLogFlavor *flavors.Flavor
)

func init() {
	Pkg.Initialize(nil)
	canLogFlavor = flavors.DefFlavor("can-log",
		map[string]slip.Object{ // instance variables
			"log-level": slip.Fixnum(1),
			"logger":    nil, // boolean
		},
		nil,
		slip.List{
			slip.Symbol(":gettable-instance-variables"),
			slip.Symbol(":inittable-instance-variables"),
			slip.Symbol(":abstract-flavor"),
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`An abstract mixin that adds support for logging.`),
			},
		},
		&Pkg,
	)
	canLogFlavor.DefMethod(":error", "", canLogErrorCaller{})
	canLogFlavor.DefMethod(":warn", "", canLogWarnCaller{})
	canLogFlavor.DefMethod(":info", "", canLogInfoCaller{})
	canLogFlavor.DefMethod(":debug", "", canLogDebugCaller{})
	canLogFlavor.DefMethod(":set-level", "", canLogSetLogLevelCaller{})
}

type canLogErrorCaller struct{}

func (caller canLogErrorCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	if logger := s.Get("logger").(slip.Instance); logger != nil {
		_ = logger.Receive(s, ":log", append(slip.List{slip.Symbol(":error")}, args...), depth)
	}
	return nil
}

func (caller canLogErrorCaller) Docs() string {
	return `__:error__ _format_ &rest _args_


Log a error message.
`
}

type canLogWarnCaller struct{}

func (caller canLogWarnCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	level := s.Get("log-level").(slip.Fixnum)
	if 1 <= level {
		if logger := s.Get("logger").(slip.Instance); logger != nil {
			_ = logger.Receive(s, ":log", append(slip.List{slip.Symbol(":warn")}, args...), depth)
		}
	}
	return nil
}

func (caller canLogWarnCaller) Docs() string {
	return `__:warn__ _format_ &rest _args_


Log a warn message if the _log-level_ is at or above 1.
`
}

type canLogInfoCaller struct{}

func (caller canLogInfoCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	level := s.Get("log-level").(slip.Fixnum)
	if 1 <= level {
		if logger := s.Get("logger").(slip.Instance); logger != nil {
			_ = logger.Receive(s, ":log", append(slip.List{slip.Symbol(":info")}, args...), depth)
		}
	}
	return nil
}

func (caller canLogInfoCaller) Docs() string {
	return `__:info__ _format_ &rest _args_


Log a info message if the _log-level_ is at or above 2.
`
}

type canLogDebugCaller struct{}

func (caller canLogDebugCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	level := s.Get("log-level").(slip.Fixnum)
	if 3 <= level {
		if logger := s.Get("logger").(slip.Instance); logger != nil {
			_ = logger.Receive(s, ":log", append(slip.List{slip.Symbol(":debug")}, args...), depth)
		}
	}
	return nil
}

func (caller canLogDebugCaller) Docs() string {
	return `__:debug__ _format_ &rest _args_


Log a debug message if the _log-level_ is at or above 3.
`
}

type canLogSetLogLevelCaller struct{}

func (caller canLogSetLogLevelCaller) Call(s *slip.Scope, args slip.List, depth int) slip.Object {
	level := -1
	switch ta := args[0].(type) {
	case slip.Fixnum:
		level = int(ta)
		if level < levelError {
			level = levelError
		} else if levelDebug < level {
			level = levelDebug
		}
	case slip.Symbol:
		switch ta {
		case slip.Symbol("error"), slip.Symbol(":error"):
			level = levelError
		case slip.Symbol("warn"), slip.Symbol(":warn"):
			level = levelWarn
		case slip.Symbol("info"), slip.Symbol(":info"):
			level = levelInfo
		case slip.Symbol("debug"), slip.Symbol(":debug"):
			level = levelDebug
		}
	}
	if level < 0 {
		slip.PanicType("level", args[0], "0", "1", "2", "3'", ":error", ":warn", ":info", "debug")
	}
	s.Set("log-level", slip.Fixnum(level))

	return slip.Fixnum(level)
}

func (caller canLogSetLogLevelCaller) Docs() string {
	return `__:set-log-level__ _level_ => _fixnum_
   _level_ [fixnum] to set the _log-level_ to. Can be a fixnum between 0 and 3
inclusive or :error, :warn, :info, or :debug.


Set the log-level and returns the _log-level_ as a _fixnum_.
`
}
