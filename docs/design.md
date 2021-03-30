Design
======

Context and scope
-----------------

pod-security-admission is a [Kubernetes Admission Webhook](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/) 
to ensure [Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/).

pod-security-admission aims to be a simple [Pod Security Policy](https://kubernetes.io/docs/concepts/policy/pod-security-policy/) replacement.

### Background

The enforcement of security settings for Pods has previously been achieved using Pod Security Policy (PSP).
However, it was [announced](https://github.com/kubernetes/kubernetes/pull/97171) that PSP is going to be removed in Kubernetes 1.25.
We need technology to replace PSP.

### Goals

- Enforcement of Pod Security Standards
- Easy to use without complicated settings
- Can specify the group of Pods to which each policy is applied
- Can customize rules to be applied for each policy 

### Non-goals

- Providing flexible policy engine, such as adding a new rule, customizing rules beyond Pod Security Standards

Webhooks
--------

Pod Security Standards define three policy types: `Privileged`, `Baseline` and `Restricted`.
pod-security-admission will provide ValidatingWebhook and MutatingWebhook endpoints for each of these policies.
However, for `Priviliged`, the policy can be archived by not applying the webhook endpoint.

Thus, by default, pod-security-admission serves four endpoints:
- Validating webhook for `Baseline`
- Mutating webhook for `Baseline`
- Validating webhook for `Restricted`
- Mutating webhook for `Restricted`

How to specify the group of Pods to which each policy is applied
----------------------------------------------------------------

In pod-security-admission, you can specify the Pods to which each policy is applied
with the label of their namespace.

Therefore, if the creator of the Pod could modify the namespace label,
they would be able to avoid applying the policy.

Cluster administrator must properly manage the permissions of namespace resources.

In PSP, you can do it with binding ServiceAccount to the profile using [RBAC](https://kubernetes.io/docs/reference/access-authn-authz/rbac/).
pod security admission does not provide that mechanism.
