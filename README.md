# Replica reloader

[![reuse compliant](https://reuse.software/badge/reuse-compliant.svg)](https://reuse.software/)

## Introduction

The replica-reloader is a simple controller that watches a Kubernetes `Deployment` and starts a `COMMAND` and adds the `Deployment` replica count as last argument:

```console
$(terminal-1) kubectl create deployment --image=nginx nginx
deployment.apps/nginx created
```

and running

```console
$(terminal-2) replica-reloader --namespace=default --deployment-name=nginx -- sleep
```

would start a `sleep 1` process.
If the watched deployment is scaled, then the controller stops the previous process and
starts a new one:

```console
$(terminal-1) kubectl scale deployment my-dep --replicas=10
deployment.apps/my-dep scaled

$(terminal-1) ps | grep sleep
61191 ttys003    0:00.00 sleep 10`,
```

## Why it was created

It was originally created as to be used in situations where you want to run some process with static count of replicas as flag (e.g. [apiserver-network-proxy](https://github.com/kubernetes-sigs/apiserver-network-proxy) that has a `--server-count` flag). In some, cases the process must be restarted with the new server count.

## Docker images

Docker images for the `replica-reloader` and `replica-reloader` bundled with `apiserver-network-proxy` are available at:

- `eu.gcr.io/gardener-project/gardener/replica-reloader:latest`
- `eu.gcr.io/gardener-project/gardener/replica-reloader:v0.2.0`
- `eu.gcr.io/gardener-project/gardener/replica-reloader:v0.2.0-konnectivity-server-v0.0.14`
- `eu.gcr.io/gardener-project/gardener/replica-reloader:v0.2.0-konnectivity-server-v0.0.15`
