# vaultpipe watch

The `watch` command polls a `.env` file and prints its contents the file changes.

## Usage

```bash
vaultpipe watch <file> [flags]
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--interval` | `-i` | `2s` | Polling interval (e.g. `500ms`, `5s`) |
| `--quiet` | `-q` | `false` | Suppress the startup banner |

## Examples

### Basic watch

```bash
vaultpipe watch .env
# Watching .env every 2s (Ctrl+C to stop)
# [changed] .env — 3 keys
#   DB_HOST=localhost
#   DB_PORT=5432
#   APP_ENV=development
```

### Fast polling

```bash
vaultpipe watch .env --interval 500ms
```

### Quiet mode (useful in scripts)

```bash
vaultpipe watch .env --quiet 2>/dev/null
```

## How it works

1. On startup the file is hashed with SHA-256.
2. Every `--interval` the file is re-hashed.
3. If the hash differs the file is parsed and `OnChange` fires.
4. Errors (missing file, bad syntax) are written to stderr via `OnError`.
5. The loop exits cleanly on `SIGINT` / `SIGTERM`.

## Notes

- The watcher uses **polling**, not `inotify`, so it works on all platforms including network-mounted volumes.
- Combine with `vaultpipe sync` in a shell loop to keep secrets fresh without a daemon process.
