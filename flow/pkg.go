// Copyright (c) 2024, Peter Ohler, All rights reserved.

package flow

import (
	"fmt"
	"strings"
	"time"

	"github.com/ohler55/ojg"
	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/slip"
)

var (
	// Pkg is the message package.
	Pkg = slip.Package{
		Name:      "flow",
		Nicknames: []string{"flow"},
		Doc:       "Home of symbols defined for the flow functions, variables, and constants.",
		PreSet:    slip.DefaultPreSet,
	}

	options = ojg.DefaultOptions
)

func init() {
	Pkg.Initialize(map[string]*slip.VarVal{
		"*flow-box-time-format*": {
			Get:    getTimeFormat,
			Set:    setTimeFormat,
			Export: true,
			Doc:    "is the format for writing time as a string in box-format and the format for parsing time.",
		},
		"*flow-box-time-wrap*": {
			Get:    getTimeWrap,
			Set:    setTimeWrap,
			Export: true,
			Doc: `if non-nil then the writing and parsing of time is as a hash-map with a key of
the _*flow-box-time_wrap*_ value and the time encoded according to the _*flow-box-time-format*_.`,
		},
	})
	defBox()
	defCanLog()
	defDeleteFileActor()
	defExitActor()
	defFlow()
	defForeachActor()
	defGlobActor()
	defGroup()
	defHttpClientActor()
	defInspectActor()
	defJumpActor()
	defLogErrorActor()
	defMergeActor()
	defReadCsvActor()
	defReadFileActor()
	defReadJsonActor()
	defReadXmlActor()
	defSplitActor()
	defTaskActor()
	defTask()
	defTrack()
	defWriteFileActor()

	slip.DefConstant(slip.Symbol("*flow*"), &Pkg, "")
	Pkg.Initialize(nil, &event{}) // lock
	slip.AddPackage(&Pkg)
	slip.UserPkg.Use(&Pkg)
}

func getTimeFormat() slip.Object {
	if 0 < len(options.TimeFormat) {
		return slip.String(options.TimeFormat)
	}
	return nil
}

func setTimeFormat(value slip.Object) {
	switch tv := value.(type) {
	case nil:
		options.TimeFormat = ""
	case slip.Symbol:
		options.TimeFormat = string(tv)
	case slip.String:
		options.TimeFormat = string(tv)
	default:
		slip.PanicType("*bag-time-format*", value, "string")
	}
	updateConverter()
}

func getTimeWrap() slip.Object {
	if 0 < len(options.TimeWrap) {
		return slip.String(options.TimeWrap)
	}
	return nil
}

func setTimeWrap(value slip.Object) {
	switch tv := value.(type) {
	case nil:
		options.TimeWrap = ""
	case slip.Symbol:
		options.TimeWrap = string(tv)
	case slip.String:
		options.TimeWrap = string(tv)
	default:
		slip.PanicType("*bag-time-wrap*", value, "string")
	}
	updateConverter()
}

func updateConverter() {
	if len(options.TimeFormat) == 0 {
		options.Converter = nil
		return
	}
	if 0 < len(options.TimeWrap) {
		options.Converter = &alt.Converter{
			Map: []func(val map[string]interface{}) (interface{}, bool){
				func(val map[string]interface{}) (interface{}, bool) {
					if len(val) == 1 {
						switch tv := val[options.TimeWrap].(type) {
						case string:
							for _, layout := range []string{
								time.RFC3339Nano,
								time.RFC3339,
								"2006-01-02",
								options.TimeFormat,
							} {
								if t, err := time.ParseInLocation(layout, tv, time.UTC); err == nil {
									return t, true
								}
							}
						case int64:
							return time.Unix(0, tv).UTC(), true
						}
					}
					return val, false
				},
			},
		}
	} else {
		switch options.TimeFormat {
		case "nano":
			options.Converter = &ojg.TimeNanoConverter
		case time.RFC3339Nano, "rfc3339":
			options.Converter = &ojg.TimeRFC3339Converter
		case "second":
			options.Converter = &ojg.Converter{
				Float: []func(val float64) (interface{}, bool){
					func(val float64) (interface{}, bool) {
						if 946684800.0 <= val && val <= 2524608000.0 { // 2000-01-01 <= val <= 2050-01-01
							sec := int64(val)
							nano := int64((val - float64(sec)) * 1_000_000_000.0)
							return time.Unix(sec, nano).UTC(), true
						}
						return val, false
					},
				},
			}
		default:
			options.Converter = &ojg.Converter{
				String: []func(val string) (interface{}, bool){
					func(val string) (interface{}, bool) {
						if 6 <= len(val) { // minimal len for year month and day.
							if t, err := time.ParseInLocation(options.TimeFormat, val, time.UTC); err == nil {
								return t, true
							}
						}
						return val, false
					},
				},
			}
		}
	}
}

func methodDocFromFunc(method, funcName, flavor, obj string) string {
	var b []byte
	fd := slip.DescribeFunction(slip.Symbol(funcName))
	if fd != nil {
		b = fmt.Appendf(b, "__%s__ ", method)
		for _, da := range fd.Args[1:] { // first arg is always the instance
			if da.Name[0] == '&' {
				b = fmt.Appendf(b, "%s ", da.Name)
			} else {
				b = fmt.Appendf(b, "_%s_ ", da.Name)
			}
		}
		if 0 < len(fd.Return) {
			b = fmt.Appendf(b, "=> _%s_\n", fd.Return)
		} else {
			b = append(b, '\n')
		}
		for _, da := range fd.Args[1:] {
			if da.Name[0] != '&' {
				b = fmt.Appendf(b, "   _%s_ [%s] %s\n", da.Name, da.Type, da.Text)
			}
		}
		b = fmt.Appendf(b, "\n\nThe __%s__ method", method)
		b = append(b, fd.Text[strings.IndexByte(fd.Text, ' '):]...)
		if 0 < len(fd.Examples) {
			b = append(b, "\n\n\nExamples:\n"...)
			pat := fmt.Sprintf("(%s %s", funcName, obj)
			rep := fmt.Sprintf("(send %s %s", obj, method)
			for _, ex := range fd.Examples {
				b = append(b, ' ', ' ')
				b = append(b, strings.Replace(ex, pat, rep, 1)...)
				b = append(b, '\n')
			}
		}
		b = fmt.Appendf(b, "\n\nSee also: __%s__\n", funcName)
	}
	return string(b)
}
