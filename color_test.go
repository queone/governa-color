package color

import (
	"io"
	"os"
	"strings"
	"testing"
)

// In test environments stdout is not a TTY, so enabled == false and all
// color functions return the input string unchanged. The tests below verify
// both the no-color path (input preserved) and the wrap helper directly.
//
// No test in this file uses t.Parallel() because several tests mutate the
// package-level enabled and color256 variables. Running any test concurrently
// with those mutations would be a data race.

// allHelpers returns a labeled list of every named color helper in v1.0.
// Used by tests to drive uniform property checks.
func allHelpers() []struct {
	name string
	fn   func(any) string
} {
	return []struct {
		name string
		fn   func(any) string
	}{
		{"Gra0", Gra0}, {"Gra1", Gra1}, {"Gra2", Gra2}, {"Gra3", Gra3}, {"Gra4", Gra4},
		{"Gra5", Gra5}, {"Gra6", Gra6}, {"Gra7", Gra7}, {"Gra8", Gra8}, {"Gra9", Gra9}, {"Gra10", Gra10},
		{"Red0", Red0}, {"Red1", Red1}, {"Red2", Red2}, {"Red3", Red3}, {"Red4", Red4},
		{"Red5", Red5}, {"Red6", Red6}, {"Red7", Red7}, {"Red8", Red8}, {"Red9", Red9}, {"Red10", Red10},
		{"Org0", Org0}, {"Org1", Org1}, {"Org2", Org2}, {"Org3", Org3}, {"Org4", Org4},
		{"Org5", Org5}, {"Org6", Org6}, {"Org7", Org7}, {"Org8", Org8}, {"Org9", Org9}, {"Org10", Org10},
		{"Yel0", Yel0}, {"Yel1", Yel1}, {"Yel2", Yel2}, {"Yel3", Yel3}, {"Yel4", Yel4},
		{"Yel5", Yel5}, {"Yel6", Yel6}, {"Yel7", Yel7}, {"Yel8", Yel8}, {"Yel9", Yel9}, {"Yel10", Yel10},
		{"Grn0", Grn0}, {"Grn1", Grn1}, {"Grn2", Grn2}, {"Grn3", Grn3}, {"Grn4", Grn4},
		{"Grn5", Grn5}, {"Grn6", Grn6}, {"Grn7", Grn7}, {"Grn8", Grn8}, {"Grn9", Grn9}, {"Grn10", Grn10},
		{"Cya0", Cya0}, {"Cya1", Cya1}, {"Cya2", Cya2}, {"Cya3", Cya3}, {"Cya4", Cya4},
		{"Cya5", Cya5}, {"Cya6", Cya6}, {"Cya7", Cya7}, {"Cya8", Cya8}, {"Cya9", Cya9}, {"Cya10", Cya10},
		{"Blu0", Blu0}, {"Blu1", Blu1}, {"Blu2", Blu2}, {"Blu3", Blu3}, {"Blu4", Blu4},
		{"Blu5", Blu5}, {"Blu6", Blu6}, {"Blu7", Blu7}, {"Blu8", Blu8}, {"Blu9", Blu9}, {"Blu10", Blu10},
		{"Pur0", Pur0}, {"Pur1", Pur1}, {"Pur2", Pur2}, {"Pur3", Pur3}, {"Pur4", Pur4},
		{"Pur5", Pur5}, {"Pur6", Pur6}, {"Pur7", Pur7}, {"Pur8", Pur8}, {"Pur9", Pur9}, {"Pur10", Pur10},
		{"Mag0", Mag0}, {"Mag1", Mag1}, {"Mag2", Mag2}, {"Mag3", Mag3}, {"Mag4", Mag4},
		{"Mag5", Mag5}, {"Mag6", Mag6}, {"Mag7", Mag7}, {"Mag8", Mag8}, {"Mag9", Mag9}, {"Mag10", Mag10},
		{"Whi0", Whi0}, {"Whi1", Whi1}, {"Whi2", Whi2}, {"Whi3", Whi3}, {"Whi4", Whi4},
		{"Whi5", Whi5}, {"Whi6", Whi6}, {"Whi7", Whi7}, {"Whi8", Whi8}, {"Whi9", Whi9}, {"Whi10", Whi10},
		{"Heat0", Heat0}, {"Heat1", Heat1}, {"Heat2", Heat2}, {"Heat3", Heat3}, {"Heat4", Heat4},
		{"Heat5", Heat5}, {"Heat6", Heat6}, {"Heat7", Heat7}, {"Heat8", Heat8}, {"Heat9", Heat9}, {"Heat10", Heat10},
	}
}

// TestColorFunctionsContainInput verifies every named helper preserves the
// input string. With color enabled, the input is wrapped in escape codes;
// with color disabled, the input is returned bare. Either way, "hello" must
// appear somewhere in the output.
func TestColorFunctionsContainInput(t *testing.T) {
	for _, tc := range allHelpers() {
		got := tc.fn("hello")
		if !strings.Contains(got, "hello") {
			t.Errorf("%s(%q) = %q, does not contain input", tc.name, "hello", got)
		}
		if got == "" {
			t.Errorf("%s(%q) returned empty string", tc.name, "hello")
		}
	}
}

