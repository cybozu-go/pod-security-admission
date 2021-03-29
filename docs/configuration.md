Configuration
=============

`pod-security-admission` can specify a profile for each webhook endpoint.

SecurityProfile 
----------------

SecurityProfile has these fields:

| Name                     | Type              | Description                                               |
| ------------------------ | ----------------- | --------------------------------------------------------- |
| name                     | string            | The name of the profile                                   |
| denyHostNamespace        | bool              | Deny sharing the host namespaces                          |
| denyPrivilegedContainers | bool              | Deny privileged containers                                |
| capabilities             | CapabilityProfile | The profile for capabilities                              |
| volumes                  | VolumeProfile     | The profile for volumes                                   |
| hostPorts                | HostPortProfile   | The profile for hostPorts                                 |
| denyUnsafeApparmor       | bool              | Deny overriding or disabling the default AppArmor profile |
| denyUnsafeSelinux        | bool              | Deny setting custom SELinux options                       |
| denyUnsafeProcMount      | bool              | Deny unmasked proc mount                                  |
| denyUnsafeSysctls        | bool              | Deny usage of unsafe sysctls                              |
| denyPrivilegeEscalation  | bool              | Deny privilege escalation                                 |
| users                    | UserProfile       | The profile for users                                     |
| denyRootGroups           | bool              | Deny running with a root primary or supplementary GID     |
| denyUnsafeSeccomp        | bool              | Deny usage of non-default Seccomp profile                 |

### CapabilityProfile

| Name                   | Type     | Description                                                                                                                                                         |
| ---------------------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| denyUnsafeCapabilities | bool     | Deny adding capabilities beyond the [default set](https://docs.docker.com/engine/reference/run/#runtime-privilege-and-linux-capabilities) and `allowedCapabilities` |
| allowedCapabilities    | []string | The list of capabilities that cab be added                                                                                                                          |

### VolumeProfile

| Name                   | Type | Description                         |
| ---------------------- | ---- | ----------------------------------- |
| denyHostPathVolumes    | bool | Deny usage of HostPath volumes      |
| denyNonCoreVolumeTypes | bool | Deny usage of non-core volume types |

### HostPortProfile

| Name             | Type        | Description                                           |
| ---------------- | ----------- | ----------------------------------------------------- |
| denyHostPorts    | bool        | Deny usage of HostPorts except for `allowedHostPorts` |
| allowedHostPorts | []PortRange | The list of host ports that can be used               |

### PortRange

| Name | Type  | Description                            |
| ---- | ----- | -------------------------------------- |
| min  | int32 | The min of host port range (inclusive) |
| max  | int32 | The max of host port range (inclusive) |

### UserProfile

| Name              | Type | Description                                          |
| ----------------- | ---- | ---------------------------------------------------- |
| denyRunAsRoot     | bool | Deny running as root users                           |
| forceRunAsNonRoot | bool | Force running with non-root users by MutatingWebhook |


Customize Profile
-----------------

By default, `pod-security-admission` uses the following configuration to enforce [Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/):

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
  users:
    denyRunAsRoot: true
  denyRootGroups: true
  denyUnsafeSeccomp: true
```

For example, administrators can customize the profile under the `Baseline` section as follows.
This profile allows to add `SYSLOG` and `NET_ADMIN` capability, use of hostPort from 1024 to 65535.
It also forces run as non-root user.

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
  users:
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
