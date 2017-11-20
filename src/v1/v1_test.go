package v1

import (
	"core"
	"testing"
)

func getAnwork(t *testing.T) *core.Anwork {
	anwork, err := core.MakeAnwork(1) // version 1
	if err != nil {
		t.Fatal("Cannot get anwork:", err)
	}
	return anwork
}

func TestCreate(t *testing.T) {
	anwork := getAnwork(t)
	defer anwork.Close()
	defer anwork.Run("reset", "-f")

	expects := []core.Expect{
		core.Expect{anwork, []string{"reset", "-f"}, []string{}},

		// Create task-a and task-b. task-a is higher priority than task-b.
		core.Expect{anwork, []string{"task", "create", "task-a", "-p", "15"}, []string{}},
		core.Expect{anwork, []string{"task", "create", "task-b", "-p", "25"}, []string{}},

		// Set task-a running and then set it blocked.
		core.Expect{anwork, []string{"task", "set-running", "task-a"}, []string{}},
		core.Expect{anwork, []string{"task", "set-blocked", "task-a"}, []string{}},

		// Set task-b running and then set it finished.
		core.Expect{anwork, []string{"task", "set-running", "task-b"}, []string{}},
		core.Expect{anwork, []string{"task", "set-finished", "task-b"}, []string{}},

		// Set task-a running and then set it finished.
		core.Expect{anwork, []string{"task", "set-running", "task-a"}, []string{}},
		core.Expect{anwork, []string{"task", "set-finished", "task-a"}, []string{}},

		// Expect the summary.
		core.Expect{anwork,
			[]string{"task", "show", "-s"},
			[]string{"RUNNING.*",
				"BLOCKED.*",
				"WAITING.*",
				"FINISHED.*",
				".*task-a",
				".*task-b"}},
	}
	core.Run(t, expects...)
}
