// Command showgrid prints a bordered side-by-side grid of every governa-color
// ramp at every step, each cell rendering a sample token. Useful for picking
// a label color or comparing close shades.
//
// Usage:
//
//	go run ./cmd/showgrid [-r | -r=N] [token]
//
// If [token] is omitted, "TOKEN" is used. With -r, cells use background-color
// SGR so each shade fills its rectangle as a solid block (handy for picking
// label-background colors). With -r=N (e.g. -r=15), N is the 256-color SGR
// index used as the text color rendered on top of those colored cells.
//
// Pass -h, --help, or -? to print this usage.
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/queone/governa-color"
)

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

func usage(out *os.File) {
	fmt.Fprintln(out, "Usage: showgrid [-r | -r=N] [token]")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Print a bordered side-by-side grid of every governa-color ramp")
	fmt.Fprintln(out, "at every step, each cell rendering [token] (default: TOKEN).")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Options:")
	fmt.Fprintln(out, "  -r            reverse video: cells use background-color SGR (48;5;N)")
	fmt.Fprintln(out, "                instead of foreground (38;5;N) — each cell becomes a solid")
	fmt.Fprintln(out, "                colored block. Text color falls back to terminal default.")
	fmt.Fprintln(out, "  -r=N          reverse video plus text-color override: N is a 256-color")
	fmt.Fprintln(out, "                SGR index (0-255) used for the text rendered on top of the")
	fmt.Fprintln(out, "                colored cells. Example: -r=15 for bright-white text.")
	fmt.Fprintln(out, "  -h, --help, -?")
	fmt.Fprintln(out, "                print this usage and exit")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Examples:")
	fmt.Fprintln(out, "  showgrid                # default token, foreground colors")
	fmt.Fprintln(out, "  showgrid HEADER         # custom token, foreground colors")
	fmt.Fprintln(out, "  showgrid -r             # default token, colored backgrounds, default text")
	fmt.Fprintln(out, "  showgrid -r=15 HEADER   # custom token, colored backgrounds, white text")
}

func main() {
	for _, a := range os.Args[1:] {
		if a == "-h" || a == "--help" || a == "-?" {
			usage(os.Stdout)
			return
		}
	}

	var rv reverseFlag
	rv.fgIndex = -1
	flag.Var(&rv, "r", "reverse video; -r=N sets text color to SGR index N (see -h)")
	flag.Usage = func() { usage(os.Stderr) }
	flag.Parse()

	token := ""
	if flag.NArg() > 0 {
		token = flag.Arg(0)
	}
	color.ShowGrid(token, rv.enabled, rv.fgIndex)
}
