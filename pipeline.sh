#!/bin/bash

set -e -o pipefail

make medius
make cluster-up

registry="$(./hack/kubevirtci.sh registry)"
kubeconfig="$(./hack/kubevirtci.sh kubeconfig)"
./bin/medius images push --force --no-fail --dry-run=false --source-registry=${registry} --insecure-skip-tls --workers 3
./bin/medius images verify --no-fail --dry-run=false --kubeconfig=${kubeconfig} --registry="registry:5000" --insecure-skip-tls --workers 3
./bin/medius images promote --dry-run=true --source-registry=${registry} --insecure-skip-tls --workers 3

make cluster-down
