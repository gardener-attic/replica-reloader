# SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

# make sure to update .devcontainer/devcontainer.json as well when changing this version.
FROM eu.gcr.io/gardener-project/3rd/golang:1.15.5 AS replica-reloader-builder

WORKDIR /workdir
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o replica-reloader cmd/replica-reloader/main.go

############# replica-reloader #############
FROM eu.gcr.io/gardener-project/3rd/alpine:3.12.1

COPY --from=replica-reloader-builder /workdir/replica-reloader .

WORKDIR /

LABEL org.opencontainers.image.authors="Gardener contributors"
LABEL org.opencontainers.image.description="Small Kubernetes controller that watches Deployments and runs commands with replica count as last argument."
LABEL org.opencontainers.image.documentation="https://github.com/gardener/replica-reloader"
LABEL org.opencontainers.image.licenses="Apache-2.0"
LABEL org.opencontainers.image.source="https://github.com/gardener/replica-reloader"
LABEL org.opencontainers.image.title="replica reloader"
LABEL org.opencontainers.image.url="https://github.com/gardener/replica-reloader"
LABEL org.opencontainers.image.vendor="Gardener contributors"

ENTRYPOINT ["/replica-reloader"]
