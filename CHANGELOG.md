# GoTest Change Log

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

## [0.9.4] UNRELEASED

## [0.9.3] 2017-02-16
### Added
- Interface `should.StructureExplorer`, a minimal JSON destructuring interface to decouple libraries using `gotest` from `github.com/Jeffail/gabs`.
- Method `should.ParseJSON()` that returns a `should.StructureExplorer` so outside libraries can write their own JSON assertions.
- method `gotest.Later()` to sketch out unimplemented tests

## [0.9.2] 2017-02-09
### Changed
- Updated dependencies with `glide`
- Added `HaveOnlyCamelcaseKeys` to 'BeJSONAPIRecord' (Limitation: BeJSONAPIRecord doesn't have the `ignore` option yet to explicitly allow some snake_case fields)

### Fixed
- camelCase detection regexp now allows numbers but requires initial lowercase letter

## [0.9.1] 2017-02-08
### Added
- `HaveOnlyCamelcaseKeys` assertions makes sure JSON object attributes aren't snake_case.

## [0.9.0] 2017-02-07
### Added
- Vendored dependencies with `glide`
