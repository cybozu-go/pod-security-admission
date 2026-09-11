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

Dependencies are updated manually. Renovate is no longer used in this
repository (the workflow was removed in #129 because it kept creating
dependency dashboard issues).

- Go modules: run `go get -u ./...` (or update specific modules) and
  `go mod tidy`. Keep `k8s.io/*` and `sigs.k8s.io/controller-runtime` on
  the versions that match the supported Kubernetes releases described above.
- Tools managed by aqua: run `aqua update` and review the result, then run
  `aqua update-checksum`.
- GitHub Actions: update the pinned commit hashes in `.github/workflows/`.
