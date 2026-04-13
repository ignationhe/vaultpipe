# `vaultpipe clone`

Clone an existing `.env` file to a new location, with optional key filtering and prefix injection.

## Usage

```bash
vaultpipe clone <src> <dst> [flags]
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--overwrite` | `false` | Replace the destination file if it already exists |
| `--keys` | _(all)_ | Comma-separated list of keys to copy |
| `--prefix` | _(none)_ | Prefix prepended to every key in the destination |

## Examples

### Basic clone

```bash
vaultpipe clone .env .env.backup
```

### Clone with key filter

```bash
vaultpipe clone .env .env.db --keys DB_HOST,DB_PORT,DB_NAME
```

### Clone with prefix

```bash
vaultpipe clone .env .env.prefixed --prefix APP
# DB_HOST → APP_DB_HOST
```

### Overwrite an existing file

```bash
vaultpipe clone .env.production .env --overwrite
```

## Notes

- The destination directory is created automatically if it does not exist.
- When `--prefix` is used, keys are uppercased and an underscore separator is
  added automatically if the supplied prefix does not already end with `_`.
- Cloning without `--overwrite` is safe: the command exits with an error if the
  destination already exists, preventing accidental data loss.
