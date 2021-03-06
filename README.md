# ANWORK TESTING

NOTE: This repo has been deprecated now that there are integration tests in the anwork repo.
Check them out here: https://github.com/ankeesler/anwork/tree/master/integration.

This repo contains the test infrastructure and tests for the ANWORK project.

[![Build Status](https://travis-ci.org/ankeesler/anwork_testing.svg?branch=master)](https://travis-ci.org/ankeesler/anwork_testing)

See [anwork](https://github.com/ankeesler/anwork) repo for source code to the ANWORK project.

## Running Tests

Tests and release packages are organized by release version. The idea is that if you want to run
the tests associated with the V*x* release, you should also run the tests from the V*i* releases,
where *i* ranges from 1 to *x*, inclusive. Whenever you run the V*x* tests, the V*x* package will
be used for all tests.

The following command will run the V*x* test where *x* is the latest release.
```
$ ./test.sh -v x
```

## Directory Structure

```
release/
  v1/
    anwork-1.zip # V1 release
  v2/
    anwork-2.zip # V2 release
  ...
core/        # Core test framework functionality
  data/      # Test data for core test framework tests
v1/
  data/      # Test data for V1 release tests
  v1_test.go # Tests related to V1 release
v2/
  data/      # Test data for V2 release tests
  v2_test.go # Tests related to V2 release
...
```
