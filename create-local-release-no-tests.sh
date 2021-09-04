#!/bin/bash

set -e

# Set environment
export LANG="C.UTF-8"
export VERSION=$1

if [ -z "$VERSION" ]; then
	VERSION=$(git rev-parse --abbrev-ref HEAD)-SNAPSHOT
fi

RUN_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $RUN_PATH

./step-clean.sh
./step-compile.sh
