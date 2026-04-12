# `vaultpipe rename`

Rename keys in an `.env` file using exact matches or regular-expression patterns.

## Usage

```
vaultpipe rename [flags]
```

## Flags

| Flag | Default | Description |
|---|---|---|
| `--input`, `-i` | `.env` | Source env file |
| `--output`, `-o` | *(same as input)* | Destination file |
| `--rule`, `-r` | *(required)* | Rename rule as `FROM=TO` (repeatable) |
| `--keep-original` | `false` | Preserve the original key alongside the new one |
| `--error-missing` | `false` | Return an error if the source key does not exist |
| `--pattern` | `false` | Treat `FROM` as a Go regular expression |

## Examples

### Exact rename

```bash
vaultpipe rename --input .env --rule DB_HOST=DATABASE_HOST
```

Before:
```
DB_HOST=localhost
```

After:
```
DATABASE_HOST=localhost
```

### Multiple rules

```bash
vaultpipe rename -i .env -r DB_HOST=DATABASE_HOST -r DB_PORT=DATABASE_PORT
```

### Pattern-based rename

Uses Go's `regexp.ReplaceAllString` syntax for the replacement string.

```bash
vaultpipe rename -i .env --pattern -r '^OLD_(.*)=LEGACY_$1'
```

Before:
```
OLD_API_KEY=abc
OLD_SECRET=xyz
```

After:
```
LEGACY_API_KEY=abc
LEGACY_SECRET=xyz
```

### Keep original key

```bash
vaultpipe rename -i .env -r OLD_KEY=NEW_KEY --keep-original
```

Both `OLD_KEY` and `NEW_KEY` will be present in the output.

## Notes

- By default missing source keys are silently skipped; use `--error-missing` to fail instead.
- Pattern rules automatically uppercase the resulting key name.
- When `--output` is omitted the source file is updated in-place.
