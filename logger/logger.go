/*
 * Copyright (c) 2020-present unTill Pro, Ltd. and Contributors
 * @author Maxim Geraskin
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package logger

import "sync/atomic"

// external log printing function
var ExtPrintFunc func(level TLogLevel, args ...interface{})

// TLogLevel s.e.
type TLogLevel int32

// Log Levels enum
const (
	LogLevelNone = TLogLevel(iota)
	LogLevelError
	LogLevelWarning
	LogLevelInfo
	LogLevelVerbose // aka Debug
	LogLevelTrace
)

func print(level TLogLevel, args ...interface{}) {
	if ExtPrintFunc != nil {
		ExtPrintFunc(level, args...)
	} else {
		printIfLevel(level, args...)
	}
}

func SetLogLevel(logLevel TLogLevel) {
	atomic.StoreInt32((*int32)(&globalLogPrinter.logLevel), int32(logLevel))
}

func Error(args ...interface{}) {
	print(LogLevelError, args...)
}

func Warning(args ...interface{}) {
	print(LogLevelWarning, args...)
}

func Info(args ...interface{}) {
	print(LogLevelInfo, args...)
}

func Verbose(args ...interface{}) {
	print(LogLevelVerbose, args...)
}

func Trace(args ...interface{}) {
	print(LogLevelTrace, args...)
}

func IsError() bool {
	return IsEnabled(LogLevelError)
}

func IsInfo() bool {
	return IsEnabled(LogLevelInfo)
}

func IsWarning() bool {
	return IsEnabled(LogLevelWarning)
}

func IsVerbose() bool {
	return IsEnabled(LogLevelVerbose)
}

func IsTrace() bool {
	return IsEnabled(LogLevelTrace)
}