// TestColorFunctionsNoTTY confirms the no-color path returns bare input.
func TestColorFunctionsNoTTY(t *testing.T) {
	if enabled {
		t.Skip("TTY detected — skipping no-color path test")
	}
	for _, tc := range allHelpers() {
		got := tc.fn("test")
		if got != "test" {
			t.Errorf("%s(%q) = %q, want %q (no-TTY path)", tc.name, "test", got, "test")
		}
	}
}

func TestWrapEmptyString(t *testing.T) {
	_ = wrap("32", "")
}

// TestWrapProducesEscapes verifies wrap() emits the documented ANSI shape.
func TestWrapProducesEscapes(t *testing.T) {
	defer SetEnabled(true)()

	got := wrap("38;5;46", "ok")
	want := "\x1b[38;5;46mok\x1b[0m"
	if got != want {
		t.Fatalf("wrap(\"38;5;46\", \"ok\") = %q, want %q", got, want)
	}
}

// TestColorFunctions256Codes spot-checks the documented 256-color escape
// for canonical (step 5) hues and a few endpoints.
func TestColorFunctions256Codes(t *testing.T) {
	defer SetEnabled(true)()

	cases := []struct {
		name string
		fn   func(any) string
		code string
	}{
		{"Gra5", Gra5, "38;5;245"},
		{"Red5", Red5, "38;5;196"},
		{"Org5", Org5, "38;5;208"},
		{"Yel5", Yel5, "38;5;220"},
		{"Grn5", Grn5, "38;5;46"},
		{"Cya5", Cya5, "38;5;51"},
		{"Blu5", Blu5, "38;5;21"},
		{"Pur5", Pur5, "38;5;92"},
		{"Mag5", Mag5, "38;5;201"},
		{"Whi5", Whi5, "38;5;231"},
		{"Heat0", Heat0, "38;5;196"},
		{"Heat5", Heat5, "38;5;226"},
		{"Heat10", Heat10, "38;5;46"},
		{"Grn0", Grn0, "38;5;16"},
		{"Grn10", Grn10, "38;5;231"},
	}
	for _, tc := range cases {
		got := tc.fn("x")
		wantPrefix := "\x1b[" + tc.code + "m"
		if !strings.HasPrefix(got, wantPrefix) {
			t.Errorf("%s: got %q, want prefix %q", tc.name, got, wantPrefix)
		}
		if !strings.HasSuffix(got, "\x1b[0m") {
			t.Errorf("%s: got %q, want suffix %q", tc.name, got, "\x1b[0m")
		}
	}
}

// TestBoldComposesWithSingleColor verifies Bold(Grn5("x")) renders bold
// green: bold attribute prefix, color escape, content, single-reset suffix.
func TestBoldComposesWithSingleColor(t *testing.T) {
	defer SetEnabled(true)()

	got := Bold(Grn5("x"))
	want := "\x1b[1m\x1b[38;5;46mx\x1b[0m\x1b[1m\x1b[0m"
	if got != want {
		t.Errorf("Bold(Grn5(\"x\")) = %q, want %q", got, want)
	}
}

// TestBoldComposesWithMultipleColors verifies Bold's reset-rewrite makes the
// modifier survive across nested color blocks.
func TestBoldComposesWithMultipleColors(t *testing.T) {
	defer SetEnabled(true)()

	got := Bold(Grn5("a") + Red5("b"))
	// The inner reset between segments is rewritten so bold re-applies for "b".
	if !strings.Contains(got, "\x1b[1m\x1b[38;5;46ma\x1b[0m\x1b[1m\x1b[38;5;196mb\x1b[0m") {
		t.Errorf("Bold composition didn't rewrap inner resets: got %q", got)
	}
}

// TestReverseComposesWithColor — Reverse mirrors Bold's behavior.
func TestReverseComposesWithColor(t *testing.T) {
	defer SetEnabled(true)()

	got := Reverse(Grn5("x"))
	want := "\x1b[7m\x1b[38;5;46mx\x1b[0m\x1b[7m\x1b[0m"
	if got != want {
		t.Errorf("Reverse(Grn5(\"x\")) = %q, want %q", got, want)
	}
}

// TestBoldNoTTY: when color is disabled, Bold returns input unchanged.
func TestBoldNoTTY(t *testing.T) {
	if enabled {
		t.Skip("TTY detected — skipping no-color path test")
	}
	if got := Bold("hello"); got != "hello" {
		t.Errorf("Bold(%q) = %q, want %q", "hello", got, "hello")
	}
	if got := Reverse("hello"); got != "hello" {
		t.Errorf("Reverse(%q) = %q, want %q", "hello", got, "hello")
	}
}

