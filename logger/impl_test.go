/*
 * Copyright (c) 2020-present unTill Pro, Ltd.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

package logger

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_BasicUsage(t *testing.T) {

	// "Hello world"
	{
		Error("Hello world", "arg1", "arg2")
		Warning("My warning")
		Info("My info")

		// IsDebug() is used to avoid unnecessary calculations
		if IsDebug() {
			Debug("!!! You should NOT see it since default level is INFO")
		}
	}
	// Changing LogLevelgolan
	{
		SetLogLevel(LogLevelDebug)
		if IsDebug() {
			Debug("Now you should see my Debug")
		}
		SetLogLevel(LogLevelError)
		Debug("!!! You should NOT see my Debug")
		Warning("!!! You should NOT see my warning")
		SetLogLevel(LogLevelInfo)
		Warning("You should see my warning")
		Warning("You should see my info")
	}
	// Let see how it looks when using from methods
	{
		m := mystruct{}
		m.iWantToLog()
	}
}

func Test_MsgFormatter(t *testing.T) {
	var out string

	out = globalLogPrinter.getFormattedMsg("", "sync_op.doSync", 120, "line1")
	assert.True(t, strings.Contains(out, ": [sync_op.doSync:120]: line1"))

	out = globalLogPrinter.getFormattedMsg("", "", 121, "line1", "line2")
	assert.True(t, strings.Contains(out, ": [:121]: line1 line2"))

	out = globalLogPrinter.getFormattedMsg("m1:m2/m3", "sync_op.doSync", 126, "line1", "line2", "line3")
	assert.True(t, strings.Contains(out, "m1:m2/m3: [sync_op.doSync:126]: line1 line2 line3"))

	out = globalLogPrinter.getFormattedMsg("m1:m2/m3", "sync_op.doSync", 127, "line/1", "line/2", "line/3")
	assert.True(t, strings.Contains(out, "m1:m2/m3: [sync_op.doSync:127]: line/1 line/2 line/3"))
}

func Test_CheckSetLevels(t *testing.T) {
	// 1. Info LogLevel = LogLevelInfo
	SetLogLevel(LogLevelInfo)
	assert.False(t, IsDebug())
	assert.True(t, IsEnabled(LogLevelInfo))
	assert.True(t, IsEnabled(LogLevelWarning))
	assert.True(t, IsEnabled(LogLevelError))

	// 2. Debug LogLevel = LogLevelDebug
	SetLogLevel(LogLevelDebug)
	assert.True(t, IsDebug())
	assert.True(t, IsEnabled(LogLevelInfo))
	assert.True(t, IsEnabled(LogLevelWarning))
	assert.True(t, IsEnabled(LogLevelError))

	// 3. Warning LogLevel = LogLevelWarning
	SetLogLevel(LogLevelWarning)
	assert.False(t, IsDebug())
	assert.False(t, IsEnabled(LogLevelInfo))
	assert.True(t, IsEnabled(LogLevelWarning))
	assert.True(t, IsEnabled(LogLevelError))

	// 4. Error LogLevel = LogLevelError
	SetLogLevel(LogLevelError)
	assert.False(t, IsDebug())
	assert.False(t, IsEnabled(LogLevelInfo))
	assert.False(t, IsEnabled(LogLevelWarning))
	assert.True(t, IsEnabled(LogLevelError))
}

func Test_CheckRightPrefix(t *testing.T) {
	// 1. Info LogLevel = LogLevelInfo
	SetLogLevel(LogLevelInfo)
	assert.Equal(t, getLevelPrefix(globalLogPrinter.logLevel), infoPrefix)

	// 2. Debug LogLevel = LogLevelDebug
	SetLogLevel(LogLevelDebug)
	assert.Equal(t, getLevelPrefix(globalLogPrinter.logLevel), debugPrefix)

	// 3. Warning LogLevel = LogLevelWarning
	SetLogLevel(LogLevelWarning)
	assert.Equal(t, getLevelPrefix(globalLogPrinter.logLevel), warningPrefix)

	// 4. Error LogLevel = LogLevelError
	SetLogLevel(LogLevelError)
	assert.Equal(t, getLevelPrefix(globalLogPrinter.logLevel), errorPrefix)

	// 5. Unexisting level
	SetLogLevel(7)
	assert.Equal(t, getLevelPrefix(globalLogPrinter.logLevel), "")

	SetLogLevel(LogLevelInfo)
}

func Test_GetFuncName(t *testing.T) {
	funcName, line := globalLogPrinter.getFuncName(2)
	assert.Equal(t, funcName, "testing.tRunner")
	assert.True(t, line > 0)
}

type mystruct struct {
}

func (m *mystruct) iWantToLog() {
	Error("OOPS")
}

func TestMultithread(t *testing.T) {
	require := require.New(t)
	r, w, err := os.Pipe()
	require.Nil(err)
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()
	os.Stdout = w
	wg := sync.WaitGroup{}
	wg.Add(1000)

	toLog := []string{}
	for i := 0; i < 100; i++ {
		toLog = append(toLog, strings.Repeat(strconv.Itoa(i), 10))
	}

	for i := 0; i < 1000; i++ {
		go func() {
			for i := 0; i < 100; i++ {
				Info(toLog[i])
			}
			wg.Done()
		}()
	}

	stdout := ""
	wait := make(chan struct{})
	go func() {
		buf := bytes.NewBuffer(nil)
		_, err := io.Copy(buf, r)
		require.Nil(err)
		stdout = buf.String()
		close(wait)
	}()
	wg.Wait()
	w.Close()
	<-wait

	logged := strings.Split(stdout, "\n")
outer:
	for _, loggedActual := range logged {
		if len(loggedActual) == 0 {
			continue
		}
		for _, loggedExpected := range toLog {
			if strings.Contains(loggedActual, loggedExpected) {
				continue outer
			}
		}
		t.Fatal(loggedActual)
	}
}

func Benchmark_FuncForPC(b *testing.B) {
	var funcName string

	start := time.Now()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		funcName, _ = globalLogPrinter.getFuncName(2)
	}
	assert.Equal(b, funcName, "testing.(*B).runN")

	elapsed := time.Since(start).Seconds()
	b.ReportMetric(float64(b.N)/elapsed, "rps")
}
