# jutge

Browse and search the [Jutge.org](https://jutge.org) competitive programming problem archive from the command line.

`jutge` is a single pure-Go binary. No login or API key required.

## Install

```bash
go install github.com/tamnd/jutge-cli/cmd/jutge@latest
```

Or grab a prebuilt binary from the [releases](https://github.com/tamnd/jutge-cli/releases), or run the container image:

```bash
docker run --rm ghcr.io/tamnd/jutge:latest --help
```

## Usage

```bash
# List all problems (4600+)
jutge list

# List first 20 problems in table format
jutge list -n 20 -o table

# Search problems by title or code
jutge search "path"
jutge search "shortest" -n 10

# Output formats
jutge list -o json
jutge list -o csv -n 50
jutge search "sort" -o jsonl
```

## Commands

| Command | Description |
|---------|-------------|
| `list` | List all problems in the Jutge archive |
| `search <query>` | Search problems by title or code (case-insensitive) |
| `version` | Show version information |

## Global flags

```
-o, --output string    output format: table|json|jsonl|csv|tsv|url|raw (default "auto")
-n, --limit int        limit number of records (0 = all)
    --fields strings   comma-separated columns to include
    --no-header        omit header row
    --template string  Go text/template per record
    --timeout duration per-request timeout (default 1m0s)
    --delay duration   minimum spacing between requests
    --retries int      retry attempts on 429/5xx (default 3)
```

## License

Apache-2.0. See [LICENSE](LICENSE).
