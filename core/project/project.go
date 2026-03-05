package project

import (
	"fmt"
	"log"
	"os"

	"github.com/Nykenik24/foundry/core/utils"
	"github.com/pelletier/go-toml/v2"
)

type ProjectConfig struct {
	Name string
}

type Project struct {
	RootPath string
	Config   ProjectConfig
}

func NewProject(rootPath string) *Project {
	if !utils.FileExists(fmt.Sprintf("%s/.foundry", rootPath)) {
		log.Fatalf(".foundry not found in root directory '%s'\n", rootPath)
	}

	if !utils.FileExists(fmt.Sprintf("%s/.foundry/project.toml", rootPath)) {
		log.Fatalf("project.toml not found in foundry dir '%s/.foundry'", rootPath)
	}

	configBytes, err := os.ReadFile(fmt.Sprintf("%s/.foundry/project.toml", rootPath))
	if err != nil {
		log.Fatalf("Error when retrieving configuration: %v", err)
	}

	var conf ProjectConfig
	err = toml.Unmarshal(configBytes, &conf)
	if err != nil {
		log.Fatalf("Error when unmarshaling configuration: %v", err)
	}

	return &Project{
		RootPath: rootPath,
		Config:   conf,
	}
}
