Configuration
=============

`pod-security-admission` can specify a profile for each webhook endpoint.

SecurityProfile 
----------------

SecurityProfile has these fields:

| Name                 | Type              | Description                                                |
| -------------------- | ----------------- | ---------------------------------------------------------- |
| name                 | string            | The name of the profile                                    |
| hostNamespace        | bool              | Allow sharing the host namespaces                          |
| privilegedContainers | bool              | Allow privileged containers                                |
| capabilities         | CapabilityProfile | The profile for capabilities                               |
| volumes              | VolumeProfile     | The profile for volumes                                    |
| hostPorts            | HostPortProfile   | The profile for hostPorts                                  |
| unsafeApparmor       | bool              | Allow overriding or disabling the default AppArmor profile |
| unsafeSelinux        | bool              | Allow setting custom SELinux options                       |
| unsafeProcMount      | bool              | Allow unmasked proc mount                                  |
| unsafeSysctls        | bool              | Allow usage of unsafe sysctls                              |
| privilegeEscalation  | bool              | Allow privilege escalation                                 |
| users                | UserProfile       | The profile for users                                      |
| rootGroups           | bool              | Allow running with a root primary or supplementary GID     |
| unsafeSeccomp        | bool              | Allow usage of non-default Seccomp profile                 |

### CapabilityProfile

| Name                | Type     | Description                                                                                                                                                          |
| ------------------- | -------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| unsafeCapabilities  | bool     | allow adding capabilities beyond the [default set](https://docs.docker.com/engine/reference/run/#runtime-privilege-and-linux-capabilities) and `allowedCapabilities` |
| allowedCapabilities | []string | The list of capabilities that cab be added                                                                                                                           |

### VolumeProfile

| Name               | Type | Description                          |
| ------------------ | ---- | ------------------------------------ |
| hostPathVolumes    | bool | Allow usage of HostPath volumes      |
| nonCoreVolumeTypes | bool | Allow usage of non-core volume types |

### HostPortProfile

| Name             | Type        | Description                                            |
| ---------------- | ----------- | ------------------------------------------------------ |
| hostPorts        | bool        | Allow usage of HostPorts except for `allowedHostPorts` |
| allowedHostPorts | []PortRange | The list of host ports that can be used                |

### PortRange

| Name | Type  | Description                            |
| ---- | ----- | -------------------------------------- |
| min  | int32 | The min of host port range (inclusive) |
| max  | int32 | The max of host port range (inclusive) |

### UserProfile

| Name              | Type | Description                                          |
| ----------------- | ---- | ---------------------------------------------------- |
| runAsRoot         | bool | Allow running as root users                          |
| forceRunAsNonRoot | bool | Force running with non-root users by MutatingWebhook |


Customize Profile
-----------------

By default, `pod-security-admission` uses the following configuration to enforce [Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/):
The profile is described only the items to be allowed. A false or no description indicates that the item is to be denied.


```yaml
- name: baseline
  volumes:
    nonCoreVolumeTypes: true
  privilegeEscalation: true
  users:
    runAsRoot: true
  rootGroups: true
  unsafeSeccomp: true
- name: restricted
  users:
    forceRunAsNonRoot: true
```

For example, administrators can customize the profile.
The following `Baseline` profile allows hostNamespaces, to add `SYSLOG` and `NET_ADMIN` capabilities, use of hostPorts from 1024 to 65535.

```yaml
- name: baseline
  hostNamespaces: true
  capabilities:
    allowedCapabilities:
      - SYSLOG
      - NET_ADMIN
  volumes:
    nonCoreVolumeTypes: true
  hostPorts:
    allowedHostPorts:
      - min: 1024
        max: 65535
  privilegeEscalation: true
  rootGroups: true
  users:
    runAsRoot: true
  unsafeSeccomp: true
- name: restricted
  users:
    forceRunAsNonRoot: true
```

### Webhook Configuration

You can configure [ValidatingWebhookConfiguration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#validatingwebhookconfiguration-v1-admissionregistration-k8s-io) and [MutatingWebhookConfiguration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#mutatingwebhookconfiguration-v1-admissionregistration-k8s-io) to change which Pods are covered by each policy.

See [manifests.yaml](../config/webhook/manifests.yaml) in details.

The endpoint of the webhook is determined by the name of the policy, such as `baseline` and `restricted`. 
A validating webhook will be `/validate-` + the policy name, a mutating webhook will be `/mutate-` + the policy name.
For example, the endpoint of the validating webhook for `baseline` will be `/validate-baseline`, mutating webhook will be `/mutate-baseline`.
