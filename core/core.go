// This package contains utilities for running Anwork tests.
//
// This core package makes 3 different contributions to a test package.
//
// 1. The RunTests function MUST be called from a TestMain function in the test package that wants to
// use this test framework. Here is an example of a test that uses this test framework.
//
//   var version int // global variable
//   ...
//   func TestMain(m *testing.M) {
//     core.RunTests(m, &version)
//   }
//   ...
//   func TestSomething(t testing.T) {
//     anwork, err := core.MakeAnwork(version)
//     ...
//   }
//
// 2. The Anwork type should be used to abstract the zipped anwork releases found in this repo. See
// the Anwork type and MakeAnwork for more information.
//
// 3. The Expect type should be used to pass commands to an Anwork instance and assert that responses
// were printed out from the executable.
package core

import (
	"flag"
	"os"
	"testing"
)

func RunTests(m *testing.M, version *int) {
	flag.IntVar(version, "v", 0, "The anwork version that should be used with these tests")
	flag.Parse()

	if *version == 0 {
		panic("Version (-v) must be passed with a legitimate anwork version number")
	}

	os.Exit(m.Run())
}
