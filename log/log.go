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

func Dbg(str string) {
	if Level == DebugLevel || Level == InfoLevel {
		fmt.Println(_Colorize(str))
	}
}

func Info(str string) {
	if Level == InfoLevel {
		fmt.Println(_Colorize(str))
	}
}

func Err(str string, err error) {
	fmt.Println(_Colorize(str))
	fmt.Printf("%sError:%s %v\n", ErrorColor, ResetColor, err)
}
