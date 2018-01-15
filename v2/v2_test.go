package v2

import (
	"fmt"
	"testing"
	"time"

	"github.com/ankeesler/anwork_testing/core"
)

const (
	taskAName  = "task-a"
	taskANote0 = "Note a 0"
	taskANote1 = "Note a 1"
	taskBName  = "task-b"
	taskBNote0 = "Note b 0"
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
		// Create task-a and task-b.
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

	// Make sure that the tasks are shown with the correct stuff.
	expects = []core.Expect{
		core.Expect{anwork,
			[]string{"show", taskAName},
			[]string{"Name: " + taskAName, "ID: 0", "State: WAITING"}},
		core.Expect{anwork,
			[]string{"show", taskBName},
			[]string{"Name: " + taskBName, "ID: 1", "State: WAITING"}},
	}
	core.Run(t, expects...)
}

func TestPriority(t *testing.T) {
	t.Parallel()

	anwork := getAnwork(t)
	defer anwork.Close()

	// Create task-a and task-b, give task-b a priority higer than task-a.
	expects := []core.Expect{
		core.Expect{anwork, []string{"create", taskAName}, []string{}},
		core.Expect{anwork, []string{"create", taskBName}, []string{}},
		core.Expect{anwork, []string{"set-priority", taskAName, "15"}, []string{}},
		core.Expect{anwork, []string{"set-priority", taskBName, "5"}, []string{}},
		core.Expect{anwork,
			[]string{"show"},
			[]string{"RUNNING.*",
				"BLOCKED.*",
				"WAITING.*",
				".*" + taskBName + ".*",
				".*" + taskAName + ".*",
				"FINISHED.*"}},
	}
	core.Run(t, expects...)

	// Change task-a to have a higher priority than task-b.
	expects = []core.Expect{
		core.Expect{anwork, []string{"set-priority", taskAName, "10"}, []string{}},
		core.Expect{anwork, []string{"set-priority", taskBName, "20"}, []string{}},
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

	// Make sure the individual journals for these tasks reflect the priority changes.
	expects = []core.Expect{
		core.Expect{anwork,
			[]string{"journal", taskAName},
			[]string{".*priority.*" + taskAName + ".*to 10",
				".*priority.*" + taskAName + ".*to 15",
				".*Created.*" + taskAName + ".*",
			}},

		core.Expect{anwork,
			[]string{"journal", taskBName},
			[]string{".*priority.*" + taskBName + ".*to 20",
				".*priority.*" + taskBName + ".*to 5",
				".*Created.*" + taskBName + ".*",
			}},
	}
	core.Run(t, expects...)

	// Check that the combined journal displays these 6 individual entries.
	expects = []core.Expect{
		core.Expect{anwork,
			[]string{"journal"},
			[]string{".*", ".*", ".*", ".*", ".*", ".*"}},
	}
	core.Run(t, expects...)

	// Make sure that the tasks are shown with the correct priority.
	expects = []core.Expect{
		core.Expect{anwork,
			[]string{"show", taskAName},
			[]string{"Name: " + taskAName, "Priority: 10"}},
		core.Expect{anwork,
			[]string{"show", taskBName},
			[]string{"Name: " + taskBName, "Priority: 20"}},
	}
	core.Run(t, expects...)
}

func TestState(t *testing.T) {
	t.Parallel()

	anwork := getAnwork(t)
	defer anwork.Close()

	// Create 2 tasks and set them to different states.
	expects := []core.Expect{
		core.Expect{anwork, []string{"create", taskAName}, []string{}},
		core.Expect{anwork, []string{"create", taskBName}, []string{}},
		core.Expect{anwork, []string{"set-running", taskAName}, []string{}},
		core.Expect{anwork, []string{"set-blocked", taskBName}, []string{}},
		core.Expect{anwork,
			[]string{"show"},
			[]string{"RUNNING.*",
				".*" + taskAName + ".*",
				"BLOCKED.*",
				".*" + taskBName + ".*",
				"WAITING.*",
				"FINISHED.*"}},
	}
	core.Run(t, expects...)

	// Set the tasks to new states.
	expects = []core.Expect{
		core.Expect{anwork, []string{"set-running", taskBName}, []string{}},
		core.Expect{anwork, []string{"set-finished", taskAName}, []string{}},
		core.Expect{anwork,
			[]string{"show"},
			[]string{"RUNNING.*",
				".*" + taskBName + ".*",
				"BLOCKED.*",
				"WAITING.*",
				"FINISHED.*",
				".*" + taskAName + ".*"}},
	}
	core.Run(t, expects...)

	// Make sure the individual journals for these tasks reflect the state changes.
	expects = []core.Expect{
		core.Expect{anwork,
			[]string{"journal", taskAName},
			[]string{".*state.*" + taskAName + ".*to Finished",
				".*state.*" + taskAName + ".*to Running",
				".*Created.*" + taskAName + ".*",
			}},

		core.Expect{anwork,
			[]string{"journal", taskBName},
			[]string{".*state.*" + taskBName + ".*to Running",
				".*state.*" + taskBName + ".*to Blocked",
				".*Created.*" + taskBName + ".*",
			}},
	}
	core.Run(t, expects...)

	// Check that the combined journal displays these 6 individual entries.
	expects = []core.Expect{
		core.Expect{anwork,
			[]string{"journal"},
			[]string{".*", ".*", ".*", ".*", ".*", ".*"}},
	}
	core.Run(t, expects...)

	// Make sure that the tasks are shown with the correct states.
	expects = []core.Expect{
		core.Expect{anwork,
			[]string{"show", taskAName},
			[]string{"Name: " + taskAName, "State: FINISHED"}},
		core.Expect{anwork,
			[]string{"show", taskBName},
			[]string{"Name: " + taskBName, "State: RUNNING"}},
	}
	core.Run(t, expects...)
}

func TestNote(t *testing.T) {
	t.Parallel()

	anwork := getAnwork(t)
	defer anwork.Close()

	// Create 2 tasks and add a note to one of them.
	expects := []core.Expect{
		core.Expect{anwork, []string{"create", taskAName}, []string{}},
		core.Expect{anwork, []string{"create", taskBName}, []string{}},
		core.Expect{anwork, []string{"note", taskAName, taskANote0}, []string{}},
		core.Expect{anwork,
			[]string{"journal", taskAName},
			[]string{".*" + taskANote0 + ".*", ".*Created.*" + taskAName + ".*"}},
	}
	core.Run(t, expects...)

	// Add a note to the other task.
	expects = []core.Expect{
		core.Expect{anwork, []string{"note", taskBName, taskBNote0}, []string{}},
		core.Expect{anwork,
			[]string{"journal", taskAName},
			[]string{".*" + taskANote0 + ".*", ".*Created.*" + taskAName + ".*"}},
		core.Expect{anwork,
			[]string{"journal", taskBName},
			[]string{".*" + taskBNote0 + ".*", ".*Created.*" + taskBName + ".*"}},
	}
	core.Run(t, expects...)

	// Add a second note to the first task.
	expects = []core.Expect{
		core.Expect{anwork, []string{"note", taskAName, taskANote1}, []string{}},
		core.Expect{anwork,
			[]string{"journal", taskAName},
			[]string{".*" + taskANote1 + ".*",
				".*" + taskANote0 + ".*",
				".*Created.*" + taskAName + ".*"}},
		core.Expect{anwork,
			[]string{"journal", taskBName},
			[]string{".*" + taskBNote0 + ".*",
				".*Created.*" + taskBName + ".*"}},
	}
	core.Run(t, expects...)

	// Check that the combined journal displays these 5 individual entries.
	expects = []core.Expect{
		core.Expect{anwork,
			[]string{"journal"},
			[]string{".*", ".*", ".*", ".*", ".*"}},
	}
	core.Run(t, expects...)
}

func TestDelete(t *testing.T) {
	t.Parallel()

	anwork := getAnwork(t)
	defer anwork.Close()

	// Create 2 tasks and delete one of them.
	expects := []core.Expect{
		core.Expect{anwork, []string{"create", taskAName}, []string{}},
		core.Expect{anwork, []string{"create", taskBName}, []string{}},
		core.Expect{anwork, []string{"delete", taskAName}, []string{}},
		core.Expect{anwork,
			[]string{"journal"},
			[]string{".*Deleted.*" + taskAName + ".*",
				".*Created.*" + taskBName + ".*",
				".*Created.*" + taskAName + ".*",
			}},
		core.Expect{anwork,
			[]string{"show"},
			[]string{"RUNNING.*",
				"BLOCKED.*",
				"WAITING.*",
				".*" + taskBName + ".*",
				"FINISHED.*"}},
	}
	core.Run(t, expects...)

	// Delete the other task.
	expects = []core.Expect{
		core.Expect{anwork, []string{"delete", taskBName}, []string{}},
		core.Expect{anwork,
			[]string{"journal"},
			[]string{".*Deleted.*" + taskBName + ".*",
				".*Deleted.*" + taskAName + ".*",
				".*Created.*" + taskBName + ".*",
				".*Created.*" + taskAName + ".*",
			}},
		core.Expect{anwork,
			[]string{"show"},
			[]string{"RUNNING.*",
				"BLOCKED.*",
				"WAITING.*",
				"FINISHED.*"}},
	}
	core.Run(t, expects...)
}

func TestDeleteAll(t *testing.T) {
	t.Parallel()

	anwork := getAnwork(t)
	defer anwork.Close()

	// Create 2 tasks and delete one of them.
	expects := []core.Expect{
		core.Expect{anwork, []string{"create", taskAName}, []string{}},
		core.Expect{anwork, []string{"create", taskBName}, []string{}},
		core.Expect{anwork, []string{"delete", taskAName}, []string{}},
		core.Expect{anwork,
			[]string{"journal"},
			[]string{".*Deleted.*" + taskAName + ".*",
				".*Created.*" + taskBName + ".*",
				".*Created.*" + taskAName + ".*",
			}},
		core.Expect{anwork,
			[]string{"show"},
			[]string{"RUNNING.*",
				"BLOCKED.*",
				"WAITING.*",
				".*" + taskBName + ".*",
				"FINISHED.*"}},
	}
	core.Run(t, expects...)

	// Delete all of the tasks. Make sure they are gone.
	expects = []core.Expect{
		core.Expect{anwork, []string{"delete-all"}, []string{}},
		core.Expect{anwork,
			[]string{"journal"},
			[]string{".*Deleted.*", ".*Deleted.*", ".*Created.*", ".*Created.*"}},
		core.Expect{anwork,
			[]string{"show"},
			[]string{"RUNNING.*",
				"BLOCKED.*",
				"WAITING.*",
				"FINISHED.*"}},
	}
	core.Run(t, expects...)
}

func TestReset(t *testing.T) {
	t.Parallel()

	anwork := getAnwork(t)
	defer anwork.Close()

	// Create 2 tasks.
	expects := []core.Expect{
		core.Expect{anwork, []string{"create", taskAName}, []string{}},
		core.Expect{anwork, []string{"create", taskBName}, []string{}},
	}
	core.Run(t, expects...)

	// Reset everything. Make sure it is gone.
	expects = []core.Expect{
		core.Expect{anwork, []string{"reset", "y"}, []string{}},
		core.Expect{anwork,
			[]string{"journal"},
			[]string{}},
		core.Expect{anwork,
			[]string{"show"},
			[]string{"RUNNING.*",
				"BLOCKED.*",
				"WAITING.*",
				"FINISHED.*"}},
	}
	core.Run(t, expects...)
}

func TestSummary(t *testing.T) {
	t.Parallel()

	anwork := getAnwork(t)
	defer anwork.Close()

	// Create 2 tasks.
	expects := []core.Expect{
		core.Expect{anwork, []string{"create", taskAName}, []string{}},
		core.Expect{anwork, []string{"create", taskBName}, []string{}},
	}
	core.Run(t, expects...)

	// Wait one second.
	time.Sleep(time.Second)

	// Set one of the tasks as finished. They should be reported in the summary.
	expects = []core.Expect{
		core.Expect{anwork, []string{"set-finished", taskAName}, []string{}},
		core.Expect{anwork,
			[]string{"summary", "1"},
			[]string{"\\[.*\\]:.*" + taskAName + ".*", "  took \\ds"}},
	}
	core.Run(t, expects...)
}

func TestIdUniqueness(t *testing.T) {
	t.Parallel()

	anwork := getAnwork(t)
	defer anwork.Close()

	// Create a task. Make sure it has id 0.
	expects := []core.Expect{
		core.Expect{anwork, []string{"create", taskAName}, []string{}},
		core.Expect{anwork, []string{"show", taskAName}, []string{"ID: 0"}},
	}
	core.Run(t, expects...)

	// Delete the task.
	expects = []core.Expect{
		core.Expect{anwork, []string{"delete", taskAName}, []string{}},
	}
	core.Run(t, expects...)

	// Create another task. Make sure it has id 1.
	expects = []core.Expect{
		core.Expect{anwork, []string{"create", taskBName}, []string{}},
		core.Expect{anwork, []string{"show", taskBName}, []string{"ID: 1"}},
	}
	core.Run(t, expects...)
}

func BenchmarkCreate(b *testing.B) {
	b.N = 5
	core.RunBenchmark(b, version, func(a *core.Anwork, i int) {
		name := fmt.Sprintf("task-%d", i)
		a.Run("create", name)
	})
}

func BenchmarkCrud(b *testing.B) {
	b.N = 5
	core.RunBenchmark(b, version, func(a *core.Anwork, i int) {
		name := fmt.Sprintf("task-%d", i)
		a.Run("create", name)
		a.Run("show")
		a.Run("set-finished", name)
		a.Run("delete", name)
	})
}
