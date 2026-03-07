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

var showLogs bool

var runsAllCmd = &cobra.Command{
	Use:   "all",
	Short: "List all runs",
	Run: func(cmd *cobra.Command, args []string) {
		list := runs.RetrieveRuns()
		runs.PrintRuns(list, fmt.Sprintf("Found %d run(s)", len(list)), showLogs)
	},
}

var runsSuccessfulCmd = &cobra.Command{
	Use:   "success",
	Short: "List successful runs",
	Run: func(cmd *cobra.Command, args []string) {
		list := runs.QuerySuccessful(runs.RetrieveRuns()).Result
		runs.PrintRuns(list, fmt.Sprintf("Found %d succesful run(s)", len(list)), showLogs)
	},
}

var runsFailedCmd = &cobra.Command{
	Use:   "failed",
	Short: "List failed runs",
	Run: func(cmd *cobra.Command, args []string) {
		list := runs.QueryFailed(runs.RetrieveRuns()).Result
		runs.PrintRuns(list, fmt.Sprintf("Found %d failed run(s)", len(list)), showLogs)
	},
}

var timeLayouts = []string{
	time.DateTime,
	time.RFC3339,
	time.DateOnly,
	time.TimeOnly,
	time.Kitchen,
	"15:04",
	"2006",
}

func parseTimestamp(input string) (time.Time, error) {
	for _, layout := range timeLayouts {
		if t, err := time.Parse(layout, input); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("Invalid timestamp: %s", input)
}

var runsBeforeCmd = &cobra.Command{
	Use:  "before <timestamp>",
	Args: cobra.ExactArgs(1),

	Short: "List runs before a timestamp.",
	Long: `List runs before a timestamp.

	Format for timestamp is YEAR-MONTH-DAY HOUR:MINUTE:SECOND`,
	Run: func(cmd *cobra.Command, args []string) {
		ts, err := parseTimestamp(args[0])
		if err != nil {
			utils.PrintFatal("Error when parsing timestamp: %v", err)
		}
		list := runs.QueryBefore(runs.RetrieveRuns(), ts).Result
		runs.PrintRuns(list, fmt.Sprintf("Found %d run(s) before %s", len(list), args[0]), showLogs)
	},
}

var runsAfterCmd = &cobra.Command{
	Use:  "after <timestamp>",
	Args: cobra.ExactArgs(1),

	Short: "List runs after a timestamp.",
	Long: `List runs after a timestamp.

	Format for timestamp is YEAR-MONTH-DAY HOUR:MINUTE:SECOND`,
	Run: func(cmd *cobra.Command, args []string) {
		ts, err := parseTimestamp(args[0])
		if err != nil {
			utils.PrintFatal("Error when parsing timestamp: %v", err)
		}
		list := runs.QueryAfter(runs.RetrieveRuns(), ts).Result
		runs.PrintRuns(list, fmt.Sprintf("Found %d run(s) after %s", len(list), args[0]), showLogs)
	},
}

func init() {
	runsAllCmd.Flags().BoolVarP(&showLogs, "show-logs", "L", false, "Show run's logs")
	runsCmd.AddCommand(runsAllCmd)

	runsSuccessfulCmd.Flags().BoolVarP(&showLogs, "show-logs", "L", false, "Show run's logs")
	runsCmd.AddCommand(runsSuccessfulCmd)

	runsFailedCmd.Flags().BoolVarP(&showLogs, "show-logs", "L", false, "Show run's logs")
	runsCmd.AddCommand(runsFailedCmd)

	runsBeforeCmd.Flags().BoolVarP(&showLogs, "show-logs", "L", false, "Show run's logs")
	runsCmd.AddCommand(runsBeforeCmd)

	runsAfterCmd.Flags().BoolVarP(&showLogs, "show-logs", "L", false, "Show run's logs")
	runsCmd.AddCommand(runsAfterCmd)

	rootCmd.AddCommand(runsCmd)
}