// TestShowPaletteCoversAllSections captures ShowPalette output and checks
// every documented section header is present.
func TestShowPaletteCoversAllSections(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	ShowPalette()
	w.Close()
	os.Stdout = oldStdout

	buf, _ := io.ReadAll(r)
	output := string(buf)

	for _, header := range []string{
		"Color palette v1.0",
		"Standard 16-color palette",
		"256-color cube",
		"Grayscale ramp",
		"Hue ramps",
		"Heat ramp",
	} {
		if !strings.Contains(output, header) {
			t.Errorf("ShowPalette() output missing section %q", header)
		}
	}

	// Spot-check hue labels and heat-ramp indices appear.
	for _, label := range []string{"Gra ", "Red ", "Org ", "Yel ", "Grn ", "Cya ", "Blu ", "Pur ", "Mag ", "Whi "} {
		if !strings.Contains(output, label) {
			t.Errorf("ShowPalette() output missing hue label %q", label)
		}
	}
	for _, idx := range []string{"196", "226", "46"} {
		if !strings.Contains(output, idx) {
			t.Errorf("ShowPalette() output missing heat-ramp index %q", idx)
		}
	}
}

func TestFormatUsage(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		got := FormatUsage("prog [flags]", []UsageLine{
			{"-v", "verbose output"},
			{"-o string", "output file"},
		}, "")
		if !strings.HasPrefix(got, "Usage: prog [flags]\n") {
			t.Errorf("heading mismatch: %q", got)
		}
		if !strings.Contains(got, "-v") || !strings.Contains(got, "verbose output") {
			t.Errorf("missing flag line: %q", got)
		}
		if !strings.Contains(got, "-o string") || !strings.Contains(got, "output file") {
			t.Errorf("missing type-suffix flag line: %q", got)
		}
		if strings.HasSuffix(got, "\n\n") {
			t.Errorf("unexpected trailing blank line with empty footer: %q", got)
		}
	})

	t.Run("footer_no_newline", func(t *testing.T) {
		got := FormatUsage("prog", nil, "See docs.")
		if !strings.Contains(got, "\nSee docs.\n") {
			t.Errorf("footer missing or not newline-terminated: %q", got)
		}
	})

	t.Run("footer_with_newline", func(t *testing.T) {
		got := FormatUsage("prog", nil, "See docs.\n")
		if strings.HasSuffix(got, "docs.\n\n") {
			t.Errorf("footer double-newlined: %q", got)
		}
	})

	t.Run("long_flag", func(t *testing.T) {
		got := FormatUsage("prog", []UsageLine{
			{"--very-long-flag-name-that-exceeds-column string", "desc"},
		}, "")
		if !strings.Contains(got, "string  desc") {
			t.Errorf("long flag alignment wrong: %q", got)
		}
	})
}

func TestClearCodeStripsAnsiSGR(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"\x1b[32mhello\x1b[0m", "hello"},
		{"plain text", "plain text"},
		{"\x1b[38;5;46mok\x1b[0m end", "ok end"},
		{"", ""},
		{"\x1b[1;38;5;15mbold\x1b[0m and \x1b[31mred\x1b[0m", "bold and red"},
	}
	for _, tc := range cases {
		if got := ClearCode(tc.in); got != tc.want {
			t.Errorf("ClearCode(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// TestClearCodeRoundTripWithColorHelpers renders a sample through each
// canonical (step 5) hue helper plus Bold/Reverse, then verifies ClearCode
// strips back to the original.
func TestClearCodeRoundTripWithColorHelpers(t *testing.T) {
	defer SetEnabled(true)()

	sample := "the quick brown fox"
	helpers := []struct {
		name string
		fn   func(any) string
	}{
		{"Gra5", Gra5}, {"Red5", Red5}, {"Org5", Org5}, {"Yel5", Yel5},
		{"Grn5", Grn5}, {"Cya5", Cya5}, {"Blu5", Blu5}, {"Pur5", Pur5},
		{"Mag5", Mag5}, {"Whi5", Whi5}, {"Heat5", Heat5},
	}
	for _, h := range helpers {
		rendered := h.fn(sample)
		if got := ClearCode(rendered); got != sample {
			t.Errorf("ClearCode(%s(%q)) = %q, want %q", h.name, sample, got, sample)
		}
	}

	// Composed Bold + color also strips clean.
	rendered := Bold(Grn5(sample))
	if got := ClearCode(rendered); got != sample {
		t.Errorf("ClearCode(Bold(Grn5(%q))) = %q, want %q", sample, got, sample)
	}
}

// TestColorSetEnabledTogglesAndRestores: SetEnabled flips the package-level
// enabled var and the returned closure restores it.
func TestColorSetEnabledTogglesAndRestores(t *testing.T) {
	before := enabled

	restore := SetEnabled(true)
	if !enabled {
		t.Errorf("after SetEnabled(true): enabled = false, want true")
	}
	restore()
	if enabled != before {
		t.Errorf("after restore: enabled = %v, want %v", enabled, before)
	}

	restore = SetEnabled(false)
	if enabled {
		t.Errorf("after SetEnabled(false): enabled = true, want false")
	}
	restore()
	if enabled != before {
		t.Errorf("after second restore: enabled = %v, want %v", enabled, before)
	}
}
