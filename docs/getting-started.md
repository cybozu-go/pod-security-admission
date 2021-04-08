Getting Started
===============

Prepare Kubernetes Cluster
--------------------------

pod-security-admission can run on [kind](https://kind.sigs.k8s.io) clusters using Docker.

If you don't have your Kubernetes cluster, [setup kind](https://kind.sigs.k8s.io/docs/user/quick-start/), then run:

```console
$ kind create cluster
```

Install cert-manager
--------------------

In order to use Admission Webhook, a certificate is required.
Let's install [cert-manager](https://cert-manager.io/docs/): a native Kubernetes certificate management controller.
Please see the [document](https://cert-manager.io/docs/installation/kubernetes/) for details.

Deploy pod-security-admission
-----------------------------

To make the system namespaces privileged, label those namespaces:

```console
$ kubectl label namespace/kube-system pod-security.cybozu.com/policy=privileged
$ kubectl label namespace/cert-manager pod-security.cybozu.com/policy=privileged
```

Deploy pod-security-admission:

```console
$ kubectl apply -f https://github.com/cybozu-go/pod-security-admission/releases/download/v0.0.1-alpha.0/install.yaml
```

Verification
------------

Now if you create a Pod that violates the policy, it will be rejected.

```console
$ kubectl apply -f hooks/testdata/baseline/additional-capability.yaml
Error from server (spec.containers[0].securityContext.capabilities.add[1]: Forbidden: Adding capability SYSLOG is not allowed): error when creating "hooks/testdata/baseline/additional-capability.yaml": admission webhook "baseline.vpod.kb.io" denied the request: spec.containers[0].securityContext.capabilities.add[1]: Forbidden: Adding capability SYSLOG is not allowed
```
