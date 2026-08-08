// Package color provides ANSI 256-color terminal helpers for CLI output.
//
// v1.0 redesign: the API is built around 11-step intensity ramps per hue
// (step 0 = darkest, step 5 = canonical, step 10 = brightest) plus a
// red→yellow→green Heat ramp for severity / quality scales. Modifier
// wrappers (Bold, Reverse) compose with any color helper.
//
// All color helpers emit 256-color SGR sequences (38;5;N). Terminals that
// don't support 256-color receive uncolored output. NO_COLOR=1 and TERM=dumb
// also suppress all escapes (https://no-color.org).
package color

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ─── enable detection + SetEnabled ───────────────────────────────────────────

// enabled is true when the terminal supports color output.
var enabled = func() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}()

// SetEnabled is a test helper that mutates package-level color-enablement
// state. It is NOT safe for concurrent use — tests calling this must NOT
// call t.Parallel(). The returned closure restores the prior enabled value
// and is intended for deferred invocation.
func SetEnabled(b bool) func() {
	prev := enabled
	enabled = b
	return func() { enabled = prev }
}

// color256 is true when the terminal advertises 256-color support. The color
// helpers emit 256-color SGR only when this and enabled are both true.
var color256 = func() bool {
	ct := os.Getenv("COLORTERM")
	if ct == "truecolor" || ct == "24bit" {
		return true
	}
	return strings.Contains(os.Getenv("TERM"), "256color")
}()

// ─── ClearCode (ANSI strip) ──────────────────────────────────────────────────

// ansiSGR matches a CSI SGR escape sequence (ESC[ ... m). Used by ClearCode
// to strip color/formatting escapes from already-rendered strings or input
// piped from a color-emitting source.
var ansiSGR = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// ClearCode strips ANSI/SGR escape sequences from s. Useful for stripping
// color/formatting from previously-rendered strings or from piped/file
// input before further processing.
func ClearCode(s string) string {
	return ansiSGR.ReplaceAllString(s, "")
}

// ─── core wrap ───────────────────────────────────────────────────────────────

// wrap emits "\x1b[<code>m<v>\x1b[0m" when both enabled and color256 are
// true; otherwise returns fmt.Sprint(v) unwrapped. v1.0 dropped the
// basic-ANSI fallback layer, so non-256-color terminals receive plain text
// (matching the no-TTY / NO_COLOR / TERM=dumb suppression).
func wrap(code string, v any) string {
	s := fmt.Sprint(v)
	if !enabled || !color256 {
		return s
	}
	return "\x1b[" + code + "m" + s + "\x1b[0m"
}

// bgWrap is the bg-family counterpart to wrap. Like wrap it gates on enabled
// and color256, but additionally rewrites every inner "\x1b[0m" in the
// rendered input to "\x1b[0m\x1b[<code>m" so the bg color re-applies after
// every nested reset — matching the resetCodeRe trick Bold and Reverse use.
// Without this rewrite, BgGra2(Whi9("A") + Yel7("B")) would lose the gray
// background at the first inner reset; with it, the bg persists across the
// whole compound string. Fg helpers stay on wrap: their semantics are
// intentionally per-segment (inner fg overrides outer fg).
func bgWrap(code string, v any) string {
	s := fmt.Sprint(v)
	if !enabled || !color256 {
		return s
	}
	rewrapped := resetCodeRe.ReplaceAllString(s, "\x1b[0m\x1b["+code+"m")
	return "\x1b[" + code + "m" + rewrapped + "\x1b[0m"
}

// ─── modifier wrappers ───────────────────────────────────────────────────────

// resetCodeRe matches the literal "\x1b[0m" inner-reset emitted by every
// color helper. Bold/Reverse rewrite these so the modifier survives across
// nested color blocks (otherwise the inner reset would clear the modifier
// and only the first colored segment would render bold/reverse).
var resetCodeRe = regexp.MustCompile(`\x1b\[0m`)

