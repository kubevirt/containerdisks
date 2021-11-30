#!/bin/bash

set -e -o pipefail

unset BUILD_ONLY

while getopts "b" arg; do
  case $arg in
  b)
    export BUILD_ONLY=true
    ;;
  esac
done

opts=""
if [ -z "${BUILD_ONLY}" ]; then
  opts="--dry-run=false"
fi
make medius
./bin/medius publish ${opts} --workers 3
