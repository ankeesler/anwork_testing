package main

import (
	"os/exec"
	"testing"
)

var DEBUG bool = true

// This function runs the anwork binary with the args provided, and the returns both the string
// output and the error (or nil if there was no error).
func run(t *testing.T, args []string) (string, error) {
	path := "release/anwork/bin/anwork"
	cmd := exec.Command(path, args...)
	if DEBUG {
		t.Log("running command", cmd.Args)
	}
	output, err := cmd.Output()
	return string(output), err
}

func TestSanity(t *testing.T) {
	output, err := run(t, []string{"-d", "reset", "--force"})
	if err != nil {
		t.Fatal("Failed to run anwork binary", err)
	} else {
		t.Log("output of command was ", string(output))
	}

	output, err = run(t, []string{"-d", "task", "create", "task-a"})
	if err != nil {
		t.Error("Failed to run anwork binary", err)
	} else {
		t.Log("output of command was ", string(output))
	}
}
