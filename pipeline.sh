#!/bin/bash

set -e -o pipefail

make medius
make cluster-up

FOCUS=${FOCUS:-centos-stream:9}
registry="$(./hack/kubevirtci.sh registry)"
kubeconfig="$(./hack/kubevirtci.sh kubeconfig)"
./bin/medius images push --force --focus="${FOCUS}" --no-fail --dry-run=false --source-registry="${registry}" --insecure-skip-tls
./bin/medius images verify --focus="${FOCUS}" --no-fail --dry-run=false --kubeconfig="${kubeconfig}" --registry="registry:5000" --insecure-skip-tls
./bin/medius images promote --focus="${FOCUS}" --dry-run=true --source-registry="${registry}" --insecure-skip-tls

make cluster-down
