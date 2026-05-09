# governa-color v1.0 Migration Guide

This document explains the v0.x → v1.0 API change and how consumers should migrate. Two parts: (1) **reference** — what changed, why, and the old→new mapping; (2) **AC template** — drop into your repo's docs/ as `ac<N>-governa-color-v1.0.md` and adapt the bracketed `<placeholders>`.

## Part 1 — Reference

### What changed

v1.0 replaces the ad-hoc named-color API (Gra, Grn, Yel, GrnD, GrnR, BoldW, …) with a systematic two-layer design:

1. **Hue ramps.** 10 named hues × 11 intensity steps each.
   - Hues: `Gra`, `Red`, `Org`, `Yel`, `Grn`, `Cya`, `Blu`, `Pur`, `Mag`, `Whi`.
   - Steps: `0` (darkest) → `5` (canonical) → `10` (brightest).
   - Helpers: `Gra0(v)`..`Gra10(v)`, `Red0(v)`..`Red10(v)`, etc. — 110 funcs total.
   - `<Hue>5` is the canonical shade for general-purpose use; walk toward 0 to dim, toward 10 to brighten.
2. **Heat ramp.** 11-step red→yellow→green gradient (`Heat0`..`Heat10`). `Heat0` is worst (bright red); `Heat10` is best (bright green); `Heat5` is the yellow midpoint. Designed for severity / quality scales (e.g., skout's SHS column).
3. **Modifier wrappers.** `Bold(s)` and `Reverse(s)` apply ANSI bold or reverse-video to any string, composable with hue / heat helpers — `Bold(Grn5("text"))` renders bold green. Internal resets are rewritten so the modifier survives across nested colored segments.

### What was dropped

- Every previous named helper: `Gra`, `Grn`, `GrnR`, `GrnD`, `Yel`, `Blu`, `Cya`, `Mag`, `Red`, `RedR`, `RedD`, `Whi`, `Whi2`, `BoldW`, `BoldR`, `BoldG`.
- The printf-style wrappers: `Yelf`, `Redf`, `Grnf`.
- The basic-ANSI (16-color) fallback layer. v1.0 emits 256-color SGR only; terminals without 256-color support receive uncolored output. (`NO_COLOR=1` and `TERM=dumb` continue to suppress all escapes per https://no-color.org.)

### What was kept (no migration needed)

- `ClearCode(s)` — strip ANSI/SGR sequences.
- `SetEnabled(b)` — test helper.
- `UsageLine` + `FormatUsage` — the help/flag formatter.
- `ShowPalette()` — now dumps the v1.0 palette (16-color + cube + grayscale + hue ramps + heat ramp). The `cmd/showpalette` runner (`go run github.com/queone/governa-color/cmd/showpalette`) is the easiest way to inspect from a terminal.

### Old → new mapping

Use this table as the starting point for a mechanical sweep. The "approximate" column flags cases where the v1.0 ramp doesn't have an exact byte-for-byte equivalent — eyeball the new render and adjust the step number if needed.

| v0.x helper          | v1.0 replacement                       | Notes                                                           |
|----------------------|----------------------------------------|-----------------------------------------------------------------|
| `color.Gra(v)`       | `color.Gra5(v)`                        | Canonical gray. Walk to `Gra3` if dimmer was wanted.            |
| `color.Grn(v)`       | `color.Grn5(v)`                        | Canonical green.                                                |
| `color.GrnR(v)`      | `color.Reverse(color.Grn5(v))`         | Composed: green text on green-ish background.                   |
| `color.GrnD(v)`      | `color.Grn3(v)` (approximate)          | Old `GrnD` was index 28; new `Grn3` is index 34. Eyeball.       |
| `color.Yel(v)`       | `color.Yel5(v)`                        | Canonical yellow.                                               |
| `color.Blu(v)`       | `color.Blu5(v)`                        | Canonical blue.                                                 |
| `color.Cya(v)`       | `color.Cya5(v)`                        | Canonical cyan.                                                 |
| `color.Mag(v)`       | `color.Mag5(v)`                        | Canonical magenta.                                              |
| `color.Red(v)`       | `color.Red5(v)`                        | Canonical red.                                                  |
| `color.RedR(v)`      | `color.Reverse(color.Red5(v))`         | Composed.                                                       |
| `color.RedD(v)`      | `color.Red3(v)` (approximate)          | Old `RedD` was index 124; new `Red3` is index 124 — exact match. |
| `color.Whi(v)`       | `color.Whi5(v)`                        | Canonical white. Walk toward `Whi9` for brighter.               |
| `color.Whi2(v)`      | `color.Whi9(v)` or `color.Whi10(v)`    | `Whi2` was bright white; `Whi10` is the 16-color "true white".  |
| `color.BoldW(v)`     | `color.Bold(color.Whi5(v))`            | Or `Bold(Whi9(v))` if you want brighter.                        |
| `color.BoldR(v)`     | `color.Bold(color.Red5(v))`            |                                                                 |
| `color.BoldG(v)`     | `color.Bold(color.Grn5(v))`            |                                                                 |
| `color.Yelf(f, ...)` | `color.Yel5(fmt.Sprintf(f, ...))`      | Same for `Redf` → `Red5`, `Grnf` → `Grn5`.                      |

### When to reach beyond step 5

Any time the call site is communicating *severity*, *intensity*, or a position on a scale, prefer the matching ramp step or the heat ramp:

- "Mostly fine but slightly off" → `Yel3` or `Heat3`/`Heat4`.
- "Very bad" → `Red5` (or `Heat0`).
- "Excellent" → `Grn5` (or `Heat10`).
- "Subtle context, don't draw the eye" → step 2-3 of any hue.
- "Stand out" → step 5-7 + optional `Bold(...)`.

### Migration procedure (every consumer)

1. **Find call sites:** `rg "queone/governa-color" <repo>` to list every file that imports it. For each, `rg "color\.\w+" <file>` to enumerate uses.
2. **Apply the mapping** above. For ambiguous cases (e.g., where the original code used the now-dropped `GrnD`), eyeball the v1.0 rendering and pick the step number that reads closest. The hue-ramp section of `go run github.com/queone/governa-color/cmd/showpalette` is the visual reference.
3. **Compose modifiers:** any place that needed bold or reverse in v0.x uses `Bold(...)` or `Reverse(...)` wrapping in v1.0.
4. **Bump the dep:** `go get github.com/queone/governa-color@v1.0.0` then `go mod tidy`.
5. **Run the consumer's full test/build pipeline** (e.g., `./build.sh` for skout). Visually spot-check any colored output the consumer emits to terminal.
6. **Update docs:** any reference to the old helper names in the consumer's docs/comments.

## Part 2 — AC Template

Copy the block below into your consumer repo as `docs/ac<N>-governa-color-v1.0.md` (or wherever your repo's AC convention is) and fill in the `<placeholders>`. Repos without AC discipline can drop the `## Status` and `## Director Review` sections.

---

```markdown
# AC<N> Migrate to governa-color v1.0

Bump `github.com/queone/governa-color` to `v1.0.0` and migrate every call site to the new API. Pure dependency-update work — no behavior change beyond the visual swap from the v0.x palette to the v1.0 hue ramps + heat ramp + modifier wrappers. See `MIGRATION-v1.0.md` in the governa-color repo for the full reference.

## Summary

v1.0 of governa-color replaces the ad-hoc named-color API with a systematic two-layer design (10 hues × 11-step intensity ramps + 11-step heat ramp + composable Bold/Reverse modifiers). This AC bumps `<repo-name>`'s dep and migrates every call site mechanically per the old→new mapping table in `MIGRATION-v1.0.md`. ECR-style external rank columns (or whatever colored output `<repo-name>` produces) keep their visual intent; only the helper names change.

## Objective Fit

1. **Outcome.** `<repo-name>` consumes governa-color v1.0; the v0.x API is fully retired across the codebase. Visual output is unchanged or improved.
2. **Priority.** Required to keep `<repo-name>` on a maintained governa-color version. The v1.0 rewrite removes all v0.x exports — staying on v0.x means freezing the dep and missing future fixes/features.
3. **Dependencies.** Builds on governa-color v1.0.0. No internal `<repo-name>` deps shift.

## In Scope

### Files to modify

- `go.mod` — bump `github.com/queone/governa-color` from `v<old-version>` to `v1.0.0`. Run `go mod tidy`.
- `<list every .go file that imports governa-color>` — apply the old→new mapping. Common transformations:
  - `color.Grn(...)` → `color.Grn5(...)`
  - `color.Red(...)` → `color.Red5(...)`
  - `color.GrnD(...)` → `color.Grn3(...)` (eyeball, may want `Grn2` for visibly darker)
  - `color.BoldW(...)` → `color.Bold(color.Whi5(...))`
  - `color.Yelf("%d", n)` → `color.Yel5(fmt.Sprintf("%d", n))`
  - …(see MIGRATION-v1.0.md for the full table)

### Documentation updates

- `<list any docs that reference helper names>` — e.g., glossary entries, README, comments mentioning specific helper names.

### Out Of Scope

- Adopting heat-ramp or hue-ramp granularity beyond what the v0.x API provided. (E.g., `<repo-name>` had two reds: `Red` and `RedD`. The migration keeps using two reds — `Red5` and `Red3`. Adding `Red7`, `Red9`, etc. to new call sites is a *follow-up* AC.)
- Changing visual intent. Where the migration is approximate (`GrnD` → `Grn3`), the goal is "looks the same as before"; if the new rendering looks notably different, file a follow-up to re-pick the step.
- Refactoring callers' color usage patterns beyond the mechanical name swap.

## Implementation Notes

- The v1.0 module emits 256-color SGR only — there is no 16-color fallback. If `<repo-name>` is ever run on a non-256-color terminal, output is uncolored. (Modern terminals all support 256-color; `NO_COLOR=1` and `TERM=dumb` continue to suppress as before.)
- `Bold` and `Reverse` are wrappers that compose with any hue/heat helper. `Bold(Grn5("x"))` is the v1.0 way to render bold green.
- The exact v0.x → v1.0 step numbers in the mapping table are starter picks. For any call site where the visual result looks off, the fix is to bump the step number up or down. The `cmd/showpalette` runner is the easiest way to compare side-by-side.

## Acceptance Tests

**AT1** [Automated] — `grep -rn "color\.\(Gra\|Grn\|Yel\|Blu\|Cya\|Mag\|Red\|Whi\|Whi2\|GrnR\|GrnD\|RedR\|RedD\|BoldW\|BoldR\|BoldG\|Yelf\|Redf\|Grnf\)\b" <source-tree>` returns zero matches. Verifies no v0.x names survive.

**AT2** [Automated] — `go.mod` lists `github.com/queone/governa-color v1.0.0` (or later). `go.sum` updated.

**AT3** [Automated] — `<repo-name>`'s canonical build pipeline (`<build-command>`) passes — vet, staticcheck, all tests, build.

**AT4** [Manual] [Post-release verification] — Spot-check `<list 2-3 representative commands or surfaces that emit colored output>` in a terminal. Confirm the visual output is recognizable / improved relative to pre-migration.

## Documentation Updates

- `arch.md` (if applicable) — note the dep version bump.
- `<other docs>` — update any literal helper-name references.
- `CHANGELOG.md` — release-prep, not implementation.

## Director Review

None.

## Status

`PENDING` — awaiting Director critique.
```

---

### Notes on adapting the template

- **Skout:** has full AC discipline (`docs/ac-template.md`, `docs/critique-protocol.md`); follow that flow. Drop the AC at `docs/ac<N>-governa-color-v1.0.md` and `IE`-list it in `plan.md` per skout's convention.
- **`governa` (umbrella for `governa-buildtool`, `governa-reltool`):** if those two are vendored or worked on from inside `governa`, write one combined AC at `governa/docs/...` that scopes all three modules' migrations together.
- **Light-governance repos (`utils`, `iq`):** drop the `## Status` and `## Director Review` sections; the body of the template still serves as a checklist. A simple commit message + PR description referencing `MIGRATION-v1.0.md` is fine.
