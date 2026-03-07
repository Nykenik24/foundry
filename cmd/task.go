/*
Copyright © 2026 Luca A. <Nykenik24@proton.me>
*/

package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Nykenik24/foundry/core/runs"
	"github.com/Nykenik24/foundry/core/task"
	"github.com/Nykenik24/foundry/core/utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manages tasks",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var newTaskCommand string
var newTaskStoreRuns bool

var taskNewCmd = &cobra.Command{
	Use:  "new <name>",
	Args: cobra.ExactArgs(1),

	Short: "Creates a new task",
	Run: func(cmd *cobra.Command, args []string) {
		if !utils.PathExists(".foundry") {
			utils.PrintFatal("Not a foundry project")
		}

		if !utils.PathExists(".foundry/tasks.toml") {
			if err := os.WriteFile(".foundry/tasks.toml", []byte(""), 0644); err != nil {
				utils.PrintFatal("Error when writing tasks document: %v", err)
			}
		}

		taskSrc, err := os.ReadFile(".foundry/tasks.toml")
		if err != nil {
			utils.PrintFatal("Error when reading .foundry/tasks.toml: %v", err)
		}

		name := args[0]
		tasks, err := task.RetrieveTasks(string(taskSrc))
		if err != nil {
			utils.PrintFatal("%s", err.Error())
		}

		if tasks.Tasks == nil {
			tasks.Tasks = make(map[string]*task.Task)
		}

		if _, exists := tasks.Tasks[name]; exists {
			utils.PrintFatal("Task '%s' already exists", name)
		}

		tasks.Tasks[name] = task.NewTask(name, newTaskCommand)
		tasks.Tasks[name].StoreRuns = newTaskStoreRuns

		newSrc, err := toml.Marshal(tasks)
		if err != nil {
			utils.PrintFatal("Error when marshaling new task document: %v", err)
		}
		os.WriteFile(".foundry/tasks.toml", newSrc, 0644)
	},
}

var taskRemoveCmd = &cobra.Command{
	Use:  "remove <name>",
	Args: cobra.ExactArgs(1),

	Short: "Remove a task",
	Run: func(cmd *cobra.Command, args []string) {
		if !utils.PathExists(".foundry") {
			utils.PrintFatal("Not a foundry project")
		}

		if !utils.PathExists(".foundry/tasks.toml") {
			if err := os.WriteFile(".foundry/tasks.toml", []byte(""), 0644); err != nil {
				utils.PrintFatal("Error when writing tasks document: %v", err)
			}
		}

		taskSrc, err := os.ReadFile(".foundry/tasks.toml")
		if err != nil {
			utils.PrintFatal("Error when reading .foundry/tasks.toml: %v", err)
		}

		name := args[0]
		tasks, err := task.RetrieveTasks(string(taskSrc))
		if err != nil {
			utils.PrintFatal("%s", err.Error())
		}

		if tasks.Tasks == nil {
			tasks.Tasks = make(map[string]*task.Task)
		}

		if _, exists := tasks.Tasks[name]; !exists {
			utils.PrintFatal("Task '%s' doesn't exist", name)
		}

		delete(tasks.Tasks, name)

		newSrc, err := toml.Marshal(tasks)
		if err != nil {
			utils.PrintFatal("Error when marshaling new task document: %v", err)
		}
		os.WriteFile(".foundry/tasks.toml", newSrc, 0644)
	},
}

