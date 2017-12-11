package core

import (
	"testing"
)

func TestBadExpectStruct(t *testing.T) {
	t.Parallel()

	anwork := mustGetAnwork(t)
	defer anwork.Close()

	bads := []Expect{
		Expect{Anwork: nil, Command: []string{"foo"}, Regexes: make([]string, 2)},
		Expect{Anwork: &Anwork{}, Command: []string{"foo"}, Regexes: make([]string, 2)},
	}

	for _, bad := range bads {
		if _, err := bad.Run(t); err == nil {
			t.Error("We should have returned an error for an invalid Expect struct!")
		} else {
			t.Logf("Successfully received error from invalid Expect struct: %s", err)
		}
	}
}

func TestIgnoreOutput(t *testing.T) {
	t.Parallel()

	anwork := mustGetAnwork(t)
	defer anwork.Close()

	expects := []Expect{
		Expect{anwork, []string{"reset", "-f"}, []string{}},         // don't care about output
		Expect{anwork, []string{"version"}, []string{}},             // don't care about output
		Expect{anwork, []string{"task", "create", "a"}, []string{}}, // don't care about output
		Expect{anwork, []string{"task", "create", "b"}, []string{}}, // don't care about output
		Expect{anwork, []string{"task", "show"}, []string{}},        // don't care about output
	}
	for _, expect := range expects {
		matched, err := expect.Run(t)
		if err != nil {
			t.Errorf("Received fatal error when trying to run expect %s: %s", expect, err)
		} else if len(matched) > 0 {
			t.Errorf("Received unexpected matched lines when running expect struct: %s", expect)
		}
	}
}

func TestMatchOutput(t *testing.T) {
	t.Parallel()

	anwork := mustGetAnwork(t)
	defer anwork.Close()

	expects := []Expect{
		Expect{anwork, []string{"reset", "-f"}, []string{}},
		Expect{anwork, []string{"version"}, []string{"Version = 1"}},
		Expect{anwork, []string{"-d", "task", "create", "task-a"}, []string{"debug:.*created.*task-a"}},
		Expect{anwork, []string{"-d", "task", "create", "task-b"}, []string{"debug:.*created.*task-b"}},
		Expect{anwork, []string{"task", "show"}, []string{"WAITING.*", ".*task-a"}},
	}
	Run(t, expects...)
}

func TestMakeOutputLines(t *testing.T) {
	t.Parallel()

	data := []struct {
		output string
		lines  []string
	}{
		{"", []string{}},
		{"one line", []string{"one line"}},
		{"two\nlines", []string{"two", "lines"}},
		{"many and\nmany\nand many\nlines\n", []string{"many and", "many", "and many", "lines"}},
	}

	for _, datum := range data {
		lines := makeOutputLines(datum.output)
		if !areSlicesEqual(datum.lines, lines) {
			t.Errorf("Wanted %s lines from '%s', got %s", datum.lines, datum.output, lines)
		}
	}
}

func TestGetMatchedLines(t *testing.T) {
	t.Parallel()

	data := []struct {
		lines   []string
		regexes []string
		matched []string
	}{
		// Matching nothing.
		{[]string{}, []string{}, []string{}},
		{[]string{}, []string{".*"}, []string{}},
		{[]string{}, []string{".*", "."}, []string{}},

		// Nothing matching.
		{[]string{"hey", "there"}, []string{}, []string{}},
		{[]string{"hey", "there"}, []string{".*tuna.*"}, []string{}},
		{[]string{"hey", "there"}, []string{".*fish.*", ".*marlin.*"}, []string{}},
		{[]string{"hey", "there"}, []string{"hye", "^ey$"}, []string{}},
		{[]string{"hey", "there"}, []string{"nope", ".*"}, []string{}},

		// Some things matching.
		{[]string{"hey", "there"}, []string{".*ey$", "nope"}, []string{"hey"}},
		{[]string{"hey", "there"}, []string{".*ere$"}, []string{"there"}},

		// Too many regexes.
		{[]string{"hey", "there"}, []string{".*", ".*", ".*"}, []string{"hey", "there"}},

		// Everything matching.
		{[]string{"hey", "there"}, []string{"hey", ".*"}, []string{"hey", "there"}},
		{[]string{"hey", "there", "foo"}, []string{"hey", "foo"}, []string{"hey", "foo"}},
	}

	for _, datum := range data {
		matched, err := getMatchedLines(datum.lines, datum.regexes)
		if err != nil {
			t.Errorf("Got error from matching %s lines against %s regexes: %s",
				datum.lines, datum.regexes, err)
		} else if !areSlicesEqual(datum.matched, matched) {
			t.Errorf("Wanted %s output from %s lines and %s regexes, but got %s",
				datum.matched, datum.lines, datum.regexes, matched)
		}
	}
}

func TestBadRegex(t *testing.T) {
	t.Parallel()

	anwork := mustGetAnwork(t)
	defer anwork.Close()
	expect := Expect{anwork, []string{"version"}, []string{"["}} // regex is missing a closing ']'
	_, err := expect.Run(t)
	if err == nil {
		t.Error("Expected an error from bad regex!")
	}
}

func areSlicesEqual(slice1, slice2 []string) bool {
	if slice1 == nil {
		return slice2 == nil
	} else if slice2 == nil {
		return false
	}

	if len(slice1) != len(slice2) {
		return false
	}

	for index := range slice1 {
		if slice1[index] != slice2[index] {
			return false
		}
	}

	return true
}

func mustGetAnwork(t *testing.T) *Anwork {
	anwork, err := MakeAnwork(1) // version
	if err != nil {
		t.Fatalf("Received fatal error when trying to make anwork struct: %s", err)
	}
	return anwork
}
