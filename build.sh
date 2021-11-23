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

(
	source rhcos/interface.sh
	SHASUM=$(containerdisks::needs_update 4.9)
	if [ -n "${SHASUM}" ]; then
		docker build -t quay.io/containerdisks/rhcos:4.9 -t quay.io/containerdisks/rhcos:4.9-${TAG_SUFFIX} --label "shasum=${SHASUM}" -f rhcos/4.9/Dockerfile .
		if [ -z "${BUILD_ONLY}" ]; then
			docker push quay.io/containerdisks/rhcos:4.9-${TAG_SUFFIX}
			docker push quay.io/containerdisks/rhcos:4.9
		fi
	fi
)
