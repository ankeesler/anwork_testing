#!/bin/sh

# This script is the interface to running anwork tests.

ME=`basename $0`

usage() {
    echo "usage: $ME [-v X] [-t X] [-l] [-a]"
    echo
    echo "-a     Run all tests for every version"
    echo "-l     Run all the tests with the latest version"
    echo "-t X   Only run the tests in package vX"
    echo "-v X   Run the tests for version X"
    echo
    echo "If no flags are passed, -a is assumed"
    echo
    echo "Example: $ME -v 15      # Run all tests with version 15"
    echo "Example: $ME -v 15 -t 1 # Run tests in v1 package with version 15"
    echo "Example: $ME -t 1       # Run test v1 for all versions"
    echo "Example: $ME -a         # Run all tests for all versions"
    echo "Example: $ME -a -v 15   # Ignore version flag and run all tests for all versions"
    echo "Example: $ME -a -t 15   # Ignore test flag and run all tests for all versions"
    echo "Example: $ME -l         # Run all tests for latest version"
    echo "Example: $ME -l -v 15   # Ignore version flag and run all tests for latest version"
    echo "Example: $ME -l -t 15   # Ignore test flag and run all tests for latest version"
}

note() {
    echo ">>> $ME: NOTE: $1"
}

error() {
    echo ">>> $ME: ERROR: $1"
    exit 1
}

runtest() {
    command="go test github.com/ankeesler/anwork_testing/$1 -args -v $2"
    note "running command: $command"
    output="$($command)"
    note "output:\n$output"
}

all=0
latest=0
tehst=
version=
while getopts alt:v: o
do
    case "$o" in
        a)   all=1;;
        l)   latest=1;;
        t)   tehst="$OPTARG";;
        v)   version="$OPTARG";;
        [?]) usage && exit 1;;
    esac
done

if [ "$all" -eq 0 ] && [ "$latest" -eq 0 ] && [ -z "$tehst" ] && [ -z "$version" ]; then
    note "no flags provided, assuming -a"
    all=1
fi

if [ "$all" -eq 1 ]; then
    error "implement -a flag behavior!!!"
fi

if [ "$latest" -eq 1 ]; then
    error "implement -l flag behavior!!!"
fi

if [ ! -z "$version" ] && [ ! -d "release/v$version" ]; then
    error "unknown version: $version"
fi

if [ ! -z "$tehst" ] && [ ! -d "v$tehst" ]; then
    error "unknown test: $tehst"
fi

if [ -z "$tehst" ]; then
    for dir in ./v*; do
        runtest "$(basename $dir)" "$version"
    done
else
    runtest "v$tehst" "$version"
fi
