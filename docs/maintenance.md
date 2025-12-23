# Maintenance

## How to update supported Kubernetes

pod-security-admission supports the three latest Kubernetes versions.
If a new Kubernetes is released, please update the following files.

- Update Kubernetes version in `Makefile`.
- Update `k8s.io/*` and `sigs.k8s.io/controller-runtime` packages version in `go.mod`.
- Update `aqua.yaml` by running `aqua update --select-version`.
    - It has been observed that automatic `aqua update` selects inappropriate packages such as `kustomize@kustomize/v5.8.0` and `controller-tools/controller-gen@envtest-v1.35.0`.

If Kubernetes or controller-runtime API has changed, please fix the relevant source code.

## How to update dependencies

Renovate will create PRs that update dependencies once a week.
However, Kubernetes is only updated with patched versions.