// Bold wraps s with the ANSI bold attribute. Composable with hue / heat
// helpers: color.Bold(color.Grn5("text")) renders bold-green. Internal
// resets inside s are rewritten so the bold attribute survives across
// multiple colored segments — color.Bold(color.Grn5("a")+color.Red5("b"))
// renders both segments bold.
func Bold(s string) string {
	if !enabled || !color256 {
		return s
	}
	rewrapped := resetCodeRe.ReplaceAllString(s, "\x1b[0m\x1b[1m")
	return "\x1b[1m" + rewrapped + "\x1b[0m"
}

// Reverse wraps s with the ANSI reverse-video attribute (foreground and
// background swapped). Composable with hue / heat helpers; same internal-
// reset rewrite as Bold.
func Reverse(s string) string {
	if !enabled || !color256 {
		return s
	}
	rewrapped := resetCodeRe.ReplaceAllString(s, "\x1b[0m\x1b[7m")
	return "\x1b[7m" + rewrapped + "\x1b[0m"
}

// ─── hue ramps (10 hues × 11 steps = 110 funcs) ──────────────────────────────
//
// Step convention: 0 = darkest, 5 = canonical, 10 = brightest. Step 5 of any
// hue is the "default" shade (e.g., Grn5 is the canonical green). Walk
// toward 0 to dim, toward 10 to brighten. Indices are curated from the
// 256-color cube + grayscale ramp; ramps may share extreme indices across
// hues (e.g., Grn0 == Red0 == 16 = near-black) since the cube has limited
// pure-hue shades and the dim/bright extremes converge to gray/white.

// Gra — gray ramp.
func Gra0(v any) string  { return wrap("38;5;16", v) }
func Gra1(v any) string  { return wrap("38;5;233", v) }
func Gra2(v any) string  { return wrap("38;5;236", v) }
func Gra3(v any) string  { return wrap("38;5;239", v) }
func Gra4(v any) string  { return wrap("38;5;242", v) }
func Gra5(v any) string  { return wrap("38;5;245", v) }
func Gra6(v any) string  { return wrap("38;5;248", v) }
func Gra7(v any) string  { return wrap("38;5;251", v) }
func Gra8(v any) string  { return wrap("38;5;254", v) }
func Gra9(v any) string  { return wrap("38;5;255", v) }
func Gra10(v any) string { return wrap("38;5;231", v) }

// Red — red ramp.
func Red0(v any) string  { return wrap("38;5;16", v) }
func Red1(v any) string  { return wrap("38;5;52", v) }
func Red2(v any) string  { return wrap("38;5;88", v) }
func Red3(v any) string  { return wrap("38;5;124", v) }
func Red4(v any) string  { return wrap("38;5;160", v) }
func Red5(v any) string  { return wrap("38;5;196", v) }
func Red6(v any) string  { return wrap("38;5;203", v) }
func Red7(v any) string  { return wrap("38;5;210", v) }
func Red8(v any) string  { return wrap("38;5;217", v) }
func Red9(v any) string  { return wrap("38;5;224", v) }
func Red10(v any) string { return wrap("38;5;231", v) }

// Org — orange ramp.
func Org0(v any) string  { return wrap("38;5;16", v) }
func Org1(v any) string  { return wrap("38;5;88", v) }
func Org2(v any) string  { return wrap("38;5;130", v) }
func Org3(v any) string  { return wrap("38;5;166", v) }
func Org4(v any) string  { return wrap("38;5;172", v) }
func Org5(v any) string  { return wrap("38;5;208", v) }
func Org6(v any) string  { return wrap("38;5;214", v) }
func Org7(v any) string  { return wrap("38;5;215", v) }
func Org8(v any) string  { return wrap("38;5;222", v) }
func Org9(v any) string  { return wrap("38;5;229", v) }
func Org10(v any) string { return wrap("38;5;230", v) }

