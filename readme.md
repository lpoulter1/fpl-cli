# FPL API CLI

`fpl` is a Go-based CLI that wraps the public Fantasy Premier League APIs so you can inspect player performance without leaving your terminal.

## Installation

```bash
go install github.com/lpoulter1/fpl-cli/cmd/fplcli@latest
```

This pulls from the public GitHub repo and drops the `fpl` binary in your Go bin directory. Ensure `$(go env GOPATH)/bin` is on your `PATH`.

## Usage

The CLI exposes a root command plus a `player` subcommand. Run `fpl --help` or `fpl player --help` at any time for the latest, auto-generated docs.

Common examples:

```bash
# Resolve by ID
fpl player --id 123

# Fuzzy name matching (web name, full name, or known-as)
fpl player --name "Haaland"

# Filter to a single GW or inclusive ranges
fpl player --name "Salah" --gw 1 --gw 3-5
fpl player --name "Watkins" --gw 2|4|8-10

# JSON output for scripting
fpl player --name "Saka" --gw 1-3 --json | jq
```

### Gameweek Filters

- `--gw` accepts single values (`--gw 1`), inclusive ranges (`--gw 1-3`), or delimited lists (`--gw 1|4|6-8`).  
- Repeat the flag to add more ranges; overlapping values are merged automatically.  
- When omitted, all available gameweeks are returned.

### Output

The default view prints:

- Player metadata (name, team, position, cost, form, ICT, selected by).  
- A per-GW table with opponent, minutes, goals, assists, clean sheets, and points.  
- Aggregated totals for the selected gameweeks.  

Add `--json` to emit the same data structure in machine-friendly JSON (handy for piping into `jq` or other tooling).

## Development

```bash
go test ./...

# local build / run
go build ./cmd/fplcli
./fplcli player --name "Haaland" --gw 1-3
```

The project targets the latest stable Go release. Continuous integration (see `.github/workflows/ci.yml`) mirrors the same `go test ./...` invocation to keep local and remote runs aligned.
