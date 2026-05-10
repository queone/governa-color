# Changelog

| Version | Summary |
|---------|---------|
| 1.1.0 | feat: add `ShowGrid(token string, reverse bool, fgIndex int)` — prints a bordered side-by-side grid of every ramp (hue ramps then Heat) at every step (0..10). Column headers use the function-name pattern with `X` as the step placeholder (RedX, BluX, HeatX, …); row labels are the X value, with row 5 flagged `*` as canonical. Each cell renders `token` (defaulting to `"TOKEN"` when empty) padded so adjacent shades separate cleanly along borders. With `reverse=true`, cells use background-color SGR (48;5;N) instead of foreground (38;5;N), turning each cell into a solid colored rectangle — useful for evaluating shades as label backgrounds; `fgIndex >= 0` then sets the text color to that 256-color SGR index (negative = terminal default). Reusable by any consumer for picking label colors or comparing close shades. New `cmd/showgrid` runner accepts `-r` (reverse, default text), `-r=N` (reverse with text color N), an optional token arg, and `-h` / `--help` / `-?` for usage. |
| 1.0.1 | bug fix: wrap / Bold / Reverse now also gate on color256, so 16-color TTYs receive plain text (matching the v1.0 doc). v1.0 was emitting 256-color SGR on any TTY regardless of 256-support. |
| 1.0.0 | API redesign: 10 hues × 11-step intensity ramps (Hue0..Hue10, Hue5 canonical) + 11-step Heat ramp (Heat0..Heat10) + composable Bold/Reverse modifier wrappers. Drops all previous named helpers (Gra/Grn/Yel/Blu/Cya/Mag/Red/Whi/Whi2/GrnR/GrnD/RedR/RedD/BoldW/BoldR/BoldG and Yelf/Redf/Grnf). Drops the basic-ANSI fallback layer (256-color terminals only; non-256 sees uncolored). ShowPalette rewritten to dump 16-color + cube + grayscale + hue ramps + heat ramp; cmd/showpalette runner. See MIGRATION-v1.0.md for the old→new mapping. |
| 0.3.0 | add BoldR + BoldG (bold-red / bold-green companions to BoldW) |
| 0.2.0 | add ClearCode + Mag + Yelf/Redf/Grnf (utils IE14 surface) |
| 0.1.2 | doc: README Why section + CHANGELOG table format (align with governa's CHANGELOG shape) |
| 0.1.1 | fix: remove stray governa-internal AC reference in SetEnabled godoc |
| 0.1.0 | initial extraction from governa internal/color |
