/*
* Copyright (c) 2023-present unTill Pro, Ltd.
* @author Maxim Geraskin
 */

package cobrau

import (
	"context"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/untillpro/goutils/logger"
)

func ExecCommandAndCatchInterrupt(cmd *cobra.Command) error {

	done := make(chan struct{})

	cmdExec := func(ctx context.Context) (err error) {
		err = cmd.ExecuteContext(ctx)
		done <- struct{}{}
		return
	}

	return goAndCatchInterrupt(cmdExec, done)

}

func goAndCatchInterrupt(f func(ctx context.Context) error, done chan struct{}) (err error) {

	var signals = make(chan os.Signal, 1)

	ctx, cancel := context.WithCancel(context.Background())
	signal.Notify(signals, os.Interrupt)

	go func() {
		err = f(ctx)
	}()

infinitcycle:
	for {
		select {
		case sig := <-signals:
			if ctx.Err() == nil {
				logger.Verbose("signal received:", sig)
				cancel()
			}
		case <-done:
			logger.Verbose("item received from `done` channel")
			cancel()
			break infinitcycle
		}
	}
	return
}
