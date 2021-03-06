package v1

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ankeesler/anwork_testing/core"
)

const (
	taskAName        = "task-a"
	taskADescription = "This is the description for " + taskAName
	taskAPriority    = "15"
	taskANote1       = "This is the first note for " + taskAName
	taskANote2       = "This is the second note for " + taskAName

	taskBName        = "task-b"
	taskBDescription = "This is the description for " + taskBName
	taskBPriority    = "25"
	taskBNote1       = "This is the first note for " + taskBName
	taskBNote2       = "This is the second note for " + taskBName

	taskCName        = "task-c"
	taskCDescription = "This is the description for " + taskCName
	taskCPriority    = "20"
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
	if version != 1 {
		t.Skipf("Skipping TestNote for non version 1 packages because of updated 'create' command flags")
	}

	anwork := getAnwork(t)
	defer anwork.Close()

	expects := []core.Expect{
		// Create task-a and task-b. task-a is higher priority than task-b.
		core.Expect{anwork,
			[]string{"task", "create", taskAName, "-p", taskAPriority, "--description", taskADescription},
			[]string{}},
		core.Expect{anwork,
			[]string{"task", "create", taskBName, "-p", taskBPriority, "--description", taskBDescription},
			[]string{}},
		core.Expect{anwork,
			[]string{"task", "show"},
			[]string{"RUNNING.*",
				"BLOCKED.*",
				"WAITING.*",
				".*" + taskAName + ".*",
				".*priority.*" + taskAPriority + ".*",
				//".*" + taskADescription + ".*", no description?
				".*" + taskBName + ".*",
				".*priority.*" + taskBPriority + ".*",
				//".*" + taskBDescription + ".*", no description?
				"FINISHED.*"}},
	}
	core.Run(t, expects...)
}

func TestSetState(t *testing.T) {
	t.Parallel()
	if version != 1 {
		t.Skipf("Skipping TestNote for non version 1 packages because of updated command names")
	}

	anwork := getAnwork(t)
	defer anwork.Close()

	expects := []core.Expect{
		// Create 3 tasks and set them all to different states.
		core.Expect{anwork, []string{"task", "create", taskAName}, []string{}},
		core.Expect{anwork, []string{"task", "create", taskBName}, []string{}},
		core.Expect{anwork, []string{"task", "create", taskCName}, []string{}},
		core.Expect{anwork, []string{"task", "set-running", taskAName}, []string{}},
		core.Expect{anwork, []string{"task", "set-blocked", taskBName}, []string{}},
		core.Expect{anwork, []string{"task", "set-waiting", taskCName}, []string{}},
		core.Expect{anwork,
			[]string{"task", "show", "-s"},
			[]string{"RUNNING.*",
				".*" + taskAName + ".*",
				"BLOCKED.*",
				".*" + taskBName + ".*",
				"WAITING.*",
				".*" + taskCName + ".*",
				"FINISHED.*"}},

		// Set the 3 tasks to new states.
		core.Expect{anwork, []string{"task", "set-blocked", taskAName}, []string{}},
		core.Expect{anwork, []string{"task", "set-running", taskBName}, []string{}},
		core.Expect{anwork, []string{"task", "set-finished", taskCName}, []string{}},
		core.Expect{anwork,
			[]string{"task", "show", "-s"},
			[]string{"RUNNING.*",
				".*" + taskBName + ".*",
				"BLOCKED.*",
				".*" + taskAName + ".*",
				"WAITING.*",
				"FINISHED.*",
				".*" + taskCName + ".*"}},
	}
	core.Run(t, expects...)
}

func TestChangePriority(t *testing.T) {
	t.Parallel()
	if version != 1 {
		t.Skipf("Skipping TestNote for non version 1 packages because of simplified printing format")
	}

	anwork := getAnwork(t)
	defer anwork.Close()

	expects := []core.Expect{
		// Create 2 tasks, change their priorities, and make sure they are what we expect.
		core.Expect{anwork, []string{"task", "create", taskAName, "-p", taskAPriority}, []string{}},
		core.Expect{anwork, []string{"task", "create", taskBName, "-p", taskBPriority}, []string{}},
		core.Expect{anwork, []string{"task", "set-priority", taskAName, taskBPriority}, []string{}},
		core.Expect{anwork, []string{"task", "set-priority", taskBName, taskAPriority}, []string{}},
		core.Expect{anwork,
			[]string{"task", "show"},
			[]string{"RUNNING.*",
				"BLOCKED.*",
				"WAITING.*",
				".*" + taskBName + ".*",
				".*priority.*" + taskAPriority + ".*",
				".*" + taskAName + ".*",
				".*priority.*" + taskBPriority + ".*",
				"FINISHED.*"}},
	}
	core.Run(t, expects...)
}

