package service

import (
	"encoding/json"
	"net"
	"os"

	"github.com/Nykenik24/foundry/core/utils"
)

type Request struct {
	Action  string `json:"action"`
	Service string `json:"service"`
}

type ResponseService struct {
	Name   string `json:"name"`
	Status bool   `json:"status"`
}

type Response struct {
	Status   string            `json:"status,omitempty"`
	Services []ResponseService `json:"services,omitempty"`
	Error    string            `json:"error,omitempty"`
}

type Daemon struct {
	manager *ServiceManager
}

func NewDaemon(manager *ServiceManager) *Daemon {
	return &Daemon{manager: manager}
}

func (d *Daemon) Run(socket string) error {
	os.Remove(socket)

	l, err := net.Listen("unix", socket)
	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		go d.handle(conn)
	}
}

func (d *Daemon) handle(conn net.Conn) {
	defer conn.Close()

	var req Request
	if err := json.NewDecoder(conn).Decode(&req); err != nil {
		json.NewEncoder(conn).Encode(Response{
			Error: err.Error(),
		})
		return
	}

	switch req.Action {

	case "start":
		err := d.manager.Start(req.Service)
		if err != nil {
			json.NewEncoder(conn).Encode(Response{
				Error: err.Error(),
			})
			utils.PrintError("When starting service: %v", err)
			return
		}

		json.NewEncoder(conn).Encode(Response{
			Status: "ok",
		})

	case "stop":
		err := d.manager.Stop(req.Service)
		if err != nil {
			json.NewEncoder(conn).Encode(Response{
				Error: err.Error(),
			})
			utils.PrintError("When stopping service: %v", err)
			return
		}

		json.NewEncoder(conn).Encode(Response{
			Status: "stopped",
		})

	case "list":
		services := d.manager.List()

		var names []ResponseService
		for _, svc := range services {
			names = append(names, ResponseService{
				Name:   svc.Name,
				Status: svc.Running,
			})
		}

		json.NewEncoder(conn).Encode(Response{
			Services: names,
		})

	case "shutdown":
		for _, svc := range d.manager.List() {
			svc.logger.Log("Daemon shutdown", utils.LogProc)
		}

		json.NewEncoder(conn).Encode(Response{Status: "shutting down"})
		go func() {
			os.Exit(0)
		}()

	default:
		json.NewEncoder(conn).Encode(Response{
			Error: "unknown action",
		})
	}
}
