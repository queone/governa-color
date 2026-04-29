# governa-color

ANSI terminal color helpers for CLI output. Emits 256-color SGR sequences when the terminal advertises 256-color support (via `COLORTERM` or `TERM` containing `256color`); falls back to basic ANSI codes (30–97) otherwise. Suppresses colors entirely when stdout is not a terminal or `NO_COLOR` is set ([no-color.org](https://no-color.org)).

## Install

    go get github.com/queone/governa-color

## Usage

```go
import "github.com/queone/governa-color"

func main() {
    fmt.Println(color.BoldW("governa"))
    fmt.Println(color.Gra("dim helper text"))
    fmt.Println(color.Red("error: something went wrong"))
}
```

`color.ShowPalette()` prints a labeled swatch of every color function for terminal verification.

## Versioning

This library follows the policy in [governa/docs/library-policy.md](https://github.com/queone/governa/blob/main/docs/library-policy.md). See `CHANGELOG.md` for version history and deprecations.
