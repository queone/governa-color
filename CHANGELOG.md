# Changelog

| Version | Summary |
|---------|---------|
| Unreleased | |
| 1.0.0 | API redesign: 10 hues × 11-step intensity ramps (Hue0..Hue10, Hue5 canonical) + 11-step Heat ramp (Heat0..Heat10) + composable Bold/Reverse modifier wrappers. Drops all previous named helpers (Gra/Grn/Yel/Blu/Cya/Mag/Red/Whi/Whi2/GrnR/GrnD/RedR/RedD/BoldW/BoldR/BoldG and Yelf/Redf/Grnf). Drops the basic-ANSI fallback layer (256-color terminals only; non-256 sees uncolored). ShowPalette rewritten to dump 16-color + cube + grayscale + hue ramps + heat ramp; cmd/showpalette runner. See MIGRATION-v1.0.md for the old→new mapping. |
| 0.3.0 | add BoldR + BoldG (bold-red / bold-green companions to BoldW) |
| 0.2.0 | add ClearCode + Mag + Yelf/Redf/Grnf (utils IE14 surface) |
| 0.1.2 | doc: README Why section + CHANGELOG table format (align with governa's CHANGELOG shape) |
| 0.1.1 | fix: remove stray governa-internal AC reference in SetEnabled godoc |
| 0.1.0 | initial extraction from governa internal/color |
