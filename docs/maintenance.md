# Maintenance

## How to update supported Kubernetes

pod-security-admission supports the three latest Kubernetes versions.
If a new Kubernetes is released, please update the following files.

- Update Kubernetes version in `Makefile`.
- Update `k8s.io/*` and `sigs.k8s.io/controller-runtime` packages version in `go.mod`.

If Kubernetes or controller-runtime API has changed, please fix the relevant source code.

## How to update dependencies

TBD.
