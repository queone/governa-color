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

`color.ShowPalette()` prints a labeled swatch of every color function for terminal verification. `color.ShowGrid(token, reverse, fgIndex)` prints a bordered side-by-side grid of every ramp at every step (RedX..HeatX columns × 0..10\* rows), each cell rendering `token` — useful for picking a label color or comparing close shades. Empty `token` defaults to `"TOKEN"`. With `reverse=true` each cell is a solid background-color block (handy for label-background colors); `fgIndex >= 0` then overrides the text color to that 256-color SGR index (e.g. 15 for bright white).

## showgrid utility

`cmd/showgrid` is a stand-alone CLI exposing `color.ShowGrid` directly. Install:

```bash
go install github.com/queone/governa-color/cmd/showgrid@latest
```

Then run `showgrid [token]` for the default grid, `showgrid -r` for reverse-video cells (solid colored blocks), `showgrid -r=N` to override the text color to 256-color SGR index N, `showgrid -v` for the program version, or `showgrid -h` for full usage.

## Versioning

This library follows the policy in [governa/docs/library-policy.md](https://github.com/queone/governa/blob/main/docs/library-policy.md). See `CHANGELOG.md` for version history and deprecations.
