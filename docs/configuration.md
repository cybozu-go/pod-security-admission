Configuration
=============

Customize Policy
----------------

You can customize the rules for each policy.

By default, the following configuration is used:

```yaml
- name: baseline
  validators:
    - denyHostNamespace
    - denyPrivilegedContainers
    - denyBeyondDefaultCapabilities
    - denyHostPathVolumes
    - denyHostPorts
    - allowOnlyDefaultAppArmor
    - denyCustomSELinux
    - allowOnlyDefaultProcMount
    - allowOnlySafeSysctls
  mutators: []
- name: restricted
  validators:
    - denyNonCoreVolumeTypes
    - denyPrivilegeEscalation
    - denyRunAsRoot
    - denyRootGroups
    - allowOnlyDefaultSeccomp
  mutators:
    - mutateRunAsNonRoot
```

For example, if you want to enforce `denyRunAsRoot` and `mutateRunAsNonRoot` in `Baseline`,
administrators can add the rules under the `Baseline` section: 

```yaml
- name: baseline
  validators:
    - denyHostNamespace
    - denyPrivilegedContainers
    - denyBeyondDefaultCapabilities
    - denyHostPathVolumes
    - denyHostPorts
    - allowOnlyDefaultAppArmor
    - denyCustomSELinux
    - allowOnlyDefaultProcMount
    - allowOnlySafeSysctls
    - denyRunAsRoot
  mutators:
    - mutateRunAsNonRoot
- name: restricted
  validators:
    - denyNonCoreVolumeTypes
    - denyPrivilegeEscalation
    - denyRootGroups
    - allowOnlyDefaultSeccomp
  mutators: []
```
