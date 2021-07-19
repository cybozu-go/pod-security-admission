Configuration
=============

`pod-security-admission` can specify a profile for each webhook endpoint.

SecurityProfile 
----------------

SecurityProfile has these fields:

| Name                     | Type              | Description                                                                                                                                 |
| ------------------------ | ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------- |
| name                     | string            | The name of the profile                                                                                                                     |
| hostNamespace            | bool              | Allow sharing the host namespaces                                                                                                           |
| privileged               | bool              | Allow privileged containers                                                                                                                 |
| capabilities             | bool              | Allow adding capabilities beyond the [default set](https://docs.docker.com/engine/reference/run/#runtime-privilege-and-linux-capabilities). |
| additionalCapabilities   | []string          | The list of capabilities that cab be added. If `capabilities` is true, this list will be ignored.                                           |
| hostPathVolumes          | bool              | Allow usage of HostPath volumes                                                                                                             |
| allowedHostPaths         | []AllowedHostPath | The list of host paths that can be used. If `hostPathVolumes` is true, this list will be ignored.                                           |
| nonCoreVolumeTypes       | bool              | Allow usage of non-core volume types, except HostPath volumes                                                                               |
| hostPorts                | bool              | Allow usage of all HostPorts                                                                                                                |
| allowedHostPorts         | []PortRange       | The list of host ports that can be used. If `hostPorts` is true, this list will be ignored.                                                 |
| appArmor                 | bool              | Allow overriding or disabling the default AppArmor profile                                                                                  |
| seLinux                  | bool              | Allow setting custom SELinux options                                                                                                        |
| procMount                | bool              | Allow unmasked proc mount                                                                                                                   |
| sysctls                  | bool              | Allow usage of unsafe sysctls                                                                                                               |
| allowPrivilegeEscalation | bool              | Allow privilege escalation                                                                                                                  |
| runAsRoot                | bool              | Allow running as root users                                                                                                                 |
| forceRunAsNonRoot        | bool              | Force running with non-root users by MutatingWebhook                                                                                        |
| rootGroups               | bool              | Allow running with a root primary or supplementary GID                                                                                      |
| seccomp                  | bool              | Allow usage of non-default Seccomp profile                                                                                                  |

#### AllowedHostPath

| Name       | Type   | Description                   |
| ---------- | ------ | ----------------------------- |
| pathPrefix | string | The path prefix to be allowed |

### PortRange

| Name | Type  | Description                            |
| ---- | ----- | -------------------------------------- |
| min  | int32 | The min of host port range (inclusive) |
| max  | int32 | The max of host port range (inclusive) |

Customize Profile
-----------------

By default, `pod-security-admission` uses the following configuration to enforce [Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/):
The profile describes only the allowed items. A false or no description indicates that the item is to be denied.


```yaml
- name: baseline
  nonCoreVolumeTypes: true
  allowPrivilegeEscalation: true
  runAsRoot: true
  rootGroups: true
  seccomp: true
- name: restricted
  forceRunAsNonRoot: true
```

Administrators can customize the profile if needed.
The following `Baseline` profile allows hostNamespaces, to add `SYSLOG` and `NET_ADMIN` capabilities, use of hostPorts from 1024 to 65535.

```yaml
- name: baseline
  hostNamespace: true
  additionalCapabilities:
    - SYSLOG
    - NET_ADMIN
  nonCoreVolumeTypes: true
  allowedHostPorts:
    - min: 1024
      max: 65535
  allowPrivilegeEscalation: true
  rootGroups: true
  runAsRoot: true
  seccomp: true
- name: restricted
  forceRunAsNonRoot: true
```

### Webhook Configuration

You can configure [ValidatingWebhookConfiguration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#validatingwebhookconfiguration-v1-admissionregistration-k8s-io) and [MutatingWebhookConfiguration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#mutatingwebhookconfiguration-v1-admissionregistration-k8s-io) to change which Pods are covered by each policy.

See [manifests.yaml](../config/webhook/manifests.yaml) for details.

The endpoint of the webhook is determined by the name of the policy, such as `baseline` and `restricted`. 
A validating webhook will be `/validate-` + the policy name, a mutating webhook will be `/mutate-` + the policy name.
For example, the endpoint of the validating webhook for `baseline` will be `/validate-baseline`, mutating webhook will be `/mutate-baseline`.
