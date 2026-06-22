// Command showgrid prints a bordered side-by-side grid of every governa-color
// ramp at every step, each cell rendering a sample token. Useful for picking
// a label color or comparing close shades.
//
// Install as a stand-alone utility:
//
//	go install github.com/queone/governa-color/cmd/showgrid@latest
//
// Usage:
//
//	showgrid [-r | -r=N | -b | -b=N] [-v] [token]
//
// If [token] is omitted, "TOKEN" is used. With -r, cells use background-color
// SGR so each shade fills its rectangle as a solid block (handy for picking
// label-background colors). With -r=N (e.g. -r=15), N is the 256-color SGR
// index used as the text color rendered on top of those colored cells.
//
// With -b / --background, prints a per-family swatch row for each Bg<Hue>
// family with auto-contrasting text (white on dark steps, black on light
// steps) so the token is always legible. With -b=N, N pins the text color
// uniformly across all swatches. -r and -b are mutually exclusive.
//
// Pass -h, --help, or -? to print this usage. Pass -v or --version to print
// the program version and exit.
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/queone/governa-color"
)

// programVersion tracks the governa-color package version. Bump in lockstep
// whenever a release touches showgrid so `showgrid -v` matches the release tag
// of the surrounding module.
const programVersion = "1.5.0"

// reverseFlag implements flag.Value with IsBoolFlag, so "-r" alone enables
// reverse mode with a default-fg text color, and "-r=N" enables reverse mode
// with text color SGR index N.
type reverseFlag struct {
	enabled bool
	fgIndex int // -1 means default foreground
}

func (rv *reverseFlag) String() string {
	if rv == nil || !rv.enabled {
		return "false"
	}
	if rv.fgIndex < 0 {
		return "true"
	}
	return strconv.Itoa(rv.fgIndex)
}

func (rv *reverseFlag) Set(s string) error {
	switch s {
	case "true":
		rv.enabled = true
		rv.fgIndex = -1
		return nil
	case "false":
		rv.enabled = false
		rv.fgIndex = -1
		return nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("-r expects an integer SGR index (0-255) or no value; got %q", s)
	}
	if n < 0 || n > 255 {
		return fmt.Errorf("-r=%d out of range (must be 0-255)", n)
	}
	rv.enabled = true
	rv.fgIndex = n
	return nil
}

// IsBoolFlag tells the flag package this flag may appear without a value
// (treated as "true"). Values are still accepted via -r=N.
func (rv *reverseFlag) IsBoolFlag() bool { return true }

// bgFlag mirrors reverseFlag's flag.Value + IsBoolFlag shape for the -b /
// --background flag: bare "-b" or "--background" enables the bg-inspector
// with auto-contrasted text, "-b=N" pins the text color to SGR index N.
type bgFlag struct {
	enabled bool
	fgIndex int // -1 means auto-contrast
}

func (bg *bgFlag) String() string {
	if bg == nil || !bg.enabled {
		return "false"
	}
	if bg.fgIndex < 0 {
		return "true"
	}
	return strconv.Itoa(bg.fgIndex)
}

func (bg *bgFlag) Set(s string) error {
	switch s {
	case "true":
		bg.enabled = true
		bg.fgIndex = -1
		return nil
	case "false":
		bg.enabled = false
		bg.fgIndex = -1
		return nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("-b expects an integer SGR index (0-255) or no value; got %q", s)
	}
	if n < 0 || n > 255 {
		return fmt.Errorf("-b=%d out of range (must be 0-255)", n)
	}
	bg.enabled = true
	bg.fgIndex = n
	return nil
}

func (bg *bgFlag) IsBoolFlag() bool { return true }

func usage(out *os.File) {
	fmt.Fprintln(out, color.Whi5("Usage"))
	fmt.Fprintln(out, "  showgrid [-r | -r=N | -b | -b=N] [-v] [token]")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Print a bordered side-by-side grid of every governa-color ramp")
	fmt.Fprintln(out, "at every step, each cell rendering [token] (default: TOKEN).")
	fmt.Fprintln(out)
	fmt.Fprintln(out, color.Whi5("Options"))
	fmt.Fprintln(out, "  -r            reverse video: cells use background-color SGR (48;5;N)")
	fmt.Fprintln(out, "                instead of foreground (38;5;N) — each cell becomes a solid")
	fmt.Fprintln(out, "                colored block. Text color falls back to terminal default.")
	fmt.Fprintln(out, "  -r=N          reverse video plus text-color override: N is a 256-color")
	fmt.Fprintln(out, "                SGR index (0-255) used for the text rendered on top of the")
	fmt.Fprintln(out, "                colored cells. Example: -r=15 for bright-white text.")
	fmt.Fprintln(out, "  -b, --background")
	fmt.Fprintln(out, "                per-Bg-family swatch view: one row per Bg<Hue> family")
	fmt.Fprintln(out, "                (BgGra, BgRed, …, BgHeat) with 11 step swatches each.")
	fmt.Fprintln(out, "                Text is auto-contrasted per step for readability.")
	fmt.Fprintln(out, "                Mutually exclusive with -r.")
	fmt.Fprintln(out, "  -b=N          same swatch view with the text color pinned to 256-color")
	fmt.Fprintln(out, "                SGR index N (0-255). Example: -b=15 for bright-white text.")
	fmt.Fprintln(out, "  -v, --version")
	fmt.Fprintln(out, "                print the program version and exit")
	fmt.Fprintln(out, "  -h, --help, -?")
	fmt.Fprintln(out, "                print this usage and exit")
	fmt.Fprintln(out)
	fmt.Fprintln(out, color.Whi5("Examples"))
	fmt.Fprintln(out, "  showgrid                # default token, foreground colors")
	fmt.Fprintln(out, "  showgrid HEADER         # custom token, foreground colors")
	fmt.Fprintln(out, "  showgrid -r             # default token, colored backgrounds, default text")
	fmt.Fprintln(out, "  showgrid -r=15 HEADER   # custom token, colored backgrounds, white text")
	fmt.Fprintln(out, "  showgrid -b             # per-family bg swatches, auto-contrasted text")
	fmt.Fprintln(out, "  showgrid -b=15 HEADER   # per-family bg swatches, white text")
}

func main() {
	for _, a := range os.Args[1:] {
		if a == "-h" || a == "--help" || a == "-?" {
			usage(os.Stdout)
			return
		}
		if a == "-v" || a == "--version" {
			fmt.Printf("showgrid v%s\n", programVersion)
			return
		}
	}

	var rv reverseFlag
	rv.fgIndex = -1
	var bg bgFlag
	bg.fgIndex = -1
	flag.Var(&rv, "r", "reverse video; -r=N sets text color to SGR index N (see -h)")
	flag.Var(&bg, "b", "per-family bg swatches; -b=N pins text color to SGR index N (see -h)")
	flag.Var(&bg, "background", "long-form alias of -b")
	flag.Usage = func() { usage(os.Stderr) }
	flag.Parse()

	if rv.enabled && bg.enabled {
		fmt.Fprintln(os.Stderr, "showgrid: -r and -b are mutually exclusive; pick one")
		os.Exit(1)
	}

	token := ""
	if flag.NArg() > 0 {
		token = flag.Arg(0)
	}

	if bg.enabled {
		color.ShowBgRamps(token, bg.fgIndex)
		return
	}
	color.ShowGrid(token, rv.enabled, rv.fgIndex)
}
