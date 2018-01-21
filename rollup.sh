#!/bin/sh

ME=`basename $0`

note() {
    echo ">>> $ME: NOTE: $1"
}

error() {
    echo ">>> $ME: ERROR: $1"
    exit 1
}

HERE=`basename $PWD`
if [ "$HERE" != "anwork_testing" ]; then
    error "this script must be run from the anwork_testing directory"
fi

note "updating anwork submodule"
git -C submodules/anwork pull
if [ "$?" -ne 0 ]; then
    error "failed to update anwork submodule"
fi

note "creating package"
./submodules/anwork/ci/package.sh
if [ "$?" -ne 0 ]; then
    error "failed to run package script"
fi

version="$(awk '/const Version =/ {print $NF}' submodules/anwork/cmd/anwork/command/command.go)"
note "using version $version"

note "copying package"
dir="release/v$version"
if [ ! -d "$dir" ]; then
    mkdir "$dir"
fi
mv submodules/anwork/anwork-$version.zip "$dir"

note "commiting"
hash="$(git -C submodules/anwork log -1 --oneline | awk '{print $1}')"
git commit -a -m "Rollup anwork to $hash (version $version)."

