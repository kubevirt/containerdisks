#!/bin/bash

set -e -o pipefail

export KUBEVIRT_MEMORY_SIZE=9216M
PROMOTE_DRY_RUN=${PROMOTE_DRY_RUN:-false}

make medius
make cluster-up
trap '{ make cluster-down; }' EXIT SIGINT SIGTERM

registry="$(./hack/kubevirtci.sh registry)"
kubeconfig="$(./hack/kubevirtci.sh kubeconfig)"

medius_push_extra_args=()
if [ "$FORCE_REBUILD" = "true" ]; then
  medius_push_extra_args+=("--force")
fi

./bin/medius images push --no-fail --dry-run=false --target-registry="${registry}" --insecure-skip-tls --workers 3 "${medius_push_extra_args[@]}"
./bin/medius images verify --no-fail --dry-run=false --kubeconfig="${kubeconfig}" --registry="registry:5000" --insecure-skip-tls --workers 3
./bin/medius images promote --dry-run="${PROMOTE_DRY_RUN}" --source-registry="${registry}" --insecure-skip-tls --workers 3

if [ -n "${QUAY_OAUTH_TOKEN}" ]; then
    ./bin/medius docs publish --dry-run=false --quay-token-file="${QUAY_OAUTH_TOKEN}"
fi
