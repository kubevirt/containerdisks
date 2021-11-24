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

export TAG_SUFFIX=$(date +"%y%m%d%H%M")

function build() {
  export distro=$1
  export tag=$2
  (
    source ${distro}/interface.sh
    SHASUM=$(containerdisks::needs_update ${tag})
    if [ -n "${SHASUM}" ]; then
      docker build -t quay.io/containerdisks/${distro}:${tag} -t quay.io/containerdisks/${distro}:${tag}-${TAG_SUFFIX} --label "shasum=${SHASUM}" -f ${distro}/${tag}/Dockerfile .
      if [ -z "${BUILD_ONLY}" ]; then
        docker push quay.io/containerdisks/${distro}:${tag}-${TAG_SUFFIX}
        docker push quay.io/containerdisks/${distro}:${tag}
      fi
    fi
  )
}

build rhcos 4.9
build fedora 35
