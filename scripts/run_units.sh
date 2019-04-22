#!/usr/bin/env bash
#
# Copyright (c) 2019 LG Electronics Inc.
# SPDX-License-Identifier: Apache-2.0
#

set -e -u

cd $(dirname $0)/..

echo "Setting up ginkgo and gomega"
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega

ginkgo  -r \
  -skipPackage controller \
  "$@"
