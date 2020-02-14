#!/bin/bash

OS="$(go env GOHOSTOS)"
ARCH="$(go env GOARCH)"
export PATH=$PWD/dist/bin/pkg/messages/generate/${OS}_${ARCH}_stripped:$PATH

bazel build //pkg/messages/generate
go generate ./pkg/...
