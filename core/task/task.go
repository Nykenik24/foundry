package task

import (
	"log"

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

func TaskFromTOML(name, tomlSrc string) *Task {
	var task Task
	err := toml.Unmarshal([]byte(tomlSrc), &task)
	if err != nil {
		log.Fatalf("Error when unmarshaling task: %v", err)
	}

	return &task
}

func RetrieveTasks(tomlDoc string) *TasksDocument {
	var tasksDoc TasksDocument
	err := toml.Unmarshal([]byte(tomlDoc), &tasksDoc)
	if err != nil {
		log.Fatalf("Error when unmarshaling task document: %v", err)
	}

	return &tasksDoc
}
