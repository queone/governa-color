# Changelog

All notable changes to this project are documented in this file.

The format is based on Keep a Changelog and this project adheres to Semantic Versioning.

## [0.1.1] - 2026-04-29

### Fixed

- Removed a stray `(AC62)` reference in `SetEnabled`'s godoc comment. It was a governa-internal AC pointer that came along with the verbatim copy and is meaningless to library consumers. No behavior change.

## [0.1.0] - 2026-04-29

### Added

- Initial extraction from `github.com/queone/governa` `internal/color`. Verbatim move under the library policy in [governa/docs/library-policy.md](https://github.com/queone/governa/blob/main/docs/library-policy.md). No API changes from the source package.
- Color helpers: `Gra`, `Grn`, `GrnR`, `GrnD`, `Yel`, `Blu`, `Cya`, `Red`, `RedR`, `RedD`, `Whi`, `Whi2`, `BoldW`. Each emits a 256-color SGR sequence when the terminal supports it, basic ANSI otherwise.
- `SetEnabled(bool) func()` test helper for toggling color enablement.
- `ShowPalette()` for terminal verification.
- `FormatUsage(heading string, lines []UsageLine, footer string) string` for CLI help formatting; `UsageLine{Flag, Desc}` is the input pair.
