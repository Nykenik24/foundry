package runs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Nykenik24/foundry/core/utils"
	"github.com/pelletier/go-toml/v2"
)

type RunInfo struct {
	RanBy   string   `toml:"ran_by"`
	Logs    []string `toml:"logs"`
	Success bool     `toml:"success"`
}

type Run struct {
	Timestamp time.Time
	Info      RunInfo
}

func NewRun(timestamp time.Time, info RunInfo) *Run {
	return &Run{
		Timestamp: timestamp,
		Info:      info,
	}
}

func RetrieveRuns() []*Run {
	if !utils.PathExists(".foundry") {
		utils.PrintFatal("Not in a foundry project")
	}

	if !utils.PathExists(".foundry/runs") {
		if err := os.Mkdir(".foundry/runs", 0755); err != nil {
			utils.PrintFatal("Error creating runs directory: %v", err)
		}
	}

	entries, err := os.ReadDir(".foundry/runs")
	if err != nil {
		utils.PrintFatal("Error when reading runs directory: %v", err)
	}

	var runs []*Run

	for _, entry := range entries {
		name, ok := strings.CutSuffix(entry.Name(), ".toml")
		if !ok {
			continue
		}

		timestamp, err := time.Parse(time.RFC3339, name)
		if err != nil {
			continue
		}

		rawSrc, err := os.ReadFile(".foundry/runs/" + entry.Name())
		if err != nil {
			utils.PrintFatal("Error reading %s: %v", entry.Name(), err)
		}

		var info RunInfo
		if err := toml.Unmarshal(rawSrc, &info); err != nil {
			utils.PrintFatal("Error parsing %s: %v", entry.Name(), err)
		}

		runs = append(runs, NewRun(timestamp, info))
	}

	return runs
}

func PrintRuns(runs []*Run, firstLine string, showLogs bool) {
	fmt.Println(firstLine)

	for i, run := range runs {
		prefix := "├──"
		if i == len(runs)-1 {
			prefix = "└──"
		}

		var logN int
		for _, log := range run.Info.Logs {
			if log != "" {
				logN++
			}
		}

		var status string
		switch run.Info.Success {
		case true:
			status = "\033[32mSuccessful"
		case false:
			status = "\033[31mFailed"
		}
		fmt.Printf("%s Run \033[34m%s\033[0m (%s\033[0m): by \033[33m%s\033[0m (has \033[35m%d\033[0m logs)\n",
			prefix,
			run.Timestamp.Format(time.DateTime),
			status,
			run.Info.RanBy,
			logN,
		)
		if showLogs {
			for j, log := range run.Info.Logs {
				if log != "" {
					prefix := "    ├──"
					if j == len(runs)-1 {
						prefix = "    └──"
					}

					fmt.Printf("%s \033[32m%s\033[0m\n", prefix, log)
				}
			}
		}
	}
}
