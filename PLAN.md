# FPL CLI Implementation Plan

## CLI Structure
- Adopt `spf13/cobra` for command parsing to leverage automatic help text, flag management, and future shell completions.
- Define the root command `fpl` with global flags such as `--json`, `--cache`, plus Cobra’s built-in `--help` and optional `--version`.
- Add a `player` subcommand responsible for `--id`, `--name`, and gameweek-related flags; set informative `Use`, `Short`, `Long`, and `Example` strings so `fpl --help` and `fpl player --help` stay descriptive.
- Keep HTTP interactions in `internal/fpl` with typed clients for `bootstrap-static` and `element-summary/{id}`, enabling reuse as more endpoints are added.
- Include an optional in-memory cache (toggle via flag/env) to avoid repeated downloads of `bootstrap-static` within a single invocation.

## Gameweek Flag UX
- Implement a custom flag type that can parse `--gw 5`, comma-separated lists, or inclusive ranges like `--gw 1-3`.
- Accept repeated flags (e.g., `--gw 1 --gw 4-6`) and normalize them into merged ranges early to simplify downstream filtering.
- Validate input: ensure positive integers, `start <= end`, and the range falls within the active season; error messages should highlight the exact issue and suggest valid syntax.
- Filter the `element-summary` history entries by the parsed ranges, returning either aggregated totals or per-GW rows depending on output mode.

## Output & UX
- Default to a human-readable table via `text/tabwriter`, showing goals, assists, minutes, total points, form, ICT index, and cost.
- Support `--json` to emit raw filtered data for scripting; ensure the JSON schema includes both player metadata and GW-specific stats.
- When fuzzy name matching yields low confidence, print the top suggestions with confidence scores so users can refine their query.

## Testing
- Unit test the API client using `httptest.Server` with recorded fixtures for `bootstrap-static` and `element-summary`.
- Cover the GW flag parser with edge cases (single GW, overlapping ranges, invalid input) and ensure normalization works.
- Add command-layer tests using Cobra’s testing helpers to confirm flag wiring and error messages.

## Documentation
- Expand `README.md` with installation steps, usage examples, GW filter explanations, and JSON output samples that mirror the CLI help text.
- Optionally add `docs/cli.md` and use Cobra’s doc generation tooling later as the command surface grows.
- Highlight environment variables or config options (e.g., cache TTL, API base URL) for power users.

## CI & Tooling
- Initialize Go modules with the latest Go version and add `Makefile` targets (`make build`, `make test`, `make lint`).
- Configure GitHub Actions (or preferred CI) using `actions/setup-go@v5` with `go-version: stable` to always pull the newest release.
- Run `go test ./...` and linting in CI; consider adding coverage reporting once the codebase grows.

## Next Steps
1. Scaffold the Cobra project structure (`cmd/fplcli/main.go`, `cmd/player.go`, `internal/fpl` package).
2. Implement the FPL API client and fuzzy player lookup.
3. Build the GW parsing utility, integrate it into the player command, and format the output (table + JSON).
4. Flesh out documentation and CI workflow, ensuring local `make test` mirrors CI behavior.
