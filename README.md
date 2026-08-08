# governa-color

ANSI terminal color helpers for CLI output. Emits 256-color SGR sequences when the terminal advertises 256-color support (via `COLORTERM` or `TERM` containing `256color`); falls back to basic ANSI codes (30–97) otherwise. Suppresses colors entirely when stdout is not a terminal or `NO_COLOR` is set ([no-color.org](https://no-color.org)).

## Why

CLIs in the governa family share a small, opinionated set of color helpers — `BoldW`, `Gra`, `Red`, `Yel`, `FormatUsage`, etc. — that produce consistent terminal output across every tool a user touches in a governa-managed repo. Before extraction, every governa-family repo carried its own copy of this package and synced manually; that cost real coordination time and produced silent drift. The library exists to be the single source of truth: one set of helpers, semver-versioned, picked up via `go get -u`.

The package is leaf-clean — no governance coupling, no governa-internal types, no convention names. It works fine outside the governa family for any Go CLI that wants the same color heuristics (256-color preferred, basic ANSI fallback, NO_COLOR respected).

## Install

    go get github.com/queone/governa-color

## Usage

```go
import "github.com/queone/governa-color"

func main() {
    fmt.Println(color.Bold(color.Whi5("governa")))            // bold white
    fmt.Println(color.Gra5("dim helper text"))                // canonical gray
    fmt.Println(color.Red5("error: something went wrong"))    // canonical red
    fmt.Println(color.Heat0("worst") + " " + color.Heat10("best"))

    // Background helpers compose with foreground helpers:
    fmt.Println(color.BgGra2(color.Whi9("CAREER 2024-25")))   // light-gray text on dim-gray bg
}
```

Each hue exposes an 11-step ramp (`Gra0..Gra10`, `Red0..Red10`, …, `Whi0..Whi10`) with step 5 as the canonical shade; `Heat0..Heat10` is a severity/quality ramp (red → yellow → green). Every foreground helper has a `Bg<Hue><Step>` twin emitting the same color as background (`48;5;N` SGR), composable with `Bold`, `Reverse`, and any foreground helper.

## Palette inspection

Use the Rust [`swatch`](https://github.com/queone/rkit#swatch) utility from rkit for terminal palette inspection:

- Run `swatch palette` for the complete xterm palette and named ramps.
- Run `swatch grid [TOKEN]` for the bordered foreground-color grid.
- Run `swatch grid --reverse [TOKEN]` for the reverse-video grid.
- Add `--foreground INDEX` to a reverse grid to pin its text color.
- Run `swatch backgrounds [TOKEN]` for auto-contrasted background ramps.
- Add `--foreground INDEX` to background ramps to pin their text color.

## Versioning

This library follows the policy in [governa/docs/library-policy.md](https://github.com/queone/governa/blob/main/docs/library-policy.md). See `CHANGELOG.md` for version history and deprecations.
