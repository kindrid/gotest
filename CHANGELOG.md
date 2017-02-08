# Change Log

## Conventions
http://keepachangelog.com/en/0.3.0/

Each version should:
- **List**: its release date in the above format.
- **Group**: changes to describe their impact on the project, as follows:
- **Added**: for new features.
- **Changed**: for changes in existing functionality.
- **Deprecated**: for once-stable features removed in upcoming releases.
- **Removed**: for deprecated features removed in this release.
- **Fixed**: for any bug fixes.
- **Security**: to invite users to upgrade in case of vulnerabilities.

## [UNRELEASED]
### Changed
- Updated dependencies with `glide`
- Added `HaveOnlyCamelcaseKeys` to 'BeJSONAPIRecord' (Limitation: BeJSONAPIRecord doesn't have the `ignore` option yet to explicitly allow some snake_case fields)

## [0.9.1] 2017-02-08
### Added
- `HaveOnlyCamelcaseKeys` assertions makes sure JSON object attributes aren't snake_case.

## [0.9.0] 2017-02-07
### Added
- Vendored dependencies with `glide`
