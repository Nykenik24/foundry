/*
Copyright © 2026 Luca A. nykenik24@proton.me
*/

package cmd

import (
	"fmt"
	"time"

	"github.com/Nykenik24/foundry/core/runs"
	"github.com/Nykenik24/foundry/core/utils"
	"github.com/spf13/cobra"
)

var runsCmd = &cobra.Command{
	Use:   "runs",
	Short: "Manage runs",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var runsAllCmd = &cobra.Command{
	Use:   "all",
	Short: "List all runs",
	Run: func(cmd *cobra.Command, args []string) {
		list := runs.RetrieveRuns()
		runs.PrintRuns(list, fmt.Sprintf("Found %d run(s)", len(list)))
	},
}

var runsSuccessfulCmd = &cobra.Command{
	Use:   "success",
	Short: "List successful runs",
	Run: func(cmd *cobra.Command, args []string) {
		list := runs.QuerySuccessful(runs.RetrieveRuns()).Result
		runs.PrintRuns(list, fmt.Sprintf("Found %d succesful run(s)", len(list)))
	},
}

var runsFailedCmd = &cobra.Command{
	Use:   "failed",
	Short: "List failed runs",
	Run: func(cmd *cobra.Command, args []string) {
		list := runs.QueryFailed(runs.RetrieveRuns()).Result
		runs.PrintRuns(list, fmt.Sprintf("Found %d failed run(s)", len(list)))
	},
}

var runsBeforeCmd = &cobra.Command{
	Use:  "before <timestamp>",
	Args: cobra.ExactArgs(1),

	Short: "List runs before a timestamp.",
	Long: `List runs before a timestamp.

	Format for timestamp is YEAR-MONTH-DAY HOUR:MINUTE:SECOND`,
	Run: func(cmd *cobra.Command, args []string) {
		ts, err := time.Parse(time.DateTime, args[0])
		if err != nil {
			utils.PrintFatal("Error when parsing timestamp: %v", err)
		}
		list := runs.QueryBefore(runs.RetrieveRuns(), ts).Result
		runs.PrintRuns(list, fmt.Sprintf("Found %d run(s) before %s", len(list), args[0]))
	},
}

var runsAfterCmd = &cobra.Command{
	Use:  "after <timestamp>",
	Args: cobra.ExactArgs(1),

	Short: "List runs after a timestamp.",
	Long: `List runs after a timestamp.

	Format for timestamp is YEAR-MONTH-DAY HOUR:MINUTE:SECOND`,
	Run: func(cmd *cobra.Command, args []string) {
		ts, err := time.Parse(time.DateTime, args[0])
		if err != nil {
			utils.PrintFatal("Error when parsing timestamp: %v", err)
		}
		list := runs.QueryAfter(runs.RetrieveRuns(), ts).Result
		runs.PrintRuns(list, fmt.Sprintf("Found %d run(s) after %s", len(list), args[0]))
	},
}

func init() {
	runsCmd.AddCommand(runsAllCmd)
	runsCmd.AddCommand(runsSuccessfulCmd)
	runsCmd.AddCommand(runsFailedCmd)
	runsCmd.AddCommand(runsBeforeCmd)
	runsCmd.AddCommand(runsAfterCmd)

	rootCmd.AddCommand(runsCmd)
}
