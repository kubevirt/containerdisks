#!/bin/bash
#
# This file is part of the KubeVirt project
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# Copyright The KubeVirt Authors.
#

set -ex

export KUBEVIRT_MEMORY_SIZE="${KUBEVIRT_MEMORY_SIZE:-16G}"
export KUBEVIRT_DEPLOY_CDI="true"
export KUBEVIRT_VERSION=${KUBEVIRT_VERSION:-$(curl -fL https://storage.googleapis.com/kubevirt-prow/devel/release/kubevirt/kubevirt/stable.txt)}
export KUBEVIRTCI_TAG=${KUBEVIRTCI_TAG:-$(curl -fL https://raw.githubusercontent.com/kubevirt/kubevirt/"${KUBEVIRT_VERSION}"/kubevirtci/cluster-up/version.txt)}

_base_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
_cluster_up_dir="${_base_dir}/_cluster-up"
_kubectl="${_cluster_up_dir}/cluster-up/kubectl.sh"
_kubessh="${_cluster_up_dir}/cluster-up/ssh.sh"
_kubevirtcicli="${_cluster_up_dir}/cluster-up/cli.sh"
_action=$1
shift

function kubevirtci::fetch_kubevirtci() {
	if [[ ! -d "${_cluster_up_dir}" ]]; then
    git clone --depth 1 --branch "${KUBEVIRTCI_TAG}" https://github.com/kubevirt/kubevirtci.git "${_cluster_up_dir}"
  fi
}

function kubevirtci::up() {
  make cluster-up -C "${_cluster_up_dir}"
  KUBECONFIG=$(kubevirtci::kubeconfig)
  export KUBECONFIG

  echo "adding kubevirtci registry to cdi-insecure-registries"
  ${_kubectl} patch cdis/cdi --type merge -p '{"spec": {"config": {"insecureRegistries": ["registry:5000"]}}}'

  echo "installing kubevirt..."
  ${_kubectl} apply -f "https://github.com/kubevirt/kubevirt/releases/download/${KUBEVIRT_VERSION}/kubevirt-operator.yaml"
  ${_kubectl} apply -f "https://github.com/kubevirt/kubevirt/releases/download/${KUBEVIRT_VERSION}/kubevirt-cr.yaml"

  # Remove once kubevirtci in the stable release uses Kubernetes 1.36+ by default
  echo "disable ImageVolume feature gate"
  ${_kubectl} patch kv/kubevirt -n kubevirt --type merge -p '{"spec":{"configuration":{"developerConfiguration": {"disabledFeatureGates": ["ImageVolume"]}}}}'

  echo "waiting for kubevirt to become ready, this can take a few minutes..."
  ${_kubectl} -n kubevirt wait kv kubevirt --for condition=Available --timeout=15m
}

function kubevirtci::down() {
  make cluster-down -C "${_cluster_up_dir}"
}

function kubevirtci::kubeconfig() {
  "${_cluster_up_dir}/cluster-up/kubeconfig.sh"
}

function kubevirtci::registry() {
  port=$(${_kubevirtcicli} ports registry 2>/dev/null)
  echo "localhost:${port}"
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
  "ssh")
    ${_kubessh} "$@"
    ;;
  "kubectl")
    ${_kubectl} "$@"
    ;;
  *)
    echo "No command provided, known commands are 'up', 'down', 'kubeconfig', 'registry', 'ssh' and 'kubectl'"
    exit 1
    ;;
esac
