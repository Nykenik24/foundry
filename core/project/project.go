package project

import (
	"fmt"
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

func NewProject(rootPath string) (*Project, error) {
	if !utils.PathExists(fmt.Sprintf("%s/.foundry", rootPath)) {
		return nil, fmt.Errorf(".foundry not found in root directory '%s'\n", rootPath)
	}

	if !utils.PathExists(fmt.Sprintf("%s/.foundry/project.toml", rootPath)) {
		return nil, fmt.Errorf("project.toml not found in foundry dir '%s/.foundry'\n", rootPath)
	}

	configBytes, err := os.ReadFile(fmt.Sprintf("%s/.foundry/project.toml\n", rootPath))
	if err != nil {
		return nil, fmt.Errorf("Error when retrieving configuration: %v\n", err)
	}

	var conf ProjectConfig
	err = toml.Unmarshal(configBytes, &conf)
	if err != nil {
		return nil, fmt.Errorf("Error when unmarshaling configuration: %v\n", err)
	}

	return &Project{
		RootPath: rootPath,
		Config:   conf,
	}, nil
}
