package v1

import (
	"testing"
	"util"
)

func getAnwork(t *testing.T) *util.Anwork {
	anwork, err := util.MakeAnwork(1) // version 1
	if err != nil {
		t.Fatal("Cannot get anwork:", err)
	}
	return anwork
}

func TestCreate(t *testing.T) {
	anwork := getAnwork(t)
	defer anwork.Close()
	output, err := anwork.Run("reset", "--force")
	if err != nil {
		t.Fatal("Failed to reset context:", err, output)
	}

	output, err = anwork.Run("task", "create", "task-a")
	if err != nil {
		t.Fatal("Failed to create task:", err, output)
	}
}
