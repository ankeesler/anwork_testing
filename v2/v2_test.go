package v2

import (
	"testing"

	"github.com/ankeesler/anwork_testing/core"
)

const (
	taskAName = "task-a"
	taskBName = "task-b"
)

var version int

func TestMain(m *testing.M) {
	core.RunTests(m, &version)
}

func getAnwork(t *testing.T) *core.Anwork {
	anwork, err := core.MakeAnwork(version)
	if err != nil {
		t.Fatal("Cannot get anwork:", err)
	}
	return anwork
}

func TestCreate(t *testing.T) {
	t.Parallel()

	anwork := getAnwork(t)
	defer anwork.Close()

	expects := []core.Expect{
		// Create task-a and task-b. task-a is higher priority than task-b.
		core.Expect{anwork,
			[]string{"create", taskAName},
			[]string{}},
		core.Expect{anwork,
			[]string{"create", taskBName},
			[]string{}},
		core.Expect{anwork,
			[]string{"show"},
			[]string{"RUNNING.*",
				"BLOCKED.*",
				"WAITING.*",
				".*" + taskAName + ".*",
				".*" + taskBName + ".*",
				"FINISHED.*"}},
	}
	core.Run(t, expects...)
}
