# ANWORK TESTING

This repo contains the test infrastructure and tests for the ANWORK project.

[![Build Status](https://travis-ci.org/ankeesler/anwork_testing.svg?branch=master)](https://travis-ci.org/ankeesler/anwork_testing)

See [anwork](https://github.com/ankeesler/anwork) repo for source code to the ANWORK project.

## Running Tests

Tests and release packages are organized by release version. The idea is that if you want to run
the tests associated with the V_x_ release, you should also run the tests from the V_i_ releases,
where _i_ ranges from 1 to _x_, inclusive. Whenever you run the V_x_ tests, the V_x_ package will
be used for all tests.

The following command will run the V_x_ test where _x_ is the latest release.
```
$ ./run-tests.sh
```

## Directory Structure

```
release/
  v1/
    anwork.zip # V1 release
  v2/
    anwork.zip # V2 release
  ...
  latest/
    anwork.zip # Latest successfull build
src/
  v1/
    v1_test.go # Tests related to V1 release
  v2/
    v2_test.go # Tests related to V2 release
  ...
```