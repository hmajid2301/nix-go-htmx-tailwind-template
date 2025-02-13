# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.1] - 2025-02-23

### Fixed

- Replace `{{service_prefix}}` with `{{project_slug}}` in gitlab CI.

## [0.2.0] - 2025-02-23

### Added

- Postgres support via `pgx`.
- Telemtry to collect metrics and traces via otel.

### Changed

- Added changes from `banterbus`, various learning after using this template.
- Moved `sqlc` files to `db` package.

## [0.1.1] - 2024-10-20

### Fixed

- Command to generate copier template.
- Not generating `.copier-answers.yml` properly.

## [0.1.0] - 2024-10-20

### Added

- Initial version released.

[unreleased]: https://gitlab.com/hmajid2301/nix-go-htmx-tailwind-template/compare/main
[0.2.1]: https://gitlab.com/hmajid2301/nix-go-htmx-tailwind-template/releases/tag/v0.2.0...v0.2.1
[0.2.0]: https://gitlab.com/hmajid2301/nix-go-htmx-tailwind-template/releases/tag/v0.1.1...v0.2.0
[0.1.1]: https://gitlab.com/hmajid2301/nix-go-htmx-tailwind-template/releases/tag/v0.1.0...v0.1.1
[0.1.0]: https://gitlab.com/hmajid2301/nix-go-htmx-tailwind-template/releases/tag/v0.1.0
