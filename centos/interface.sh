#!/bin/bash

set -e -o pipefail

function containerdisks::needs_update() {
  local name=centos
  local tag="$1"
  local base_url="https://cloud.centos.org/centos/${tag}/x86_64/images/"

  local os_image="$(wget -O - "${base_url}" | grep "GenericCloud.*qcow2" | sed -e 's/.*>\(.*qcow2\)<.*/\1/' | sort -r | head -n 1)"
  local repo="quay.io/containerdisks/centos"


  local latest_shasum="$(wget -O - "${base_url}CHECKSUM" | grep ${os_image} | tail -n 1 | grep -o '[a-Z0-9]*$')"
  echo >&2 "Latest shasum for ${name}:${tag}: ${latest_shasum}"
  local release_shasum="$(skopeo inspect --no-tags docker://${repo}:${tag} --format '{{ .Labels.shasum }}')"
  echo >&2 "Published shasum for ${name}:${tag}: ${release_shasum}"
  if [[ "${latest_shasum}" == "${release_shasum}" ]]; then
    echo >&2 "No update for ${name}:${tag} required"
  else
    echo >&2 "Update for ${name}:${tag} required"
    echo "${latest_shasum}"
  fi
}
