/*
 * Copyright (c) 2020-present unTill Pro, Ltd.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package logger

import (
	"fmt"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

// TLogLevel s.e.
type TLogLevel int32

const (
	skipStackFramesCount = 4
	normalLineLength     = 60
)

// Log Levels enum
const (
	LogLevelNone = TLogLevel(iota)
	LogLevelError
	LogLevelWarning
	LogLevelInfo
	LogLevelDebug
)

const (
	errorPrefix   = "*****"
	warningPrefix = "!!!"
	infoPrefix    = "==="
	debugPrefix   = "---"
)

var globalLogPrinter = logPrinter{logLevel: LogLevelInfo}

// SetLogLevel s.e.
func SetLogLevel(logLevel TLogLevel) {
	atomic.StoreInt32((*int32)(&globalLogPrinter.logLevel), int32(logLevel))
}

// IsEnabled s.e.
func IsEnabled(logLevel TLogLevel) bool {
	curLogLevel := TLogLevel(atomic.LoadInt32((*int32)(&globalLogPrinter.logLevel)))
	return curLogLevel >= logLevel
}

// IsDebug s.e.
func IsDebug() bool {
	return IsEnabled(LogLevelDebug)
}

type logPrinter struct {
	logLevel TLogLevel
}

func (p *logPrinter) getFuncName(skipCount int) (funcName string, line int) {
	var fn string
	pc, _, line, ok := runtime.Caller(skipCount)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		elems := strings.Split(details.Name(), "/")
		if len(elems) > 0 {
			fn = elems[len(elems)-1]
		}
	}
	return fn, line
}

func (p *logPrinter) getFormattedMsg(msgType string, funcName string, line int, args ...interface{}) string {
	t := time.Now()
	out := fmt.Sprint(t.Format("01/02 15:04:05.000"))
	out += fmt.Sprint(": " + msgType)
	out += fmt.Sprintf(": [%v:%v]:", funcName, line)
	if len(args) > 0 {
		var s string
		for _, arg := range args {
			s = s + fmt.Sprint(" ", arg)
		}
		for i := len(s); i < normalLineLength; i++ {
			s = s + " "
		}
		out += fmt.Sprint(s)
	}
	return out
}

func (p *logPrinter) print(msgType string, args ...interface{}) {
	funcName, line := p.getFuncName(skipStackFramesCount)
	out := p.getFormattedMsg(msgType, funcName, line, args...)
	fmt.Println(out)
}

func getLevelPrefix(level TLogLevel) string {
	switch level {
	case LogLevelError:
		return errorPrefix
	case LogLevelWarning:
		return warningPrefix
	case LogLevelInfo:
		return infoPrefix
	case LogLevelDebug:
		return debugPrefix
	}
	return ""
}

func printIfLevel(level TLogLevel, args ...interface{}) {
	if IsEnabled(level) {
		globalLogPrinter.print(getLevelPrefix(level), args...)
	}
}

// Error s.e.
func Error(args ...interface{}) {
	printIfLevel(LogLevelError, args...)
}

// Warning s.e.
func Warning(args ...interface{}) {
	printIfLevel(LogLevelWarning, args...)
}

// Info s.e.
func Info(args ...interface{}) {
	printIfLevel(LogLevelInfo, args...)
}

// Debug s.e.
func Debug(args ...interface{}) {
	printIfLevel(LogLevelDebug, args...)
}
