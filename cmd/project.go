/*
Copyright © 2026 Luca A. <Nykenik24@proton.me>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Nykenik24/foundry/core/utils"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manages projects",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("project called")
	},
}

var overwriteProject bool

var projectInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new project",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatalf("Expected 1 arg (name) in 'project init', got %d args", len(args))
		}

		name := args[0]
		if !utils.FileExists(".foundry") {
			if err := os.Mkdir(".foundry", 0755); err != nil {
				log.Fatalf("Failed to create .foundry directory: %v", err)
			}
		}
		rawConf := fmt.Sprintf(`name = "%s"`, name)

		if utils.FileExists(".foundry/project.toml") && !overwriteProject {
			log.Fatalf("Project already exists, consider removing .foundry/project.toml or using '--overwrite' (or '-o')")
		} else {
			// os.Create(".foundry/project.toml")
			if err := os.WriteFile(".foundry/project.toml", []byte(rawConf), 0644); err != nil {
				log.Fatalf("Failed to write project.toml: %v", err)
			}
		}
	},
}

func init() {
	projectInitCmd.Flags().BoolVarP(&overwriteProject, "overwrite", "o", false, "overwrite existing project (if there is one)")
	projectCmd.AddCommand(projectInitCmd)

	rootCmd.AddCommand(projectCmd)
}
