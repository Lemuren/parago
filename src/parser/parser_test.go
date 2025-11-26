package parser

import (
	"strings"
	"errors"
	"testing"
)

func Test_Shorthand(t *testing.T) {
	yaml := `
    build: echo Building
    test: echo Testing
    `
	r := strings.NewReader(yaml)

	// Parse the yaml. It should succeed since it's valid.
	tasks, err := Parse(r)
	if err != nil {
		t.Fatalf("correct yaml\n%q\nemitted error %q", yaml, err)
	}

	// Check we've parsed exactly two tasks.
	if len(tasks) != 2 {
		for _, task := range tasks {
			t.Errorf("%+v\n", task)
		}
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}

	// Check that the tasks are exactly what we expect.
	got := make(map[string]Task)
	for _, t := range tasks {
		got[t.Id] = t
	}

	build := got["build"]
	if build.Cmd != "echo Building" {
		t.Errorf("build task not parsed correctly: %+v", build)
	}

	test := got["test"]
	if test.Cmd != "echo Testing" {
		t.Errorf("test task not parsed correctly: %+v", test)
	}
}

func Test_FullDefinition(t *testing.T) {
	yaml := `
    build:
        cmd: echo Building
    test:
        cmd: echo Testing
        deps: [build]
    `
	r := strings.NewReader(yaml)

	// Parse the yaml. It should succeed since it's valid.
	tasks, err := Parse(r)
	if err != nil {
		t.Fatalf("correct yaml\n%q\nemitted error %q", yaml, err)
	}

	// Check we've parsed exactly two tasks.
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}

	// Check that the tasks are exactly what we expect.
	got := make(map[string]Task)
	for _, t := range tasks {
		got[t.Id] = t
	}

	build := got["build"]
	if build.Cmd != "echo Building" {
		t.Errorf("build task not parsed correctly: %+v", build)
	}

	test := got["test"]
	if test.Cmd != "echo Testing" {
		t.Errorf("test task not parsed correctly: %+v", test)
	}

	// Ensure the dependency got correctly parsed.
	if len(test.Deps) != 1 {
		t.Fatalf("expected 1 dependency, got %d", len(test.Deps))
	}
	if test.Deps[0] != "build" {
		t.Errorf("dependency of test task not parsed correctly: %s",
		test.Deps[0])
	}
}

func Test_ShorthandAndFull(t *testing.T) {
	yaml := `
    build: echo Building
    test:
        cmd: echo Testing
    `
	r := strings.NewReader(yaml)

	// Parse the yaml. It should succeed since it's valid.
	tasks, err := Parse(r)
	if err != nil {
		t.Fatalf("correct yaml\n%q\nemitted error %q", yaml, err)
	}

	// Check we've parsed exactly two tasks.
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}

	// Check that the tasks are exactly what we expect.
	got := make(map[string]Task)
	for _, t := range tasks {
		got[t.Id] = t
	}

	build := got["build"]
	if build.Cmd != "echo Building" {
		t.Errorf("build task not parsed correctly: %+v", build)
	}

	test := got["test"]
	if test.Cmd != "echo Testing" {
		t.Errorf("test task not parsed correctly: %+v", test)
	}
}

func Test_MissingCmd(t *testing.T) {
	yaml := `
    foo:
        hello: there
    `
	r := strings.NewReader(yaml)

	_, err := Parse(r)
	if err == nil {
		t.Errorf("expected error from invalid yaml\n%q\n, got nil", yaml)
	}
	if !errors.Is(err, ErrMissingCmd) {
		t.Errorf("expected %v, got %v", ErrMissingCmd, err)
	}
}

func Test_CmdNotString(t *testing.T) {
	yaml := `
    foo:
        cmd: 3
    `
	r := strings.NewReader(yaml)

	_, err := Parse(r)
	if err == nil {
		t.Errorf("expected error from invalid yaml\n%q\n, got nil", yaml)
	}
	if !errors.Is(err, ErrCmdNotString) {
		t.Errorf("expected %v, got %v", ErrCmdNotString, err)
	}
}

func Test_CmdNotStringShorthand(t *testing.T) {
	yaml := `
    foo: 3
    `
	r := strings.NewReader(yaml)

	_, err := Parse(r)
	if err == nil {
		t.Errorf("expected error from invalid yaml\n%q\n, got nil", yaml)
	}
	if !errors.Is(err, ErrInvalidDef) {
		t.Errorf("expected %v, got %v", ErrInvalidDef, err)
	}
}

func Test_InvalidDepsNoArray(t *testing.T) {
	yaml := `
    foo:
        cmd: echo Hello
        deps: seven
    `
	r := strings.NewReader(yaml)

	_, err := Parse(r)
	if err == nil {
		t.Errorf("expected error from invalid yaml\n%q\n, got nil", yaml)
	}
	if !errors.Is(err, ErrDepsNotArray) {
		t.Errorf("expected %v, got %v", ErrDepsNotArray, err)
	}
}

func Test_InvalidDepsArrayNotString(t *testing.T) {
	yaml := `
    foo:
        cmd: echo Hello
        deps: [this, is, a, string, 7, but, seven, is, there]
    `
	r := strings.NewReader(yaml)

	_, err := Parse(r)
	if err == nil {
		t.Errorf("expected error from invalid yaml\n%q\n, got nil", yaml)
	}
	if !errors.Is(err, ErrDepsNotArrayOfStrings) {
		t.Errorf("expected %v, got %v", ErrDepsNotArrayOfStrings, err)
	}
}

func Test_InvalidYaml(t *testing.T) {
	yaml := `
    this is not valid yaml
    `
	r := strings.NewReader(yaml)
	// Parse the yaml. It should fail since it's invalid.
	_, err := Parse(r)
	if err == nil {
		t.Errorf("expected error from invalid yaml\n%q\n, got nil", yaml)
	}
}

