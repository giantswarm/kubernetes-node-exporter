# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project's packages adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [v0.2.0]

### Added

- Separate pod security policy for node-exporter and node-exporter-migration workloads.
- Security context with non-root user (`nobody`) for running node-exporter inside container.


