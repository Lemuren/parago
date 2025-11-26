package parser

import (
	"io"
	"fmt"
	"errors"

	"github.com/goccy/go-yaml"
)

var ErrMissingCmd  = errors.New("cmd fields missing")
var ErrCmdNotString = errors.New("cmd is not a string")
var ErrInvalidDef = errors.New("invalid definition")
var ErrDepsNotArray = errors.New("deps is not an array")
var ErrDepsNotArrayOfStrings = errors.New("deps is not an array of strings")

// The parser Task struct resembles the dag one, but has strings in
// its dependencies instead of proper Task structs. This is to keep the
// unmarshaling from YAML simple. When we start constructing the DAG these
// structs will be used as input to the more 'proper' DAG Task structs.
type Task struct {
	Id string	// task name (YAML key)
	Cmd string	// command to run
	Deps []string // other tasks this one depends on
}

func Parse(r io.Reader) ([]Task, error) {
	// First we parse the raw YAML.
	var raw map[string]interface{}
	dec := yaml.NewDecoder(r)
	if err := dec.Decode(&raw); err != nil {
		return nil, err
	}

	// Next we'll unmarshal each of the tasks from the raw
	// input into proper Task structs.
	var tasks []Task = make([]Task, 0)
	for id, def := range raw {
		// There's two primary styles of defining tasks, full and shorthand.
		switch v := def.(type) {
		// Shorthand style: id: "command".
		case string:
			t := Task{Id: id, Cmd: v, Deps: make([]string, 0)}
			tasks = append(tasks, t)

		// Full definition style.
		case map[string]interface{}:
			t := Task{Id: id, Deps: make([]string, 0)}

			if cmd, ok := v["cmd"]; ok {
				if scmd, ok := cmd.(string); ok {
					t.Cmd = scmd
				} else {
					return nil, fmt.Errorf("task %s: %w", id, ErrCmdNotString)
				}
			} else {
				return nil, fmt.Errorf("task %s: %w", id, ErrMissingCmd)
			}

			if deps, ok := v["deps"]; ok {
				if sdeps, ok := deps.([]interface{}); ok {
					for _, d := range sdeps {
						if sd, ok := d.(string); ok {
							t.Deps = append(t.Deps, sd)
						} else {
							return nil, fmt.Errorf("task %s: %w", id, ErrDepsNotArrayOfStrings)
						}
					}
				} else {
					return nil, fmt.Errorf("task %s: %w", id, ErrDepsNotArray)
				}
			}

			tasks = append(tasks, t)
		
		// We shouldn't end up here, it means the user gave something like
		// build: 3 as a task.
		default:
			return nil, fmt.Errorf("task %s: %w", id, ErrInvalidDef)
		}
	}

	return tasks, nil
}
