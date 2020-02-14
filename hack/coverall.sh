#!/bin/bash

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." >/dev/null 2>&1 && pwd & )"
ROOT=${ROOT} go run ${ROOT}/hack/concat_dat_files.go
