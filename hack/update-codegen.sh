#!/usr/bin/env bash

# Copyright 2017 The Kubernetes Authors.
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

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[@]}")/..

TOOLS_DIR=$(realpath ./hack/tools)
TOOLS_BIN_DIR="${TOOLS_DIR}/bin"
GO_INSTALL=$(realpath ./hack/go-install.sh)
CONTROLLER_GEN_VER=v0.11.1
CONTROLLER_GEN_BIN=controller-gen
CONTROLLER_GEN=${TOOLS_BIN_DIR}/${CONTROLLER_GEN_BIN}-${CONTROLLER_GEN_VER}

GOBIN=${TOOLS_BIN_DIR} ${GO_INSTALL} sigs.k8s.io/controller-tools/cmd/controller-gen ${CONTROLLER_GEN_BIN} ${CONTROLLER_GEN_VER}

CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

bash "${CODEGEN_PKG}"/generate-internal-groups.sh \
  "deepcopy,conversion,defaulter" \
  github.com/tanjunchen/tanjunchen-scheduler/pkg/generated \
  github.com/tanjunchen/tanjunchen-scheduler/apis \
  github.com/tanjunchen/tanjunchen-scheduler/apis \
  "config:v1,v1beta2,v1beta3" \
  --trim-path-prefix github.com/tanjunchen/tanjunchen-scheduler \
  --output-base "./" \
  --go-header-file "${SCRIPT_ROOT}"/hack/boilerplate/boilerplate.generatego.txt
