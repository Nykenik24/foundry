package task

import (
	"fmt"

	"github.com/pelletier/go-toml/v2"
)

type Task struct {
	Name string `toml:"name"`
	Cmd  string `toml:"cmd"`
}

type TasksDocument struct {
	Tasks map[string]Task `toml:"tasks"`
}

func NewTask(name, cmd string) *Task {
	return &Task{
		Name: name,
		Cmd:  cmd,
	}
}

func TaskFromTOML(name, tomlSrc string) (*Task, error) {
	var task Task
	err := toml.Unmarshal([]byte(tomlSrc), &task)
	if err != nil {
		return nil, fmt.Errorf("Error when unmarshaling task: %v", err)
	}

	return &task, nil
}

func RetrieveTasks(tomlDoc string) (*TasksDocument, error) {
	var tasksDoc TasksDocument
	err := toml.Unmarshal([]byte(tomlDoc), &tasksDoc)
	if err != nil {
		return nil, fmt.Errorf("Error when unmarshaling task document: %v", err)
	}

	return &tasksDoc, nil
}
