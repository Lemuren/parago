package dag

import (
	"errors"
	"fmt"

	"github.com/Lemuren/parago/parser"
)

var ErrUnknownDep = errors.New("unknown dependency")
var ErrCyclicDep = errors.New("cyclic dependency")
var ErrUnknownId = errors.New("unknown task id")

// A 'Task' represents a single workflow task.
type Task struct {
	Id string	// task name (YAML key)
	Cmd string	// command to run
	Deps []*Task // other tasks this one depends on
}

// Find the index of a Task with the given id.
func IndexOf(id string, tasks []Task) (int, error) {
	for i, t := range tasks {
		if t.Id == id {
			return i, nil
		}
	}
	return 0, fmt.Errorf("id %s: %w", id, ErrUnknownId)
}

func Parse(input []parser.Task) ([]Task, error) {
	tasks := make([]Task, len(input))

	// Create the Task structs, but save the deps for later.
	for i, t := range input {
		tasks[i] = Task{Id: t.Id, Cmd: t.Cmd}
	}

	// With the structs created, we now insert the deps.
	// This is quite slow (we scan the list over and over again) but for now
	// it'll have to do.
	for i, t := range tasks {
		// For each "string" dependency, we find the task struct it matches to.
		deps := make([]*Task, len(input[i].Deps))
		for j, s := range input[i].Deps {
			k, err := IndexOf(s, tasks)
			if err != nil {
				return nil, fmt.Errorf("task %s: %w %s", t.Id, ErrUnknownDep, s)
			} else {
				deps[j] = &tasks[k]
			}
		}
		tasks[i].Deps = deps
	}

	return tasks, nil
}

func hasCyclicDependency() {
}
