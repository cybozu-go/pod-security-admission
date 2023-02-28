# Change Log

All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]

## [0.5.0] - 2023-02-28

### Changed

- Support Kubernetes 1.25 ([#29](https://github.com/cybozu-go/pod-security-admission/pull/29))
    - Build with go 1.20
    - Update Ubuntu to 22.04
    - Update dependencies

## [0.4.0] - 2023-01-31

### Added

- Allow net.ipv4.ip_unprivileged_port_start (#27)

## [0.3.0] - 2022-10-06

### Added

- Support for Ephemeral Container (#25)

## [0.2.4] - 2022-07-25

### Changed

- Update supported k8s version to 1.24 (#23)
- Build with Go 1.18 (#23)
- Update dependencies (#23)

## [0.2.3] - 2021-12-10

### Changed

- update supported k8s version to 1.22 (#19)

## [0.2.2] - 2021-09-17

### Changed

- Update supported k8s to 1.21 (#15)

## [0.2.1] - 2021-09-02

### Changed

- Deploy in "kube-system" namespace instead of "psa-system" (#12)

## [0.2.0] - 2021-07-20

### Added

- Add allowedHostPaths configuration (#10)

## [0.1.0] - 2021-04-08

This is the first release.

[Unreleased]: https://github.com/cybozu-go/pod-security-admission/compare/v0.5.0...HEAD
[0.5.0]: https://github.com/cybozu-go/pod-security-admission/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/cybozu-go/pod-security-admission/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/cybozu-go/pod-security-admission/compare/v0.2.4...v0.3.0
[0.2.4]: https://github.com/cybozu-go/pod-security-admission/compare/v0.2.3...v0.2.4
[0.2.3]: https://github.com/cybozu-go/pod-security-admission/compare/v0.2.2...v0.2.3
[0.2.2]: https://github.com/cybozu-go/pod-security-admission/compare/v0.2.1...v0.2.2
[0.2.1]: https://github.com/cybozu-go/pod-security-admission/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/cybozu-go/pod-security-admission/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/cybozu-go/pod-security-admission/compare/1468d8fc5862faccd4c0444b1d7721798ffe6080...v0.1.0
