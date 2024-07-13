// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/ohler55/slip"
	"github.com/ohler55/slip/pkg/flavors"
)

var (
	writeFileActorFlavor *flavors.Flavor
)

func defWriteFileActor() {
	writeFileActorFlavor = flavors.DefFlavor("flow-write-file-actor",
		map[string]slip.Object{},
		nil,
		slip.List{
			slip.List{
				slip.Symbol(":documentation"),
				slip.String(`A flow-write-file-actor writes a file before forwarding the received _box_ on
the "ok" link.`),
			},
			slip.List{
				slip.Symbol(":init-keywords"),
				slip.Symbol(":filename"),
				slip.Symbol(":content"),
				slip.Symbol(":overwrite"),
				slip.Symbol(":append"),
				slip.Symbol(":permissions"),
			},
		},
		&Pkg,
	)
	writeFileActorFlavor.DefMethod(":init", "", writeFileInitCaller{})
	writeFileActorFlavor.DefMethod(":start", "", writeFileActorStartCaller{})
	writeFileActorFlavor.DefMethod(":perform", "", writeFileActorPerformCaller{})
	writeFileActorFlavor.DefMethod(":init-key-values", "", writeFileActorInitKeyValuesCaller{})
}

type writeFileCtx struct {
	task     *task
	filename strCaller
	content  strCaller
	flag     int
	perm     fs.FileMode
}

type writeFileInitCaller struct{}

func (caller writeFileInitCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	self := s.Get("self").(*flavors.Instance)
	if 0 < len(args) {
		args = args[0].(slip.List)
	}
	wfc := writeFileCtx{
		perm: 0664,
		flag: os.O_CREATE | os.O_WRONLY,
	}
	for pos := 0; pos < len(args)-1; pos += 2 {
		sym, _ := args[pos].(slip.Symbol)
		switch string(sym) {
		case ":filename":
			wfc.filename.extract(s, args[pos+1])
		case ":content":
			wfc.content.extract(s, args[pos+1])
		case ":overwrite":
			if args[pos+1] != nil {
				wfc.flag |= os.O_TRUNC
			}
		case ":append":
			if args[pos+1] != nil {
				wfc.flag |= os.O_APPEND
			}
		case ":permissions":
			switch ta := args[pos+1].(type) {
			case slip.String:
				wfc.perm = parsePerm(string(ta))
			case slip.Symbol:
				wfc.perm = parsePerm(string(ta))
			case slip.Fixnum:
				wfc.perm = fs.FileMode(ta)
			default:
				slip.PanicType(":permissions", args[pos+1], "string", "symbol", "fixnum")
			}
		}
	}
	self.Any = &wfc

	return nil
}

func (caller writeFileInitCaller) Docs() string {
	return `__:init__ &key _filename_ _content_ _overwrite_ _append_ _permissions_
   _:filename_ [string|symbol|function] of the file to write.
   _:content_ [string|function] the content to write.
   _:overwrite_ [boolean] if true the file is replaced or overwritten.
   _:append_ [boolean] if true the content is appended to the file.
   _:permissions_ [fixnum|string] if the file is to be created then this is used as the permission.
The format is the unix permission format such as "-rw-rw-r--" or number.


Sets the initial value when _make-instance_ is called.
`
}

type writeFileActorStartCaller struct{}

func (caller writeFileActorStartCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	obj.Any.(*writeFileCtx).task = args[0].(*flavors.Instance).Any.(*task)

	return nil
}

func (caller writeFileActorStartCaller) Docs() string {
	return `__:start__ _task_
   _:task_ [instance] the task that contains the actor.


Sets the context for the actor.
`
}

type writeFileActorPerformCaller struct{}

func (caller writeFileActorPerformCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	wfc := obj.Any.(*writeFileCtx)
	bi := args[0].(*flavors.Instance)

	f, err := os.OpenFile(wfc.filename.value(s, bi), wfc.flag, wfc.perm)
	if err != nil {
		panic(err)
	}
	defer func() { _ = f.Close() }()
	if _, err = f.WriteString(wfc.content.value(s, bi)); err != nil {
		panic(err)
	}
	return slip.List{slip.String("ok"), bi}
}

func (caller writeFileActorPerformCaller) Docs() string {
	return `__:perform__ _box_
   _:box_ [instance] the data to submit to the target flow.


Submits a box to the target flow.
`
}

type writeFileActorInitKeyValuesCaller struct{}

func (caller writeFileActorInitKeyValuesCaller) Call(s *slip.Scope, args slip.List, _ int) slip.Object {
	obj := s.Get("self").(*flavors.Instance)
	wfc := obj.Any.(*writeFileCtx)
	kvs := slip.List{
		slip.Symbol(":filename"), wfc.filename.raw(),
		slip.Symbol(":content"), wfc.content.raw(),
		slip.Symbol(":permissions"), slip.String(wfc.perm.String()),
	}
	if wfc.flag&os.O_TRUNC != 0 {
		kvs = append(kvs, slip.Symbol(":overwrite"), slip.True)
	} else {
		kvs = append(kvs, slip.Symbol(":overwrite"), nil)
	}
	if wfc.flag&os.O_APPEND != 0 {
		kvs = append(kvs, slip.Symbol(":append"), slip.True)
	} else {
		kvs = append(kvs, slip.Symbol(":append"), nil)
	}
	return kvs
}

func (caller writeFileActorInitKeyValuesCaller) Docs() string {
	return `__:init-key-values__ => (:filename "write-me.txt")


Returns the keywords and values needed to recreate the instance as a property list.
`
}

const permAllow = "drwxrwxrwx"

func parsePerm(s string) (perm fs.FileMode) {
	off := 0
	switch len(s) {
	case 9:
		off = 1
	case 10:
	default:
		panic(fmt.Sprintf("%s is not a valid symbolic permission", s))
	}
	for i, c := range []byte(s) {
		if permAllow[i+off] == c {
			perm |= 1 << (9 - i - off)
		} else if c != '-' {
			panic(fmt.Sprintf("%s is not a valid symbolic permission", s))
		}
	}
	return
}
