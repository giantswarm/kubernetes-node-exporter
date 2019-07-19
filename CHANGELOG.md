# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project's packages adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.5.1] 2019-07-17

### Changed

- Tolerations changed to tolerate all taints.
- Change prioty class from to `giantswarm-critical`.

## [0.4.1] 2019-06-28

### Fixed

- Fix systemd collector D-Bus connection. https://github.com/giantswarm/kubernetes-node-exporter/pull/44

## [0.4.0] 2019-06-14

### Changed

- Disabled ipvs collector.

### Fixed

- Fix monitored file system mount points.

## [0.3.0] 2019-05-24

### Changed

- Upgrade to node-exporter 0.18.0.

## [0.2.0] 2019-05-17

### Added

- Separate pod security policy for node-exporter and node-exporter-migration workloads.
- Security context with non-root user (`nobody`) for running node-exporter inside container.

[Unreleased]: https://github.com/giantswarm/kubernetes-node-exporter/compare/v0.5.0...HEAD
[Unreleased]: https://github.com/giantswarm/kubernetes-node-exporter/compare/v0.4.1...HEAD
[0.4.1]: https://github.com/giantswarm/kubernetes-node-exporter/compare/v0.4.0...v0.4.1
[0.4.0]: https://github.com/giantswarm/kubernetes-node-exporter/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/giantswarm/kubernetes-node-exporter/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/giantswarm/kubernetes-node-exporter/releases/tag/v0.2.0
