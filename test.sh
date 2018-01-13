#!/bin/sh

# This script is the interface to running anwork tests.

ME=`basename $0`

usage() {
    echo "usage: $ME -v X [-b] [-t X] [-n]"
    echo
    echo "-b     Run all benchmarks as well as tests"
    echo "-n     Don't actually run the tests, only print the commands"
    echo "-t X   Only run the tests in package vX"
    echo "-v X   Run the tests for version X"
    echo
    echo "Example: $ME -v 15      # Run all tests with version 15"
    echo "Example: $ME -v 15 -t 1 # Run tests in v1 package with version 15"
}

note() {
    echo ">>> $ME: NOTE: $1"
}

error() {
    echo ">>> $ME: ERROR: $1"
    exit 1
}

runtest() {
    command="go test"
    if [ "$bench" -eq 1 ]; then
        command="$command -bench ."
    fi
    command="$command github.com/ankeesler/anwork_testing/$1 -args -v $2"
    if [ "$norun" -ne 1 ]; then
        note "running command: $command"
        echo "$($command)"
    else
        note "NOT running command: $command"
    fi
}

bench=0
norun=0
tehst=
version=
while getopts bnt:v: o
do
    case "$o" in
        b)   bench=1;;
        n)   norun=1;;
        t)   tehst="$OPTARG";;
        v)   version="$OPTARG";;
        [?]) usage && exit 1;;
    esac
done

if [ -z "$version" ]; then
    error "version not specified via -v flag; must specify version"
fi

if [ ! -z "$version" ] && [ ! -d "release/v$version" ]; then
    error "unknown version: $version"
fi

if [ ! -z "$tehst" ] && [ ! -d "v$tehst" ]; then
    error "unknown test: $tehst"
fi

if [ -z "$tehst" ]; then
    for dir in ./v*; do
        testversion="$(echo $(basename $dir) | sed -e 's/v//')"
        if [ "$version" -ge "$testversion" ]; then
            runtest "v$testversion" "$version"
        fi
    done
else
    runtest "v$tehst" "$version"
fi