// Yel — yellow ramp.
func Yel0(v any) string  { return wrap("38;5;16", v) }
func Yel1(v any) string  { return wrap("38;5;58", v) }
func Yel2(v any) string  { return wrap("38;5;100", v) }
func Yel3(v any) string  { return wrap("38;5;142", v) }
func Yel4(v any) string  { return wrap("38;5;184", v) }
func Yel5(v any) string  { return wrap("38;5;220", v) }
func Yel6(v any) string  { return wrap("38;5;226", v) }
func Yel7(v any) string  { return wrap("38;5;227", v) }
func Yel8(v any) string  { return wrap("38;5;228", v) }
func Yel9(v any) string  { return wrap("38;5;229", v) }
func Yel10(v any) string { return wrap("38;5;230", v) }

// Grn — green ramp.
func Grn0(v any) string  { return wrap("38;5;16", v) }
func Grn1(v any) string  { return wrap("38;5;22", v) }
func Grn2(v any) string  { return wrap("38;5;28", v) }
func Grn3(v any) string  { return wrap("38;5;34", v) }
func Grn4(v any) string  { return wrap("38;5;40", v) }
func Grn5(v any) string  { return wrap("38;5;46", v) }
func Grn6(v any) string  { return wrap("38;5;83", v) }
func Grn7(v any) string  { return wrap("38;5;120", v) }
func Grn8(v any) string  { return wrap("38;5;157", v) }
func Grn9(v any) string  { return wrap("38;5;194", v) }
func Grn10(v any) string { return wrap("38;5;231", v) }

// Cya — cyan ramp.
func Cya0(v any) string  { return wrap("38;5;16", v) }
func Cya1(v any) string  { return wrap("38;5;23", v) }
func Cya2(v any) string  { return wrap("38;5;30", v) }
func Cya3(v any) string  { return wrap("38;5;37", v) }
func Cya4(v any) string  { return wrap("38;5;44", v) }
func Cya5(v any) string  { return wrap("38;5;51", v) }
func Cya6(v any) string  { return wrap("38;5;87", v) }
func Cya7(v any) string  { return wrap("38;5;123", v) }
func Cya8(v any) string  { return wrap("38;5;159", v) }
func Cya9(v any) string  { return wrap("38;5;195", v) }
func Cya10(v any) string { return wrap("38;5;231", v) }

// Blu — blue ramp.
func Blu0(v any) string  { return wrap("38;5;16", v) }
func Blu1(v any) string  { return wrap("38;5;17", v) }
func Blu2(v any) string  { return wrap("38;5;18", v) }
func Blu3(v any) string  { return wrap("38;5;19", v) }
func Blu4(v any) string  { return wrap("38;5;20", v) }
func Blu5(v any) string  { return wrap("38;5;21", v) }
func Blu6(v any) string  { return wrap("38;5;27", v) }
func Blu7(v any) string  { return wrap("38;5;33", v) }
func Blu8(v any) string  { return wrap("38;5;75", v) }
func Blu9(v any) string  { return wrap("38;5;111", v) }
func Blu10(v any) string { return wrap("38;5;195", v) }

// Pur — purple ramp.
func Pur0(v any) string  { return wrap("38;5;16", v) }
func Pur1(v any) string  { return wrap("38;5;53", v) }
func Pur2(v any) string  { return wrap("38;5;54", v) }
func Pur3(v any) string  { return wrap("38;5;55", v) }
func Pur4(v any) string  { return wrap("38;5;91", v) }
func Pur5(v any) string  { return wrap("38;5;92", v) }
func Pur6(v any) string  { return wrap("38;5;99", v) }
func Pur7(v any) string  { return wrap("38;5;105", v) }
func Pur8(v any) string  { return wrap("38;5;141", v) }
func Pur9(v any) string  { return wrap("38;5;177", v) }
func Pur10(v any) string { return wrap("38;5;219", v) }

