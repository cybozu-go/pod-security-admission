Configuration
=============

Customize Policy
----------------

`pod-security-admission` enforces policies that conform to [Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/) by default.
You can customize it.

### Validators

`pod-security-admission` provides the following validators:

| name                       | policy type | description                                               |
| -------------------------- | ----------- | --------------------------------------------------------- |
| deny-host-namespace        | baseline    | deny sharing the host namespaces                          |
| deny-privileged-containers | baseline    | deny privileged containers                                |
| deny-unsafe-capabilities   | baseline    | deny adding capabilities beyond the default set           |
| deny-host-path-volumes     | baseline    | deny usage of HostPath volumes                            |
| deny-host-ports            | baseline    | deny usage of HostPorts                                   |
| deny-unsafe-apparmor       | baseline    | deny overriding or disabling the default AppArmor profile |
| deny-unsafe-selinux        | baseline    | deny setting custom SELinux options                       |
| deny-unsafe-proc-mount     | baseline    | deny unmasked proc mount                                  |
| deny-unsafe-sysctls        | baseline    | deny usage of unsafe sysctls                              |
| deny-non-core-volume-types | restricted  | deny usage of non-core volume types                       |
| deny-privilege-escalation  | restricted  | deny privilege escalation                                 |
| deny-run-as-root           | restricted  | deny running as root users                                |
| deny-root-groups           | restricted  | deny running with a root primary or supplementary GID     |
| deny unsafe-seccomp        | restricted  | deny usage of non-default Seccomp profile                 |

### Mutators

`pod-security-admission` provides the following mutators (Not enabled by default):

| name                  | policy type | description                       |
| --------------------- | ----------- | --------------------------------- |
| force-run-as-non-root | -           | force running with non-root users |

### Customize

By default, `pod-security-admission` uses the following configuration:

```yaml
- name: baseline
  denyHostNamespace: true
  denyPrivilegedContainers: true
  capabilities:
    denyUnsafeCapabilities: true
  volumes:
    denyHostPathVolumes: true
  hostPorts:
    denyHostPorts: true
  denyUnsafeApparmor: true
  denyUnsafeSelinux: true
  denyUnsafeProcMount: true
  denyUnsafeSysctls: true
- name: restricted
  volumes:
    denyNonCoreVolumeTypes: true
  denyPrivilegeEscalation: true
  runAsRoot:
    denyRunAsRoot: true
  denyRootGroups: true
  denyUnsafeSeccomp: true
```

For example, if you want to enforce `deny-run-as-root` and `force-run-as-non-root` in `Baseline`,
administrators can add the rules under the `Baseline` section: 

```yaml
- name: baseline
  denyHostNamespace: true
  denyPrivilegedContainers: true
  capabilities:
    denyUnsafeCapabilities: true
    allowedCapabilities:
      - SYSLOG
      - NET_ADMIN
  volumes:
    denyHostPathVolumes: true
  hostPorts:
    denyHostPorts: true
    allowedHostPorts:
      - min: 1024
        max: 65535
  denyUnsafeApparmor: true
  denyUnsafeSelinux: true
  denyUnsafeProcMount: true
  denyUnsafeSysctls: true
  runAsRoot:
    denyRunAsRoot: true
    forceRunAsNonRoot: true
- name: restricted
  volumes:
    denyNonCoreVolumeTypes: true
  denyPrivilegeEscalation: true
  denyRootGroups: true
  denyUnsafeSeccomp: true
```

### Webhook Configuration

You can configure [ValidatingWebhookConfiguration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#validatingwebhookconfiguration-v1-admissionregistration-k8s-io) and [MutatingWebhookConfiguration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#mutatingwebhookconfiguration-v1-admissionregistration-k8s-io) to change which Pods are covered by each policy.

See [manifests.yaml](../config/webhook/manifests.yaml) in details.

The endpoint of the webhook is determined by the name of the policy, such as `baseline` and `restricted`. 
A validating webhook will be `/validate-` + the policy name, a mutating webhook will be `/mutate-` + the policy name.
For example, the endpoint of the validating webhook for `baseline` will be `/validate-baseline`, mutating webhook will be `/mutate-baseline`.
