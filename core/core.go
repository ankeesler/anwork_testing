// This package contains utilities for running Anwork tests.
//
// This core package makes 4 different contributions to a test package.
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
//
// 4. The RunBenchmark function is a utility provided to test packages for benchmarking. See function
// for further details.
package core

import (
	"flag"
	"os"
	"testing"
)

// This function MUST be called from a TestMain function inside the test package that wants to use
// this test framework. This function parses a version argument passed to the test executable. If no
// version argument is passed (via the -v flag), then this function will panic.
func RunTests(m *testing.M, version *int) {
	flag.IntVar(version, "v", 0, "The anwork version that should be used with these tests")
	flag.Parse()

	if *version == 0 {
		panic("Version (-v) must be passed with a legitimate anwork version number")
	}

	os.Exit(m.Run())
}

// This function allocates an Anwork struct with the provided version and then runs the provided
// function b.N number of times. It resets the b timer (with b.ResetTimer()) right before it runs
// the function. The integer argument to the function is the number of benchmark iteration that is
// being run.
func RunBenchmark(b *testing.B, version int, f func(*Anwork, int)) {
	a, err := MakeAnwork(version)
	if err != nil {
		b.Fatal(err)
	}
	defer a.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f(a, i)
	}
}
