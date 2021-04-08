[![Project Status](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#active)
[![GitHub release](https://img.shields.io/github/release/cybozu-go/pod-security-admission.svg?maxAge=60)][releases]
![CI](https://github.com/cybozu-go/pod-security-admission/workflows/CI/badge.svg)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/cybozu-go/pod-security-admission?tab=overview)](https://pkg.go.dev/github.com/cybozu-go/pod-security-admission?tab=overview)
[![Go Report Card](https://goreportcard.com/badge/github.com/cybozu-go/pod-security-admission)](https://goreportcard.com/report/github.com/cybozu-go/pod-security-admission)

***NOTE***

The PSP replacement has been [announced](https://kubernetes.io/blog/2021/04/06/podsecuritypolicy-deprecation-past-present-and-future/).
This project is just a stopgap until it is replaced.

Pod Security Admission
======================

pod-security-admission is a set of [Kubernetes Admission Webhooks](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/) to ensure [Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/).

pod-security-admission aims to be a simple [Pod Security Policy](https://kubernetes.io/docs/concepts/policy/pod-security-policy/) replacement.

This is not a policy engine and users cannot write their own policies flexibly.
If you want to do that, I recommend using a policy engine such as [OPA/GateKeeper](https://open-policy-agent.github.io/gatekeeper) and [Kyverno](https://kyverno.io).

Getting started
---------------

Please see the [getting-started.md](./docs/getting-started.md) to deploy `pod-security-admission` to your Kubernetes cluster.

Policies
--------

pod-security-admission provides 3 policy types based on Pod Security Standards.

### Privileged

The `Privileged` is an entirely unrestricted policy.
Admission webhook does nothing to the Pods in namespaces with `Privileged` label.
This policy should be applied to the Pods that are the core components for the Kubernetes cluster, such as network plugins.

This policy will be applied to Pods that belong to namespaces with the following label:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: your-namespace
  labels: 
    pod-security.cybozu.com/policy: privileged
```

### Baseline

The `Baseline` is a basic policy that applies to Pods for many applications.

This policy prohibits the creation of Pods that violate the following rules:
- https://kubernetes.io/docs/concepts/security/pod-security-standards/#baseline

This policy will be applied to Pods that belong to all namespaces except privileged.

### Restricted

The `Restricted` is a restricted policy that applies to Pods for secure applications.

In addition to the `Baseline`, this policy prohibits the creation of Pods that violate the following rules:
- https://kubernetes.io/docs/concepts/security/pod-security-standards/#restricted

Furthermore, if a Pod violates `Running as Non-root`, the mutating webhook will rewrite the securityContext forcibly.

This policy will be applied to Pods that belong to namespaces with the following label:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: your-namespace
  labels:
    pod-security.cybozu.com/policy: restricted
```

Documentation
-------------

[docs](docs/) directory contains documents about designs and specifications.

Docker images
-------------

Docker images are available on [Quay.io](https://quay.io/repository/cybozu/pod-security-admission)

[releases]: https://github.com/cybozu-go/pod-security-admission/releases
