#!/bin/bash

set -e -o pipefail

function containerdisks::needs_update() {
  local tag=$1
  local latest_shasum="$(curl https://mirror.openshift.com/pub/openshift-v4/dependencies/rhcos/${tag}/latest/sha256sum.txt | grep rhcos-openstack.x86_64.qcow2.gz | awk '{ print $1}')"
  echo >&2 "Latest shasum for rhcos:${tag}: ${latest_shasum}"
  local release_shasum="$(skopeo inspect --no-tags docker://quay.io/containerdisks/rhcos:${tag} --format '{{ .Labels.shasum }}')"
  echo >&2 "Published shasum for rhcos:${tag}: ${release_shasum}"
  if [[ "${latest_shasum}" == "${release_shasum}" ]]; then
    echo >&2 "No update for rhcos:${tag} required"
  else
    echo >&2 "Update for rhcos:${tag} required"
    echo "${latest_shasum}"
  fi
}
