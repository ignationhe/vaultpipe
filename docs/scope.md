# scope

Resolve a flat `.env` file that contains scope-prefixed keys into a single
environment-specific file.

## Motivation

It is common to keep all environment variants in a single source file:

```
APP_NAME=myapp
DB_URL=localhost:5432/dev
STAGING__DB_URL=staging.db.example.com:5432/app
PROD__DB_URL=prod.db.example.com:5432/app
PROD__SECRET_KEY=supersecret
```

`vaultpipe scope` resolves the correct values for a target environment so
you can write a clean `.env.prod` without manual editing.

## Usage

```
vaultpipe scope --scopes staging,prod --target prod --input .env --output .env.prod
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--input`, `-i` | `.env` | Source env file |
| `--output`, `-o` | stdout | Destination file |
| `--scopes`, `-s` | *(required)* | Comma-separated ordered scope list |
| `--target`, `-t` | *(required)* | Scope to resolve |
| `--sep` | `__` | Separator between and key |
| `--keep-prefix` | `false` | Preserve scope prefix in output keys |

## Resolution rules

1. **Global keys** (no scope prefix) are always included as defaults.
2. Scopes are applied **in the order declared**; later scopes override earlier ones.
3. Only keys matching the `--target` scope override globals in the output.
4. Keys belonging to *other* scopes are silently dropped.

## Example

Given `.env`:
```
LOG_LEVEL=info
PROD__LOG_LEVEL=warn
PROD__API_KEY=abc123
STAGING__API_KEY=stagingkey
```

Running:
```
vaultpipe scope -s staging,prod -t prod
```

Outputs:
```
LOG_LEVEL=warn
API_KEY=abc123
```
