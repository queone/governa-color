# governa-color Plan

## Product Direction

Describe the product direction.

## Ideas To Explore

Ideas captured for future reference. A bullet list — each line starts with `- IE<N>: ` (sequential N) for stable references. Two kinds: (a) **pre-rubric IE** — `IE<N>: <one-liner>`, awaiting director discussion and the objective-fit rubric (see `AGENTS.md` Approval Boundaries); (b) **AC-pointer** — `IE<N>: <one-liner> → docs/ac<N>-<slug>.md`, pointing at a drafted AC stub not yet through critique. A pre-rubric entry that clears the rubric converts to an AC-pointer at AC-draft time, keeping its `IE<N>` number. Remove entries when the idea is rejected, retired, or (for AC-pointers) the AC has shipped and its file deleted. Not a historical record.

- IE2: pair `showgrid -r` with a `--reverse` long-form alias to satisfy AGENTS.md `Code Style and Conventions` short+long pairing. Deferred from AC3 because that AC's `cmd/showgrid/main.go` touch was additive (added `-b`/`--background`, did not rewrite `-r`); the AGENTS.md migration trigger ("when their code is next touched") fires the next time `-r`'s code is materially edited.
