package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/Nykenik24/foundry/core/utils"
)

type ServiceManager struct {
	services map[string]*Service
	mu       sync.Mutex
}

func NewManager() *ServiceManager {
	return &ServiceManager{
		services: make(map[string]*Service),
	}
}

func (m *ServiceManager) Start(name string) error {
	m.mu.Lock()
	svc, ok := m.services[name]
	if !ok {
		m.mu.Unlock()
		return fmt.Errorf("service not found")
	}

	if svc.Running {
		m.mu.Unlock()
		return fmt.Errorf("service already running")
	}
	m.mu.Unlock()

	if !ok {
		return fmt.Errorf("service not found")
	}

	if svc.Running {
		return fmt.Errorf("service already running")
	}

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "sh", "-c", svc.Cmd)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if utils.PathExists(".foundry") {
		if !utils.PathExists(".foundry/proc-logs") {
			os.Mkdir(".foundry/proc-logs", 0755)
		}
		logfile, _ := os.Create(fmt.Sprintf(".foundry/proc-logs/%s.log", svc.Name))
		cmd.Stdout = logfile
		cmd.Stderr = logfile

		svc.logger.SetNewOutput(logfile)
		svc.logger.MinimumLevel(utils.LogProc)
		svc.logger.Log("(Foundry) service started", utils.LogProc)
	} else {
		return fmt.Errorf("Not in a foundry project")
	}

	err := cmd.Start()
	if err != nil {
		return err
	}

	svc.Process = cmd

	svc.Running = true

	go func() {
		err := cmd.Wait()

		m.mu.Lock()
		defer m.mu.Unlock()

		svc.Running = false

		if err != nil {
			utils.GlobalLogger.Log(fmt.Sprintf("service %s exited: %v\n", svc.Name, err), utils.LogInfo)
		}
	}()

	return nil
}

func (m *ServiceManager) Stop(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	svc, ok := m.services[name]
	if !ok {
		return fmt.Errorf("service not found")
	}

	if !svc.Running {
		return fmt.Errorf("service not running")
	}

	pgid, err := syscall.Getpgid(svc.Process.Process.Pid)
	if err != nil {
		return err
	}

	err = syscall.Kill(-pgid, syscall.SIGTERM)
	if err != nil {
		return err
	}

	svc.Running = false
	return nil
}

func (m *ServiceManager) List() []Service {
	m.mu.Lock()
	defer m.mu.Unlock()

	result := []Service{}
	for _, svc := range m.services {
		result = append(result, *svc)
	}
	return result
}

func (m *ServiceManager) Register(name, cmd string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.services[name] = &Service{
		Name:   name,
		Cmd:    cmd,
		logger: utils.NewLogger(os.Stderr),
	}
}

func (m *ServiceManager) Restart(name string) error {
	if err := m.Stop(name); err != nil {
		return err
	}
	return m.Start(name)
}
