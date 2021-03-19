Configuration
=============

Customize Policy
----------------

You can customize the rules for each policy.

By default, the following configuration are used:

```yaml
- name: baseline
  validators:
    - denyHostNamespace
    - denyPrivilegedContainers
    - denyCapabilities
    - denyHostPathVolumes
    - denyHostPorts
    - allowOnlyDefaultAppArmor
    - denySELinux
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

For example, if you want to enforce `denyRunAsRoot` and `mutateRunAsNonRoot` in `Baseline`, you can specify the configuration as follows:

```yaml
- name: baseline
  validators:
    - denyHostNamespace
    - denyPrivilegedContainers
    - denyCapabilities
    - denyHostPathVolumes
    - denyHostPorts
    - allowOnlyDefaultAppArmor
    - denySELinux
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
