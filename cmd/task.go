/*
Copyright © 2026 Luca A. <Nykenik24@proton.me>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Nykenik24/foundry/core/task"
	"github.com/Nykenik24/foundry/core/utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manages tasks",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("task called")
	},
}

var newTaskCommand string

var taskNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new task",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatalf("Expected 1 arg (name) in 'task new', got %d", len(args))
		}

		if !utils.FileExists(".foundry") {
			log.Fatalf("Not a foundry project")
		}

		if !utils.FileExists(".foundry/tasks.toml") {
			if err := os.WriteFile(".foundry/tasks.toml", []byte(""), 0644); err != nil {
				log.Fatalf("Error when writing tasks document: %v", err)
			}
		}

		taskSrc, err := os.ReadFile(".foundry/tasks.toml")
		if err != nil {
			log.Fatalf("Error when reading .foundry/tasks.toml: %v", err)
		}

		name := args[0]
		tasks := task.RetrieveTasks(string(taskSrc))

		if tasks.Tasks == nil {
			tasks.Tasks = make(map[string]task.Task)
		}

		if _, exists := tasks.Tasks[name]; exists {
			log.Fatalf("Task '%s' already exists", name)
		}

		tasks.Tasks[name] = *task.NewTask(name, newTaskCommand)

		newSrc, err := toml.Marshal(tasks)
		if err != nil {
			log.Fatalf("Error when marshaling new task document: %v", err)
		}
		os.WriteFile(".foundry/tasks.toml", newSrc, 0644)
	},
}

var taskRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a task",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatalf("Expected 1 arg (name) in 'task remove', got %d", len(args))
		}

		if !utils.FileExists(".foundry") {
			log.Fatalf("Not a foundry project")
		}

		if !utils.FileExists(".foundry/tasks.toml") {
			if err := os.WriteFile(".foundry/tasks.toml", []byte(""), 0644); err != nil {
				log.Fatalf("Error when writing tasks document: %v", err)
			}
		}

		taskSrc, err := os.ReadFile(".foundry/tasks.toml")
		if err != nil {
			log.Fatalf("Error when reading .foundry/tasks.toml: %v", err)
		}

		name := args[0]
		tasks := task.RetrieveTasks(string(taskSrc))

		if tasks.Tasks == nil {
			tasks.Tasks = make(map[string]task.Task)
		}

		if _, exists := tasks.Tasks[name]; !exists {
			log.Fatalf("Task '%s' doesn't exist", name)
		}

		delete(tasks.Tasks, name)

		newSrc, err := toml.Marshal(tasks)
		if err != nil {
			log.Fatalf("Error when marshaling new task document: %v", err)
		}
		os.WriteFile(".foundry/tasks.toml", newSrc, 0644)
	},
}

func init() {
	taskNewCmd.Flags().StringVar(&newTaskCommand, "cmd", "", "The command the task will have")

	taskCmd.AddCommand(taskNewCmd)
	taskCmd.AddCommand(taskRemoveCmd)
	rootCmd.AddCommand(taskCmd)
}
