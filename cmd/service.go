/*
Copyright © 2026 Luca A. Nykenik24@proton.me
*/

package cmd

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/Nykenik24/foundry/core/service"
	"github.com/Nykenik24/foundry/core/utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

var socket string = "/tmp/foundry.sock"

func daemonRunning() bool {
	_, err := net.Dial("unix", socket)
	return err == nil
}

func ensureDaemon() error {
	if daemonRunning() {
		return nil
	}

	cmd := exec.Command(os.Args[0], "service", "daemon")

	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	for range 10 {
		time.Sleep(200 * time.Millisecond)

		_, err := net.Dial("unix", socket)
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("daemon failed to start")
}

func initManagerWithConfig() *service.ServiceManager {
	if !utils.PathExists(".foundry") {
		utils.PrintFatal("Not in a foundry project")
	}

	if !utils.PathExists(".foundry/services.toml") {
		os.Create(".foundry/services.toml")
		return service.NewManager()
	}

	rawSrc, err := os.ReadFile(".foundry/services.toml")
	if err != nil {
		utils.PrintFatal("%v", err)
	}

	var serviceDoc service.ServiceDoc
	toml.Unmarshal([]byte(rawSrc), &serviceDoc)
	if serviceDoc.Services == nil {
		serviceDoc.Services = make(map[string]*service.Service)
	}

	manager := service.NewManager()
	for _, service := range serviceDoc.Services {
		manager.Register(service.Name, service.Cmd)
	}
	return manager
}

func shutdownDaemon() {
	conn, err := net.Dial("unix", "/tmp/foundry.sock")
	if err != nil {
		utils.PrintFatal("Could not connect to daemon: %v", err)
	}
	defer conn.Close()

	req := service.Request{
		Action: "shutdown",
	}

	if err := json.NewEncoder(conn).Encode(req); err != nil {
		utils.PrintFatal("Failed to send shutdown: %v", err)
	}

	var resp service.Response
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		utils.PrintFatal("Failed to read response: %v", err)
	}

	os.Remove("/tmp/foundry.sock")
	utils.GlobalLogger.Log("Shutting daemon down...", utils.LogInfo)
}

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage services",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var serviceDaemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start the service daemon",
	Run: func(cmd *cobra.Command, args []string) {
		utils.GlobalLogger.Log("Started daemon. It is suggested to fork this command.", utils.LogInfo)

		manager := initManagerWithConfig()

		daemon := service.NewDaemon(manager)
		err := daemon.Run("/tmp/foundry.sock")
		if err != nil {
			shutdownDaemon()
			utils.PrintFatal("%v", err)
		}
	},
}

var serviceStartCmd = &cobra.Command{
	Use:  "start <name>",
	Args: cobra.ExactArgs(1),

	Short: "Start a service",
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureDaemon(); err != nil {
			utils.PrintFatal("%v", err)
		}

		conn, err := net.Dial("unix", "/tmp/foundry.sock")
		if err != nil {
			utils.PrintFatal("%v", err)
		}
		defer conn.Close()

		req := service.Request{
			Action:  "start",
			Service: args[0],
		}

		err = json.NewEncoder(conn).Encode(req)
		if err != nil {
			utils.PrintFatal("%v", err)
		}
	},
}

var serviceStopCmd = &cobra.Command{
	Use:  "stop <name>",
	Args: cobra.ExactArgs(1),

	Short: "Stop a service",
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureDaemon(); err != nil {
			utils.PrintFatal("%v", err)
		}

		conn, err := net.Dial("unix", "/tmp/foundry.sock")
		if err != nil {
			utils.PrintFatal("%v", err)
		}
		defer conn.Close()

		req := service.Request{
			Action:  "stop",
			Service: args[0],
		}

		err = json.NewEncoder(conn).Encode(req)
		if err != nil {
			utils.PrintFatal("%v", err)
		}
	},
}

var serviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all services",
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureDaemon(); err != nil {
			utils.PrintFatal("%v", err)
		}

		conn, err := net.Dial("unix", "/tmp/foundry.sock")
		if err != nil {
			utils.PrintFatal("%v", err)
		}
		defer conn.Close()

		req := service.Request{
			Action: "list",
		}

		err = json.NewEncoder(conn).Encode(req)
		if err != nil {
			utils.PrintFatal("%v", err)
		}

		var resp service.Response
		err = json.NewDecoder(conn).Decode(&resp)
		if err != nil {
			utils.PrintFatal("%v", err)
		}

		if resp.Error != "" {
			utils.PrintFatal("%s", resp.Error)
		}

		fmt.Printf("Found %d service(s)\n", len(resp.Services))
		for _, svc := range resp.Services {
			var status string
			switch svc.Status {
			case true:
				status = "\033[32mRunning"
			case false:
				status = "\033[31mStopped"
			}
			fmt.Printf("└── Service \033[34m'%s'\033[0m: %s\033[0m\n", svc.Name, status)
		}
	},
}

var serviceKillCmd = &cobra.Command{
	Use:   "dkill",
	Short: "Stop the service daemon",
	Run: func(cmd *cobra.Command, args []string) {
		shutdownDaemon()
	},
}

var serviceIsActiveCmd = &cobra.Command{
	Use:   "isactive",
	Short: "Check if daemon is running",
	Run: func(cmd *cobra.Command, args []string) {
		if daemonRunning() {
			fmt.Println("Daemon is \033[32mactive\033[0m.")
		} else {
			fmt.Println("Daemon is \033[31minactive\033[0m.")
			fmt.Println("Use \033[34m`foundry service daemon`\033[0m")
		}
	},
}

func init() {
	serviceCmd.AddCommand(serviceDaemonCmd)
	serviceCmd.AddCommand(serviceStartCmd)
	serviceCmd.AddCommand(serviceStopCmd)
	serviceCmd.AddCommand(serviceListCmd)
	serviceCmd.AddCommand(serviceKillCmd)
	serviceCmd.AddCommand(serviceIsActiveCmd)

	rootCmd.AddCommand(serviceCmd)
}
