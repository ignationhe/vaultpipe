# mask

The `mask` command prints an env file to stdout with sensitive values replaced by a placeholder. This is useful for logging, debugging, or sharing configuration without exposing secrets.

## Usage

```bash
vaultpipe mask <file> [flags]
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--patterns` | — | Additional key patterns to mask (comma-separated) |
| `--placeholder` | `****` | Replacement string for masked values |
| `--partial` | `0` | Reveal first N characters before the placeholder |
| `--case-sensitive` | `false` | Use case-sensitive pattern matching |

## Default Patterns

The following substrings trigger masking by default (case-insensitive):

- `SECRET`
- `PASSWORD`
- `TOKEN`
- `KEY`
- `PRIVATE`
- `CREDENTIAL`

## Examples

### Basic masking

```bash
vaultpipe mask .env
```

Output:
```
DB_PASSWORD=****
API_TOKEN=****
APP_NAME=vaultpipe
```

### Partial reveal

```bash
vaultpipe mask .env --partial 4
```

Output:
```
API_TOKEN=tok_****
```

### Custom placeholder and extra patterns

```bash
vaultpipe mask .env --placeholder '[HIDDEN]' --patterns 'CERT,SEED'
```
