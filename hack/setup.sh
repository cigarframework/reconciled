#!/bin/bash

OS="$(go env GOHOSTOS)"
ARCH="$(go env GOARCH)"

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." >/dev/null 2>&1 && pwd & )"

rm -rf ${ROOT}/pkg/proto

cp -r ${ROOT}/bazel-out/${OS}-fastbuild/bin/_proto/${OS}_${ARCH}_stripped/proto_go_proto%/github.com/cigarframework/reconciled/pkg/proto ${ROOT}/pkg/proto
