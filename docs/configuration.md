Configuration
=============

Customize Policy
----------------

You can customize the rules for each policy.

By default, the following configuration is used:

```yaml
- name: baseline
  validators:
    - deny-host-namespace
    - deny-privileged-containers
    - deny-unsafe-capabilities
    - deny-host-path-volumes
    - deny-host-ports
    - deny-unsafe-apparmor
    - deny-unsafe-selinux
    - deny-unsafe-proc-mount
    - deny-unsafe-sysctls
  mutators: []
- name: restricted
  validators:
    - deny-non-core-volume-types
    - deny-privilege-escalation
    - deny-run-as-root
    - deny-root-groups
    - deny-unsafe-seccomp
  mutators:
    - force-run-as-non-root
```

For example, if you want to enforce `deny-run-as-root` and `force-run-as-non-root` in `Baseline`,
administrators can add the rules under the `Baseline` section: 

```yaml
- name: baseline
  validators:
    - deny-host-namespace
    - deny-privileged-containers
    - deny-unsafe-capabilities
    - deny-host-path-volumes
    - deny-host-ports
    - deny-unsafe-apparmor
    - deny-unsafe-selinux
    - deny-unsafe-proc-mount
    - deny-unsafe-sysctls
    - deny-run-as-root
  mutators:
    - force-run-as-non-root
- name: restricted
  validators:
    - deny-non-core-volume-types
    - deny-privilege-escalation
    - deny-root-groups
    - deny-unsafe-seccomp
  mutators: []
```
