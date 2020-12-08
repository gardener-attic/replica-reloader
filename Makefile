# SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

VERSION                             := $(shell cat VERSION)
EFFECTIVE_VERSION                   := $(VERSION)-$(shell git rev-parse HEAD)
REPO_ROOT                           := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

ifneq ($(strip $(shell git status --porcelain 2>/dev/null)),)
	EFFECTIVE_VERSION := $(EFFECTIVE_VERSION)-dirty
endif

.PHONY: revendor
revendor:
	@GO111MODULE=on go mod tidy

.PHONY: check
check:
	@golangci-lint run --timeout 10m --config .golangci.yaml ./...

.PHONY: format
format:
	@goimports -l -w  ./

.PHONY: test
test:
	@GO111MODULE=on go test -race ./...

.PHONY: verify
verify: check format test
