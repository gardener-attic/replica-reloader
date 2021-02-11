# SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

ARG PROXY_SERVER_VERSION=v0.0.15
FROM k8s.gcr.io/kas-network-proxy/proxy-server:${PROXY_SERVER_VERSION} AS konnectivity-server-base

# make sure to update .devcontainer/devcontainer.json as well when changing this version.
FROM eu.gcr.io/gardener-project/3rd/golang:1.15.5 AS builder

WORKDIR /workdir
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o replica-reloader cmd/replica-reloader/main.go

############# base
FROM eu.gcr.io/gardener-project/3rd/alpine:3.12.1 AS base

############# replica-reloader #############
FROM base AS replica-reloader
ARG PROXY_SERVER_VERSION=v0.0.15

LABEL org.opencontainers.image.authors="Gardener contributors"
LABEL org.opencontainers.image.description="Small Kubernetes controller that watches Deployments and runs commands with replica count as last argument. It contains the konnectivity-server."
LABEL org.opencontainers.image.documentation="https://github.com/gardener/replica-reloader"
LABEL org.opencontainers.image.licenses="Apache-2.0"
LABEL org.opencontainers.image.source="https://github.com/gardener/replica-reloader"
LABEL org.opencontainers.image.title="replica reloader with konnectivty-server"
LABEL org.opencontainers.image.url="https://github.com/gardener/replica-reloader"
LABEL org.opencontainers.image.vendor="Gardener contributors"
LABEL org.opencontainers.image.version="${PROXY_SERVER_VERSION}"

COPY --from=builder /workdir/replica-reloader .
COPY --from=konnectivity-server-base /proxy-server .

WORKDIR /

ENTRYPOINT ["/replica-reloader"]
