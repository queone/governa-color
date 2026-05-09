// Command showpalette prints every governa-color helper plus a labeled dump
// of all 256-color SGR indices to stdout. Run from a real terminal so the
// color escapes render — when stdout is not a TTY, governa-color disables
// itself and the output is uninformative plain text.
//
// Usage:
//
//	go run ./cmd/showpalette
package main

import "github.com/queone/governa-color"

func main() {
	color.ShowPalette()
}
