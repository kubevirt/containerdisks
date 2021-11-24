#!/bin/bash

set -e -o pipefail

function containerdisks::needs_update() {
  local tag=$1
  local latest_shasum="$(curl https://getfedora.org/releases.json | jq -r --arg VERSION $tag '.[] | select(.link|test(".*qcow2")) | select(.variant=="Cloud" and .arch=="x86_64" and .version==$VERSION).sha256')"
  echo >&2 "Latest shasum for fedora:${tag}: ${latest_shasum}"
  local release_shasum="$(skopeo inspect --no-tags docker://quay.io/containerdisks/fedora:${tag} --format '{{ .Labels.shasum }}')"
  echo >&2 "Published shasum for fedora:${tag}: ${release_shasum}"
  if [[ "${latest_shasum}" == "${release_shasum}" ]]; then
    echo >&2 "No update for fedora:${tag} required"
  else
    echo >&2 "Update for fedora:${tag} required"
    echo "${latest_shasum}"
  fi
}
