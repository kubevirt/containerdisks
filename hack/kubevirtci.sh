#!/bin/bash

set -e

KUBEVIRTCI_TAG=${KUBEVIRTCI_TAG:-$(curl -sfL https://storage.googleapis.com/kubevirt-prow/release/kubevirt/kubevirtci/latest)}
export KUBEVIRTCI_TAG

export KUBEVIRT_PROVIDER=${KUBEVIRT_PROVIDER:-k8s-1.23}
export KUBEVIRT_NUM_NODES=${KUBEVIRT_NUM_NODES:-1}
export KUBEVIRT_MEMORY_SIZE=${KUBEVIRT_MEMORY_SIZE:-4096M}
export KUBEVIRT_DEPLOY_CDI="true"

_base_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
_kubectl="${_base_dir}/cluster-up/cluster-up/kubectl.sh"
_action=$1
shift

function kubevirtci::fetch_kubevirtci() {
	if [[ ! -d ${_base_dir}/cluster-up ]]; then
    git clone https://github.com/kubevirt/kubevirtci.git "${_base_dir}/cluster-up"
	  (cd "${_base_dir}/cluster-up" && git checkout "$KUBEVIRTCI_TAG")
  fi
}

function kubevirtci::up() {
  make cluster-up -C "${_base_dir}/cluster-up"
  KUBECONFIG=$(kubevirtci::kubeconfig)
  export KUBECONFIG
  echo "adding kubevirtci registry to cdi-insecure-registries"
  ${_kubectl} patch configmap cdi-insecure-registries -n cdi --type merge -p '{"data":{"kubevirtci": "registry:5000"}}'
  echo "installing kubevirt..."
  LATEST=$(curl -L https://storage.googleapis.com/kubevirt-prow/devel/release/kubevirt/kubevirt/stable.txt)
  ${_kubectl} apply -f "https://github.com/kubevirt/kubevirt/releases/download/${LATEST}/kubevirt-operator.yaml"
  ${_kubectl} apply -f "https://github.com/kubevirt/kubevirt/releases/download/${LATEST}/kubevirt-cr.yaml"
  echo "waiting for kubevirt to become ready, this can take a few minutes..."
  ${_kubectl} -n kubevirt wait kv kubevirt --for condition=Available --timeout=15m
}

function kubevirtci::down() {
  make cluster-down -C "${_base_dir}/cluster-up"
}

function kubevirtci::kubeconfig() {
  "${_base_dir}/cluster-up/cluster-up/kubeconfig.sh"
}

function kubevirtci::registry() {
  echo "localhost:$("${_base_dir}/cluster-up/cluster-up/cli.sh" ports registry 2>/dev/null)"
}

kubevirtci::fetch_kubevirtci

case ${_action} in
  "up")
    kubevirtci::up
    ;;
  "down")
    kubevirtci::down
    ;;
  "kubeconfig")
    kubevirtci::kubeconfig
    ;;
  "registry")
    kubevirtci::registry
    ;;
  "kubectl")
    ${_kubectl} "$@"
    ;;
  *)
    echo "No command provided, known commands are 'up', 'down', 'kubeconfig', 'registry', 'kubectl'"
    exit 1
    ;;
esac
