/*
 * Copyright (c) 2020-present unTill Pro, Ltd.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package logger

// ILogger s.e.
type ILogger interface {
	Error(args ...interface{})
	Info(args ...interface{})
	Warning(args ...interface{})
	Debug(args ...interface{})
	IsDebug() bool
}
