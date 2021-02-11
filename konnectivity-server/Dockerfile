# SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

ARG REPLICA_RELOADER_IMAGE
ARG PROXY_SERVER_VERSION=v0.0.15

FROM k8s.gcr.io/kas-network-proxy/proxy-server:${PROXY_SERVER_VERSION} AS konnectivity-server-base

FROM ${REPLICA_RELOADER_IMAGE}
ARG PROXY_SERVER_VERSION=v0.0.15
############# konnectivity-server #############

LABEL org.opencontainers.image.authors="Gardener contributors"
LABEL org.opencontainers.image.description="Small Kubernetes controller that watches Deployments and runs commands with replica count as last argument. It contains the konnectivity-server."
LABEL org.opencontainers.image.documentation="https://github.com/gardener/replica-reloader"
LABEL org.opencontainers.image.licenses="Apache-2.0"
LABEL org.opencontainers.image.source="https://github.com/gardener/replica-reloader"
LABEL org.opencontainers.image.title="replica reloader with konnectivty-server"
LABEL org.opencontainers.image.url="https://github.com/gardener/replica-reloader"
LABEL org.opencontainers.image.vendor="Gardener contributors"
LABEL org.opencontainers.image.version="${PROXY_SERVER_VERSION}"

COPY --from=konnectivity-server-base /proxy-server .

WORKDIR /

ENTRYPOINT ["/replica-reloader"]
