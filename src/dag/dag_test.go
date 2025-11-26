package dag

import (
	"testing"
	"errors"

	"github.com/Lemuren/parago/parser"
)

func Test_IndexOf(t *testing.T) {
	tasks := []Task {
		{Id: "id1", Cmd: "cmd1"},
		{Id: "id2", Cmd: "cmd2"},
		{Id: "id3", Cmd: "cmd3"},
		{Id: "id4", Cmd: "cmd4"},
		{Id: "id5", Cmd: "cmd5"},
	}

	for i, task := range tasks {
		got, err := IndexOf(task.Id, tasks)
		if err != nil {
			t.Fatalf("error not nil when looking for %s", task.Id)
		}
		if got != i {
			t.Errorf("expected %d, got %d", i, got)
		}
	}
}

func Test_IndexOfUnknown(t *testing.T) {
	tasks := []Task {
		{Id: "id1", Cmd: "cmd1"},
		{Id: "id2", Cmd: "cmd2"},
		{Id: "id3", Cmd: "cmd3"},
		{Id: "id4", Cmd: "cmd4"},
		{Id: "id5", Cmd: "cmd5"},
	}

	_, err := IndexOf("foo", tasks)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, ErrUnknownId) {
		t.Errorf("expcted error %v, got %v", ErrUnknownId, err)
	}
}

func Test_BasicParse(t *testing.T) {
	input := []parser.Task {
		{Id: "id1", Cmd: "cmd1", Deps: []string{}},
		{Id: "id2", Cmd: "cmd2", Deps: []string{"id1"}},
		{Id: "id3", Cmd: "cmd3", Deps: []string{}},
		{Id: "id4", Cmd: "cmd4", Deps: []string{"id1", "id2"}},
		{Id: "id5", Cmd: "cmd5", Deps: []string{"id3"}},
	}
	output, err := Parse(input)

	if err != nil {
		t.Fatalf("got unexpected error %v", err)
	}

	// Check we have the correct number of tasks.
	if len(output) != 5 {
		t.Fatalf("expected 5 tasks, got %d", len(output))
	}

	// Check that each task has been correctly created, ignoring dependencies.
	for i, task := range output {
		if task.Id != input[i].Id || task.Cmd != input[i].Cmd {
			t.Errorf("expected task with id: %s and cmd: %s, but got %q and %q",
						input[i].Id, input[i].Cmd, task.Id, task.Cmd)
		}
	}

	// Check that the dependencies are correct.
	if len(output[0].Deps) != 0 {
		t.Fatalf("expected %d deps, got %d", 0, len(output[0].Deps))
	}
	if len(output[1].Deps) != 1 {
		t.Fatalf("expected %d deps, got %d", 1, len(output[1].Deps))
	}
	if len(output[2].Deps) != 0 {
		t.Fatalf("expected %d deps, got %d", 0, len(output[2].Deps))
	}
	if len(output[3].Deps) != 2 {
		t.Fatalf("expected %d deps, got %d", 2, len(output[3].Deps))
	}
	if len(output[4].Deps) != 1 {
		t.Fatalf("expected %d deps, got %d", 1, len(output[4].Deps))
	}

	if output[1].Deps[0] != &output[0] {
		t.Errorf("incorrect dependency")
	}
	if (output[3].Deps[0] != &output[0]) || (output[3].Deps[1] != &output[1]) {
		t.Errorf("incorrect dependency")
	}
	if output[4].Deps[0] != &output[2] {
		t.Errorf("incorrect dependency")
	}
}

func Test_UnknownDep(t *testing.T) {
	input := []parser.Task {
		{Id: "id1", Cmd: "cmd1", Deps: []string{}},
		{Id: "id2", Cmd: "cmd2", Deps: []string{"id1"}},
		{Id: "id3", Cmd: "cmd3", Deps: []string{}},
		{Id: "id4", Cmd: "cmd4", Deps: []string{"id1", "id2"}},
		{Id: "id5", Cmd: "cmd5", Deps: []string{"id7"}},
	}
	_, err := Parse(input)

	if err == nil {
		t.Fatalf("expected error but got nil")
	}
	if !errors.Is(err, ErrUnknownDep) {
		t.Errorf("expected error %v, but got %v", ErrUnknownDep, err)
	}
}

func Test_CyclicDep(t *testing.T) {
}

