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
	LogLevelVerbose
	LogLevelDebug
)

// SetLogLevel s.e.
func SetLogLevel(logLevel TLogLevel) {
	atomic.StoreInt32((*int32)(&globalLogPrinter.logLevel), int32(logLevel))
}

// IsEnabled s.e.
func IsEnabled(logLevel TLogLevel) bool {
	curLogLevel := TLogLevel(atomic.LoadInt32((*int32)(&globalLogPrinter.logLevel)))
	return curLogLevel >= logLevel
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

func Debug(args ...interface{}) {
	printIfLevel(LogLevelDebug, args...)
}

func IsDebug() bool {
	return IsEnabled(LogLevelDebug)
}

func IsVerbose() bool {
	return IsEnabled(LogLevelVerbose)
}
