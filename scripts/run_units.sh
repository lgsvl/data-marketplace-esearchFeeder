#!/usr/bin/env bash

set -e -u

cd $(dirname $0)/..

echo "Setting up ginkgo and gomega"
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega

ginkgo  -r \
  -skipPackage controller \
  "$@"