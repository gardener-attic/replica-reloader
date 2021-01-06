# SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

FROM eu.gcr.io/gardener-project/3rd/golang:1.15.5 AS builder

WORKDIR /workdir
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o replica-reloader cmd/replica-reloader/main.go

FROM k8s.gcr.io/kas-network-proxy/proxy-server:v0.0.14 AS konnectivity-server-base

############# base
FROM eu.gcr.io/gardener-project/3rd/alpine:3.12.1 AS base

############# replica-reloader #############
FROM base AS replica-reloader

COPY --from=builder /workdir/replica-reloader .
COPY --from=konnectivity-server-base /proxy-server .

WORKDIR /

ENTRYPOINT ["/replica-reloader"]
