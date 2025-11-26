package dag

// A 'Task' represents a single workflow task.
type Task struct {
	Id string	// task name (YAML key)
	Cmd string	// command to run
	Deps []string // other tasks this one depends on
}
