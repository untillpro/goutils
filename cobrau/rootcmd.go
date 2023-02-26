/*
* Copyright (c) 2023-present unTill Pro, Ltd.
* @author Maxim Geraskin
 */

package cobrau

import (
	"github.com/spf13/cobra"
	"github.com/untillpro/goutils/logger"
)

/*

Persistent flags:

  -v, --verbose   Print verbose output (detailed level)
      --trace     Print trace output   (most detailed level)

*/

func PrepareRootCmd(use string, short string, args []string, cmds ...*cobra.Command) *cobra.Command {

	var rootCmd = &cobra.Command{
		Use:   use,
		Short: "Cluster management utility written in golang",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if ok, _ := cmd.Flags().GetBool("trace"); ok {
				logger.SetLogLevel(logger.LogLevelTrace)
				logger.Verbose("Using logger.LogLevelTrace...")
			} else if ok, _ := cmd.Flags().GetBool("verbose"); ok {
				logger.SetLogLevel(logger.LogLevelVerbose)
				logger.Verbose("Using logger.LogLevelVerbose...")
			}
		},
	}

	rootCmd.SetArgs(args[1:])
	rootCmd.AddCommand(cmds...)
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().Bool("trace", false, "Extremely verbose output")
	rootCmd.SilenceUsage = true
	return rootCmd
}
