#!/bin/sh

ME=`basename $0`

note() {
  echo "$ME: note: $1  "
}

note "GOPATH is: '$GOPATH'"
previous=$GOPATH

export set GOPATH="$PWD"
note "Set GOPATH to: '$GOPATH'"

note "Kicking off tests"
go test core v1 -v -p 2
error=$?

export set GOPATH="$previous"
note "Reset GOPATH to: '$GOPATH'"

exit $error