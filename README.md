[![GitHub release](https://img.shields.io/github/release/cybozu-go/pod-security-admission.svg?maxAge=60)][releases]
[![CircleCI](https://circleci.com/gh/cybozu-go/pod-security-admission.svg?style=svg)](https://circleci.com/gh/cybozu-go/pod-security-admission)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/cybozu-go/pod-security-admission?tab=overview)](https://pkg.go.dev/github.com/cybozu-go/pod-security-admission?tab=overview)
[![Go Report Card](https://goreportcard.com/badge/github.com/cybozu-go/pod-security-admission)](https://goreportcard.com/report/github.com/cybozu-go/pod-security-admission)

Pod Security Admission
======================

**Project Status**: Initial development

`pod-security-admission` is composed of [validation and mutating admission webhooks](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/) to ensure [Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/).

The enforcement of security settings for Pods has previously been achieved using [Pod Security Policy (PSP)](https://kubernetes.io/docs/concepts/policy/pod-security-policy/).
However, it was [announced](https://github.com/kubernetes/kubernetes/pull/97171) that PSP is going to be removed in Kubernetes 1.25.

`pod-security-admission` provides an enforcement of [simple Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/) as admission webhook.

`pod-security-admission` is not a policy engine and users cannot write their own policies flexibly.
It does not provide flexible policy changes.
If you want to do that, I recommend using a policy engine such as [OPA GateKeeper](https://github.com/open-policy-agent/gatekeeper).

Webhooks
--------

`pod-security-admission` provides 3 policy types based on Pod Security Standards.

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

[docs](docs/) directory contains documents about specifications.

Docker images
-------------

Docker images are available on [Quay.io](https://quay.io/repository/cybozu/pod-security-admission)

[releases]: https://github.com/cybozu-go/pod-security-admission/releases
