/*
 * Copyright (c) 2020-present unTill Pro, Ltd. and Contributors
 * @author Maxim Geraskin
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package logger

import "sync/atomic"

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

func SetLogLevel(logLevel TLogLevel) {
	atomic.StoreInt32((*int32)(&globalLogPrinter.logLevel), int32(logLevel))
}

func Error(args ...interface{}) {
	printIfLevel(LogLevelError, args...)
}

func Warning(args ...interface{}) {
	printIfLevel(LogLevelWarning, args...)
}

func Info(args ...interface{}) {
	printIfLevel(LogLevelInfo, args...)
}

func Verbose(args ...interface{}) {
	printIfLevel(LogLevelVerbose, args...)
}

func Trace(args ...interface{}) {
	printIfLevel(LogLevelTrace, args...)
}

func IsError() bool {
	return isEnabled(LogLevelError)
}

func IsInfo() bool {
	return isEnabled(LogLevelInfo)
}

func IsWarning() bool {
	return isEnabled(LogLevelWarning)
}

func IsVerbose() bool {
	return isEnabled(LogLevelVerbose)
}

func IsTrace() bool {
	return isEnabled(LogLevelTrace)
}