var taskRunCmd = &cobra.Command{
	Use:  "run <name>",
	Args: cobra.ExactArgs(1),

	Short: "Run a task",
	Run: func(cmd *cobra.Command, args []string) {
		if !utils.PathExists(".foundry") {
			utils.PrintFatal("Not a foundry project")
		}

		if !utils.PathExists(".foundry/tasks.toml") {
			if err := os.WriteFile(".foundry/tasks.toml", []byte(""), 0644); err != nil {
				utils.PrintFatal("Error when writing tasks document: %v", err)
			}
		}

		taskSrc, err := os.ReadFile(".foundry/tasks.toml")
		if err != nil {
			utils.PrintFatal("Error when reading .foundry/tasks.toml: %v", err)
		}

		name := args[0]
		tasks, err := task.RetrieveTasks(string(taskSrc))
		if err != nil {
			utils.PrintFatal("%s", err.Error())
		}

		if tasks.Tasks == nil {
			tasks.Tasks = make(map[string]*task.Task)
		}

		if _, exists := tasks.Tasks[name]; !exists {
			utils.PrintFatal("Task '%s' doesn't exist", name)
		}

		targetTask := tasks.Tasks[name]
		cmdStr := exec.Command("sh", "-c", targetTask.Cmd)

		var stdout bytes.Buffer
		var stderr bytes.Buffer

		cmdStr.Stdout = &stdout
		cmdStr.Stderr = &stderr
		cmdStr.Stdin = os.Stdin

		err = cmdStr.Run()
		os.Stdout.Write([]byte(stdout.String()))
		os.Stderr.Write([]byte(stderr.String()))

		if err != nil {
			log.Printf("Error when running task: %v\n", err)
			utils.PrintFatal("Command was '%s'", tasks.Tasks[name].Cmd)
		}

		if targetTask.StoreRuns {
			utils.GlobalLogger.Log(fmt.Sprintf("Storing %s's run", targetTask.Name), utils.LogInfo)
			if !utils.PathExists(".foundry/runs") {
				os.Mkdir(".foundry/runs", 0755)
			}

			taskLogs := strings.Split(stdout.String()+"\n"+stderr.String(), "\n")

			runInfo := runs.RunInfo{
				RanBy:   "task/" + targetTask.Name,
				Logs:    taskLogs,
				Success: true,
			}

			runToml, err := toml.Marshal(runInfo)
			if err != nil {
				utils.PrintFatal("Error when marshaling run: %v", err)
			}

			utils.GlobalLogger.Log(fmt.Sprintf("TOML for run: %s", string(runToml)), utils.LogInfo)
			utils.GlobalLogger.Log(fmt.Sprintf("Name for run: %s.run", time.Now().Format(time.RFC3339)), utils.LogInfo)

			err = os.WriteFile(fmt.Sprintf(".foundry/runs/%s.toml", time.Now().Format(time.RFC3339)), runToml, 0644)
			if err != nil {
				utils.PrintFatal("Error when creating run file: %v", err)
			}
		}
	},
}

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Run: func(cmd *cobra.Command, args []string) {
		if !utils.PathExists(".foundry") {
			utils.PrintFatal("Not a foundry project")
		}

		if !utils.PathExists(".foundry/tasks.toml") {
			if err := os.WriteFile(".foundry/tasks.toml", []byte(""), 0644); err != nil {
				utils.PrintFatal("Error when writing tasks document: %v", err)
			}
		}

		taskSrc, err := os.ReadFile(".foundry/tasks.toml")
		if err != nil {
			utils.PrintFatal("Error when reading .foundry/tasks.toml: %v", err)
		}

		tasks, err := task.RetrieveTasks(string(taskSrc))
		if err != nil {
			utils.PrintFatal("%s", err.Error())
		}

		if tasks.Tasks == nil {
			tasks.Tasks = make(map[string]*task.Task)
		}

		fmt.Printf("Found %d task(s)\n", len(tasks.Tasks))
		for name, task := range tasks.Tasks {
			fmt.Printf("└── Task \033[34m'%s'\033[0m: \033[32m%s\033[0m\n", name, task.Cmd)
		}
	},
}

func init() {
	taskNewCmd.Flags().StringVar(&newTaskCommand, "cmd", "", "The command the task will have")
	taskNewCmd.Flags().BoolVar(&newTaskStoreRuns, "store-runs", false, "Store runs of the task")

	taskCmd.AddCommand(taskNewCmd)
	taskCmd.AddCommand(taskRemoveCmd)
	taskCmd.AddCommand(taskRunCmd)
	taskCmd.AddCommand(taskListCmd)
	rootCmd.AddCommand(taskCmd)
}