func TestNote(t *testing.T) {
	t.Parallel()
	if version != 1 {
		t.Skipf("Skipping TestNote for non version 1 packages because of simplified printing format")
	}

	anwork := getAnwork(t)
	defer anwork.Close()

	expects := []core.Expect{
		// Create 2 tasks and make sure the latest note is shown in "task show."
		core.Expect{anwork, []string{"task", "create", taskAName, "-p", taskAPriority}, []string{}},
		core.Expect{anwork, []string{"task", "create", taskBName, "-p", taskBPriority}, []string{}},
		core.Expect{anwork, []string{"task", "note", taskAName, taskANote1}, []string{}},
		core.Expect{anwork, []string{"task", "note", taskBName, taskBNote1}, []string{}},
		core.Expect{anwork,
			[]string{"task", "show"},
			[]string{"RUNNING.*",
				"BLOCKED.*",
				"WAITING.*",
				".*" + taskAName + ".*",
				".*" + taskANote1 + ".*",
				".*" + taskBName + ".*",
				".*" + taskBNote1 + ".*",
				"FINISHED.*"}},

		// Add different notes to the tasks.
		core.Expect{anwork, []string{"task", "note", taskAName, taskANote2}, []string{}},
		core.Expect{anwork, []string{"task", "note", taskBName, taskBNote2}, []string{}},
		core.Expect{anwork,
			[]string{"task", "show"},
			[]string{"RUNNING.*",
				"BLOCKED.*",
				"WAITING.*",
				".*" + taskAName + ".*",
				".*" + taskANote2 + ".*",
				".*" + taskBName + ".*",
				".*" + taskBNote2 + ".*",
				"FINISHED.*"}},
	}
	core.Run(t, expects...)
}

func TestJournal(t *testing.T) {
	t.Parallel()
	if version != 1 {
		t.Skipf("Skipping TestNote for non version 1 packages because of updated command names")
	}

	anwork := getAnwork(t)
	defer anwork.Close()

	expects := []core.Expect{
		// Create 2 tasks, add some notes, and set some states. There should be at least 4 journal entries
		// for each task.
		core.Expect{anwork, []string{"task", "create", taskAName}, []string{}},
		core.Expect{anwork, []string{"task", "create", taskBName}, []string{}},
		core.Expect{anwork, []string{"task", "note", taskAName, taskANote1}, []string{}},
		core.Expect{anwork, []string{"task", "note", taskBName, taskBNote1}, []string{}},
		core.Expect{anwork, []string{"task", "set-running", taskAName}, []string{}},
		core.Expect{anwork, []string{"task", "set-blocked", taskBName}, []string{}},
		core.Expect{anwork, []string{"task", "note", taskAName, taskANote2}, []string{}},
		core.Expect{anwork, []string{"task", "note", taskBName, taskBNote2}, []string{}},
		core.Expect{anwork, []string{"journal", "show", taskAName}, []string{".*", ".*", ".*", ".*"}},
		core.Expect{anwork, []string{"journal", "show", taskBName}, []string{".*", ".*", ".*", ".*"}},

		// The whole journal should contain at least 8 entries.
		core.Expect{anwork,
			[]string{"journal", "show-all"},
			[]string{".*", ".*", ".*", ".*", ".*", ".*", ".*", ".*"}},
	}
	core.Run(t, expects...)
}

func TestDelete(t *testing.T) {
	t.Parallel()
	if version != 1 {
		t.Skipf("Skipping TestNote for non version 1 packages because of updated command names")
	}

	anwork := getAnwork(t)
	defer anwork.Close()

	expects := []core.Expect{
		// Create 2 tasks and delete one of them.
		core.Expect{anwork, []string{"task", "create", taskAName}, []string{}},
		core.Expect{anwork, []string{"task", "create", taskBName}, []string{}},
		core.Expect{anwork, []string{"task", "delete", taskAName}, []string{}},
	}
	core.Run(t, expects...)

	// We should only see one of our tasks.
	expectDoesNotContain(t, anwork, taskAName)

	expects = []core.Expect{
		// Delete the remaining task.
		core.Expect{anwork, []string{"task", "delete", taskBName}, []string{}},
	}
	core.Run(t, expects...)

	// We should not see any tasks.
	expectDoesNotContain(t, anwork, taskBName)
}

func TestDeleteAll(t *testing.T) {
	t.Parallel()
	if version != 1 {
		t.Skipf("Skipping TestNote for non version 1 packages because of updated command names")
	}

	anwork := getAnwork(t)
	defer anwork.Close()

	expects := []core.Expect{
		// Create 2 tasks, and then delete them both at the same time.
		core.Expect{anwork, []string{"task", "create", taskAName}, []string{}},
		core.Expect{anwork, []string{"task", "create", taskBName}, []string{}},
		core.Expect{anwork, []string{"task", "delete-all"}, []string{}},
	}
	core.Run(t, expects...)
	expectDoesNotContain(t, anwork, taskAName)
	expectDoesNotContain(t, anwork, taskBName)
}

func expectDoesNotContain(t *testing.T, anwork *core.Anwork, taskName string) {
	output, err := anwork.Run("task", "show")
	if err != nil {
		t.Errorf("Command failed: %s", err)
	} else if strings.Contains(output, taskName) {
		t.Errorf("Didn't expect to see task name '%s' in output:\n%s", taskName, output)
	}
}

func BenchmarkCreate(b *testing.B) {
	if version != 1 {
		b.Skipf("Skipping BenchmarkCreate for non version 1 packages because of updated command names")
	}

	b.N = 5
	core.RunBenchmark(b, version, func(a *core.Anwork, i int) {
		name := fmt.Sprintf("task-%d", i)
		a.Run("task", "create", name)
	})
}

func BenchmarkCrud(b *testing.B) {
	if version != 1 {
		b.Skipf("Skipping BenchmarkCreate for non version 1 packages because of updated command names")
	}

	b.N = 5
	core.RunBenchmark(b, version, func(a *core.Anwork, i int) {
		name := fmt.Sprintf("task-%d", i)
		a.Run("task", "create", name)
		a.Run("task", "show")
		a.Run("task", "set-finished", name)
		a.Run("task", "delete", name)
	})
}
