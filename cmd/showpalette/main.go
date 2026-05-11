// Command showpalette prints every governa-color helper plus a labeled dump
// of all 256-color SGR indices to stdout. Run from a real terminal so the
// color escapes render — when stdout is not a TTY, governa-color disables
// itself and the output is uninformative plain text.
//
// Usage:
//
//	go run ./cmd/showpalette
package main

import (
	"fmt"

	"github.com/queone/governa-color"
)

// programVersion tracks the governa-color package version. Bump in lockstep
// with the package release.
const programVersion = "1.2.1"

func main() {
	fmt.Printf("governa-color v%s\n\n", programVersion)
	color.ShowPalette()
}
