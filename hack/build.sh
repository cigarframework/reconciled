#!/bin/bash

OS="$(go env GOHOSTOS)"
ARCH="$(go env GOARCH)"

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." >/dev/null 2>&1 && pwd & )"

bazel build //cmd/rd-server
mkdir ${ROOT}/bin
cp ${ROOT}/bazel-out/${OS}-fastbuild/bin/cmd/rd-server/${OS}_${ARCH}_stripped/rd-server ${ROOT}/bin/
docker build -t cigarframework/reconciled .

. ${ROOT}/hack/setup.sh
