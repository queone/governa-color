# governa-color Architecture

## Purpose

Provide reusable ANSI 256-color helpers for Go command-line applications.

## System Summary

The repository is a Go library with no installable command binaries. It exposes foreground and background ramps, composable modifiers, ANSI clearing, terminal capability detection, and usage formatting.

## Current Platform

- Go

## Major Components

- foreground, background, and heat-ramp helpers
- terminal capability and color-enablement detection
- composable bold and reverse modifiers
- ANSI clearing and usage formatting

## Core Files

- `AGENTS.md`: base governance contract
- `plan.md`: prioritized roadmap and approved direction
- `build.sh`: self-contained build, release-prep, and release tooling
- `govna/development-cycle.md`: workflow from roadmap through release
- `govna/ac-template.md`: acceptance-criteria template for new work
- `govna/build-release.md`: build, test, and release rules
- `docs/critique-protocol.md`: critique protocol (repo-specific; not a governa doc)

## Data And Control Flow

Consumers call an exported helper. The package checks terminal and environment capabilities, emits the applicable SGR wrapper, and otherwise returns uncolored text.

## Architecture Notes

- record stable system decisions here
- prefer durable structure and interfaces over transient implementation detail

## Conventions

- update this document when architecture or major workflow changes materially
- keep implementation detail in code and stable architecture here
