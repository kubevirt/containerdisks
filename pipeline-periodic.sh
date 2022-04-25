#!/bin/bash

set -e -o pipefail

make medius

export KUBEVIRT_MEMORY_SIZE=9216M
make cluster-up
trap '{ make cluster-down; }' EXIT SIGINT SIGTERM SIGSTOP

registry="$(./hack/kubevirtci.sh registry)"
kubeconfig="$(./hack/kubevirtci.sh kubeconfig)"
./bin/medius images push --no-fail --dry-run=false --target-registry=${registry} --insecure-skip-tls --workers 3
./bin/medius images verify --no-fail --dry-run=false --kubeconfig=${kubeconfig} --registry="registry:5000" --insecure-skip-tls --workers 3
./bin/medius images promote --dry-run=false --source-registry=${registry} --insecure-skip-tls --workers 3
