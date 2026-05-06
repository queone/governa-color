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

func TestColorFunctionsContainInput(t *testing.T) {
	cases := []struct {
		name string
		fn   func(any) string
	}{
		{"Gra", Gra},
		{"Grn", Grn},
		{"GrnR", GrnR},
		{"GrnD", GrnD},
		{"Yel", Yel},
		{"Blu", Blu},
		{"Cya", Cya},
		{"Mag", Mag},
		{"Red", Red},
		{"RedR", RedR},
		{"RedD", RedD},
		{"Whi", Whi},
		{"Whi2", Whi2},
		{"BoldW", BoldW},
	}
	for _, tc := range cases {
		got := tc.fn("hello")
		if !strings.Contains(got, "hello") {
			t.Errorf("%s(%q) = %q, does not contain input", tc.name, "hello", got)
		}
		if got == "" {
			t.Errorf("%s(%q) returned empty string", tc.name, "hello")
		}
	}
}

func TestColorFunctionsNoTTY(t *testing.T) {
	// In test environment, enabled is false (stdout is not a char device).
	// Functions must return the bare input string.
	if enabled {
		t.Skip("TTY detected — skipping no-color path test")
	}
	cases := []struct {
		name string
		fn   func(any) string
	}{
		{"Gra", Gra},
		{"Grn", Grn},
		{"GrnR", GrnR},
		{"GrnD", GrnD},
		{"Yel", Yel},
		{"Blu", Blu},
		{"Cya", Cya},
		{"Mag", Mag},
		{"Red", Red},
		{"RedR", RedR},
		{"RedD", RedD},
		{"Whi", Whi},
		{"Whi2", Whi2},
		{"BoldW", BoldW},
	}
	for _, tc := range cases {
		got := tc.fn("test")
		if got != "test" {
			t.Errorf("%s(%q) = %q, want %q (no-TTY path)", tc.name, "test", got, "test")
		}
	}
}

func TestWrapEmptyString(t *testing.T) {
	// wrap with an empty input should not panic regardless of TTY state.
	_ = wrap("32", "")
}

// TestWrapProduces256ColorEscapes verifies wrap() emits the exact ANSI
// 256-color escape format. This test calls wrap() directly so results are
// deterministic regardless of TTY state.
func TestWrapProduces256ColorEscapes(t *testing.T) {
	origEnabled := enabled
	enabled = true
	defer func() { enabled = origEnabled }()

	got := wrap("38;5;2", "ok")
	want := "\033[38;5;2mok\033[0m"
	if got != want {
		t.Fatalf("wrap(\"38;5;2\", \"ok\") = %q, want %q", got, want)
	}
}

// TestColorFunctions256Codes verifies every color function uses the
// documented 256-color escape code when color256 is true.
func TestColorFunctions256Codes(t *testing.T) {
	origEnabled := enabled
	orig256 := color256
	enabled = true
	color256 = true
	defer func() { enabled = origEnabled; color256 = orig256 }()

	cases := []struct {
		name string
		fn   func(any) string
		code string // expected escape code between \033[ and m
	}{
		{"Gra", Gra, "38;5;246"},
		{"Grn", Grn, "38;5;2"},
		{"GrnR", GrnR, "7;38;5;2"},
		{"GrnD", GrnD, "38;5;28"},
		{"Yel", Yel, "38;5;3"},
		{"Blu", Blu, "38;5;12"},
		{"Cya", Cya, "38;5;6"},
		{"Mag", Mag, "38;5;13"},
		{"Red", Red, "38;5;9"},
		{"RedR", RedR, "38;5;15;48;5;1"},
		{"RedD", RedD, "38;5;124"},
		{"Whi", Whi, "38;5;7"},
		{"Whi2", Whi2, "38;5;15"},
		{"BoldW", BoldW, "1;38;5;15"},
	}
	for _, tc := range cases {
		got := tc.fn("x")
		wantPrefix := "\033[" + tc.code + "m"
		if !strings.HasPrefix(got, wantPrefix) {
			t.Errorf("%s: got %q, want prefix %q", tc.name, got, wantPrefix)
		}
		wantSuffix := "\033[0m"
		if !strings.HasSuffix(got, wantSuffix) {
			t.Errorf("%s: got %q, want suffix %q", tc.name, got, wantSuffix)
		}
	}
}

// TestColorFunctionsBasicCodes verifies every color function falls back to
// basic ANSI codes when color256 is false.
func TestColorFunctionsBasicCodes(t *testing.T) {
	origEnabled := enabled
	orig256 := color256
	enabled = true
	color256 = false
	defer func() { enabled = origEnabled; color256 = orig256 }()

	cases := []struct {
		name string
		fn   func(any) string
		code string
	}{
		{"Gra", Gra, "90"},
		{"Grn", Grn, "32"},
		{"GrnR", GrnR, "7;32"},
		{"GrnD", GrnD, "32"},
		{"Yel", Yel, "33"},
		{"Blu", Blu, "94"},
		{"Cya", Cya, "36"},
		{"Mag", Mag, "95"},
		{"Red", Red, "91"},
		{"RedR", RedR, "97;41"},
		{"RedD", RedD, "31"},
		{"Whi", Whi, "37"},
		{"Whi2", Whi2, "97"},
		{"BoldW", BoldW, "1;97"},
	}
	for _, tc := range cases {
		got := tc.fn("x")
		wantPrefix := "\033[" + tc.code + "m"
		if !strings.HasPrefix(got, wantPrefix) {
			t.Errorf("%s: got %q, want prefix %q", tc.name, got, wantPrefix)
		}
		wantSuffix := "\033[0m"
		if !strings.HasSuffix(got, wantSuffix) {
			t.Errorf("%s: got %q, want suffix %q", tc.name, got, wantSuffix)
		}
	}
}

