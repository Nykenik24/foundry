package service

import (
	"fmt"
	"os/exec"

	"github.com/Nykenik24/foundry/core/utils"
)

type Service struct {
	Name    string `toml:"name"`
	Cmd     string `toml:"cmd"`
	Process *exec.Cmd
	Running bool
	logger  *utils.Logger
}

func (s *Service) String() string {
	var status string
	switch s.Running {
	case true:
		status = "Running"
	case false:
		status = "Stopped"
	}
	return fmt.Sprintf("%s %v", s.Name, status)
}

type ServiceDoc struct {
	Services map[string]*Service `toml:"services"`
}
