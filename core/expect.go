package core

import (
	"errors"
	"fmt"
	"path"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

// This structure represents a command passed to an Anwork instance and a number of expected regular
// expressions to be matched against the data printed to stdout by running the Anwork instance. This
// struct is meant to be initialized manually, like this.
//   anwork := core.MakeAnwork(1) // version
//   expect := Expect{Anwork: anwork, Command: "foo", Regexes: []string{".*bar.*", "^bat$"}}
type Expect struct {
	// This is the Anwork instance that this expect will use to run.
	Anwork *Anwork

	// This is the command that will be run on the Anwork instance.
	Command []string

	// These are the regular expressions that the Expect struct will try to match against the output
	// lines from the Anwork field. If the length of this slice is 0, then no output lines will be
	// matched.
	Regexes []string
}

// This function does the running of an Expect instance. The expect.Command will be run via the
// expect.Anwork.Run method and the expect.Regexes will be matched against the expect.Anwork.Run
// output lines. This method will return a slice of strings that represent the n lines that were
// successfully matched against n expect.Regexs. If a expect.Regex is not found in the output lines,
// then the length of the returned slice will not match the length of the passed expect.Regex slice.
func (expect *Expect) Run(t *testing.T) ([]string, error) {
	if expect.Anwork == nil {
		return nil, errors.New(fmt.Sprintf("Invalid expect struct: %#v", expect))
	}

	output, err := expect.Anwork.Run(expect.Command...)
	if err != nil {
		return nil, err
	}
	t.Logf("Got output from '%s' command:\n%s", expect.Command, output)

	outputLines := makeOutputLines(output)

	matchedLines, err := getMatchedLines(outputLines, expect.Regexes)
	if err != nil {
		return nil, err
	}
	t.Logf("Matched lines '%s' from regexes '%s'", matchedLines, expect.Regexes)

	return matchedLines, nil
}

// This is a helper method to run a bunch of Expect structs and log the errors to a testing.T struct.
func Run(t *testing.T, expects ...Expect) {
	for _, expect := range expects {
		matched, err := expect.Run(t)
		callerStr := getCallerStr()
		if err != nil {
			t.Errorf("%s: Got error when running Expect struct %s: %s", callerStr, expect, err)
		} else if len(matched) != len(expect.Regexes) {
			t.Errorf("%s: Did not match regex '%s' when running Expect struct %s",
				callerStr, expect.Regexes[len(matched)], expect)
		}
	}
}

func getCallerStr() string {
	_, file, line, ok := runtime.Caller(2) // we want the caller of the caller of this function
	if !ok {
		file = "?"
		line = 0
	} else {
		file = path.Base(file)
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func makeOutputLines(output string) []string {
	lines := strings.Split(output, "\n")
	if len(lines) > 0 && len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func getMatchedLines(outputLines []string, regexes []string) ([]string, error) {
	matchedLines := make([]string, 0, len(regexes))
	var regexpErr error = nil

	for lineI, regexI := 0, 0; regexpErr == nil && lineI < len(outputLines) && regexI < len(regexes); lineI++ {
		line := outputLines[lineI]
		regex := regexes[regexI]
		if matches, err := regexp.Match(regex, []byte(line)); err != nil {
			regexpErr = err
		} else if matches {
			matchedLines = append(matchedLines, line)
			regexI++
		}
	}

	return matchedLines, regexpErr
}