// TestShowPaletteCoversAllFunctions captures ShowPalette output and verifies
// all 13 color function labels are present.
func TestShowPaletteCoversAllFunctions(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	ShowPalette()
	w.Close()
	os.Stdout = oldStdout

	buf, _ := io.ReadAll(r)
	output := string(buf)

	for _, label := range []string{
		"Gra", "Grn", "GrnR", "GrnD",
		"Yel", "Blu", "Cya", "Mag",
		"Red", "RedR", "RedD",
		"Whi", "Whi2", "BoldW",
	} {
		if !strings.Contains(output, label) {
			t.Errorf("ShowPalette() output missing label %q", label)
		}
	}
}

// TestFormatUsage exercises heading, flag alignment, type-suffix rendering,
// long-flag overflow, and footer newline handling.
func TestFormatUsage(t *testing.T) {
	// Tests run with enabled=false, so color wrappers return plain text.

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
		// No footer => no trailing blank line.
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
		// Should not double the trailing newline.
		if strings.HasSuffix(got, "docs.\n\n") {
			t.Errorf("footer double-newlined: %q", got)
		}
	})

	t.Run("long_flag", func(t *testing.T) {
		got := FormatUsage("prog", []UsageLine{
			{"--very-long-flag-name-that-exceeds-column string", "desc"},
		}, "")
		// Long flags get 2-space gap instead of padding to column 38.
		if !strings.Contains(got, "string  desc") {
			t.Errorf("long flag alignment wrong: %q", got)
		}
	})
}

// TestClearCodeStripsAnsiSGR feeds known ANSI escape sequences and verifies
// they are removed while non-escape content is preserved.
func TestClearCodeStripsAnsiSGR(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"\x1b[32mhello\x1b[0m", "hello"},
		{"plain text", "plain text"},
		{"\x1b[38;5;2mok\x1b[0m end", "ok end"},
		{"", ""},
		{"\x1b[1;38;5;15mbold\x1b[0m and \x1b[31mred\x1b[0m", "bold and red"},
		{"no escapes here", "no escapes here"},
	}
	for _, tc := range cases {
		if got := ClearCode(tc.in); got != tc.want {
			t.Errorf("ClearCode(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// TestClearCodeRoundTripWithColorHelpers renders a sample through every
// color helper and verifies ClearCode strips back to the original input.
func TestClearCodeRoundTripWithColorHelpers(t *testing.T) {
	defer SetEnabled(true)()

	sample := "the quick brown fox"
	helpers := []struct {
		name string
		fn   func(any) string
	}{
		{"Gra", Gra},
		{"Grn", Grn},
		{"GrnR", GrnR},
		{"GrnD", GrnD},
		{"Yel", Yel},
		{"Blu", Blu},
		{"Cya", Cya},
		{"Mag", Mag},
		{"Red", Red},
		{"RedR", RedR},
		{"RedD", RedD},
		{"Whi", Whi},
		{"Whi2", Whi2},
		{"BoldW", BoldW},
	}
	for _, h := range helpers {
		rendered := h.fn(sample)
		if got := ClearCode(rendered); got != sample {
			t.Errorf("ClearCode(%s(%q)) = %q, want %q", h.name, sample, got, sample)
		}
	}
}

// TestPrintfStyleColorWrappers verifies Yelf/Redf/Grnf format their input
// through fmt.Sprintf and pass the result through their non-printf
// counterpart for coloring.
func TestPrintfStyleColorWrappers(t *testing.T) {
	defer SetEnabled(true)()

	cases := []struct {
		name    string
		fn      func(string, ...any) string
		wrapper func(any) string
	}{
		{"Yelf", Yelf, Yel},
		{"Redf", Redf, Red},
		{"Grnf", Grnf, Grn},
	}
	for _, tc := range cases {
		got := tc.fn("count=%d name=%s", 42, "alice")
		want := tc.wrapper("count=42 name=alice")
		if got != want {
			t.Errorf("%s(%q, 42, \"alice\") = %q, want %q", tc.name, "count=%d name=%s", got, want)
		}
	}
}

// SetEnabled toggles the package-level enabled var and returns a
// restore closure. Not parallel-safe by design.
func TestColorSetEnabledTogglesAndRestores(t *testing.T) {
	// Record the current state so we can confirm round-trip.
	before := enabled

	restore := SetEnabled(true)
	if !enabled {
		t.Errorf("after SetEnabled(true): enabled = false, want true")
	}
	restore()
	if enabled != before {
		t.Errorf("after restore: enabled = %v, want %v", enabled, before)
	}

	// Disable path.
	restore = SetEnabled(false)
	if enabled {
		t.Errorf("after SetEnabled(false): enabled = true, want false")
	}
	restore()
	if enabled != before {
		t.Errorf("after second restore: enabled = %v, want %v", enabled, before)
	}
}
