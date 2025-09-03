package log

import (
	"fmt"
	"strings"
)

type LEVEL int

const (
	DebugLevel LEVEL = iota
	InfoLevel
	NoLogLevel
)

type COLOR string

const (
	SuccessColor COLOR = "\033[32m"
	ErrorColor   COLOR = "\033[31m"
	DebugColor   COLOR = "\033[33m"
	InfoColor    COLOR = "\033[34m"
	ResetColor   COLOR = "\033[0m"
)

type TAG string

const (
	SuccessTag TAG = "<s>"
	InfoTag    TAG = "<i>"
	ErrorTag   TAG = "<e>"
	DebugTag   TAG = "<d>"
	RessetTag  TAG = "</>"
)

var (
	Level = InfoLevel
)

func _Colorize(text string) string {
	// success
	text = strings.ReplaceAll(text, string(SuccessTag), string(SuccessColor))

	// error
	text = strings.ReplaceAll(text, string(ErrorTag), string(ErrorColor))

	// info
	text = strings.ReplaceAll(text, string(InfoTag), string(InfoColor))

	// debug
	text = strings.ReplaceAll(text, string(DebugTag), string(DebugColor))

	// reset
	text = strings.ReplaceAll(text, string(RessetTag), string(ResetColor))

	return text
}

func _ProccessArgs(prefix string, argn ...any) {
	var args []any
	args = append(args, argn...)

	fmt.Print(_Colorize(prefix + " "))

	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			fmt.Print(_Colorize(v))
		case []byte:
			fmt.Print(_Colorize(fmt.Sprintf("[<d>LENGTH:</> <i>%d</>] ", len(v))))

			for i, b := range v {
				start := "<d>"
				if i%2 == 0 {
					start = "<i>"
				}

				fmt.Print(_Colorize(start + fmt.Sprintf("%02d ", b) + "</>"))
			}
		case byte:
			fmt.Print(_Colorize(fmt.Sprintf("<d>BYTE:</> <i>%d</>", v)))
		case error:
			fmt.Print(_Colorize(fmt.Sprintf("<e>Error:</> %v", v)))
		// case tpkt.Header:
		// 	fmt.Print(_Colorize(fmt.Sprintf("<i>[TPKT-HEADER]</> Version: <d>%d</>\tReserved: <d>%d</>\tLength: <d>%d</>", v.Version, v.Reserved, v.Length)))
		// case x224.Header:
		default:
			fmt.Print(_Colorize("<e>[UNKNOWN LOG TYPE]</>"), v)
		}
	}

	fmt.Print("\n")
}

func Zebra(text string, color COLOR) {
	for i, s := range []byte(text) {
		if i%2 == 0 {
			fmt.Print(string(color) + string(s) + string(ResetColor))
		} else {
			fmt.Print(string(s))
		}
	}

	fmt.Print("\n")
}

func Dbg(arg1 any, argn ...any) {
	if Level == DebugLevel || Level == InfoLevel {
		prefix := "<d>[DEBUG]</>"
		var args []any
		args = append(args, arg1)
		args = append(args, argn...)

		_ProccessArgs(prefix, args...)
	}
}

func Info(arg1 any, argn ...any) {
	if Level == InfoLevel {
		prefix := "<i>[INFO]</>"
		var args []any
		args = append(args, arg1)
		args = append(args, argn...)

		_ProccessArgs(prefix, args...)
	}
}

func Err(arg1 any, argn ...any) {
	prefix := "<e>[ERROR]</> "
	var args []any
	args = append(args, arg1)
	args = append(args, argn...)

	_ProccessArgs(prefix, args...)
}