// Mag — magenta ramp.
func Mag0(v any) string  { return wrap("38;5;16", v) }
func Mag1(v any) string  { return wrap("38;5;53", v) }
func Mag2(v any) string  { return wrap("38;5;90", v) }
func Mag3(v any) string  { return wrap("38;5;127", v) }
func Mag4(v any) string  { return wrap("38;5;164", v) }
func Mag5(v any) string  { return wrap("38;5;201", v) }
func Mag6(v any) string  { return wrap("38;5;207", v) }
func Mag7(v any) string  { return wrap("38;5;213", v) }
func Mag8(v any) string  { return wrap("38;5;219", v) }
func Mag9(v any) string  { return wrap("38;5;225", v) }
func Mag10(v any) string { return wrap("38;5;231", v) }

// Whi — white ramp.
func Whi0(v any) string  { return wrap("38;5;16", v) }
func Whi1(v any) string  { return wrap("38;5;59", v) }
func Whi2(v any) string  { return wrap("38;5;102", v) }
func Whi3(v any) string  { return wrap("38;5;145", v) }
func Whi4(v any) string  { return wrap("38;5;188", v) }
func Whi5(v any) string  { return wrap("38;5;231", v) }
func Whi6(v any) string  { return wrap("38;5;252", v) }
func Whi7(v any) string  { return wrap("38;5;253", v) }
func Whi8(v any) string  { return wrap("38;5;254", v) }
func Whi9(v any) string  { return wrap("38;5;255", v) }
func Whi10(v any) string { return wrap("38;5;15", v) }

// ─── heat ramp (red → yellow → green, 11 steps) ──────────────────────────────
//
// Heat0 = worst (bright red), Heat10 = best (bright green). Heat5 is the
// midpoint (yellow). Designed for severity / quality scales — e.g., skout's
// SHS column. Anchors traverse the cube's red→yellow edge (r=5, g=0..5)
// then the yellow→green edge (g=5, r=5..0).

func Heat0(v any) string  { return wrap("38;5;196", v) }
func Heat1(v any) string  { return wrap("38;5;202", v) }
func Heat2(v any) string  { return wrap("38;5;208", v) }
func Heat3(v any) string  { return wrap("38;5;214", v) }
func Heat4(v any) string  { return wrap("38;5;220", v) }
func Heat5(v any) string  { return wrap("38;5;226", v) }
func Heat6(v any) string  { return wrap("38;5;190", v) }
func Heat7(v any) string  { return wrap("38;5;154", v) }
func Heat8(v any) string  { return wrap("38;5;118", v) }
func Heat9(v any) string  { return wrap("38;5;82", v) }
func Heat10(v any) string { return wrap("38;5;46", v) }

// ─── FormatUsage ─────────────────────────────────────────────────────────────

// UsageLine is a single flag+description pair for FormatUsage.
type UsageLine struct {
	Flag string
	Desc string
}

func formatFlag(flag string) (string, int) {
	rawLen := len(flag)
	idx := strings.LastIndex(flag, " ")
	if idx < 0 {
		return flag, rawLen
	}
	suffix := flag[idx+1:]
	switch suffix {
	case "string", "int", "float", "bool", "duration":
		return flag[:idx+1] + Gra5(suffix), rawLen
	}
	return flag, rawLen
}

// FormatUsage builds a formatted help string with a heading, flag table, and optional footer.
func FormatUsage(heading string, lines []UsageLine, footer string) string {
	var b strings.Builder
	b.WriteString(Bold(Whi5("Usage:")))
	b.WriteString(" ")
	b.WriteString(heading)
	b.WriteString("\n")
	for _, l := range lines {
		flag, flagLen := formatFlag(l.Flag)
		col := 2 + flagLen
		b.WriteString("  ")
		b.WriteString(flag)
		if col < 38 {
			b.WriteString(strings.Repeat(" ", 38-col))
		} else {
			b.WriteString("  ")
		}
		b.WriteString(l.Desc)
		b.WriteString("\n")
	}
	if footer != "" {
		b.WriteString("\n")
		b.WriteString(footer)
		if !strings.HasSuffix(footer, "\n") {
			b.WriteString("\n")
		}
	}
	return b.String()
}
