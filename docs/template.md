# `vaultpipe template` — .env Template Rendering

The `template` command renders a `.env` template file by substituting
`${KEY}` placeholders with values from one or more source files and/or
the current process environment.

## Usage

```sh
vaultpipe template --src .env.tpl --dst .env [--vars base.env] [--vars overrides.env] [--missing-key error|keep|empty]
```

### Flags

| Flag | Short | Default | Description |
|---|---|---|---|
| `--src` | `-s` | *(required)* | Source template file |
| `--dst` | `-d` | *(required)* | Destination output file |
| `--vars` | `-v` | — | `.env` file(s) providing substitution values (repeatable) |
| `--missing-key` | — | `error` | Behaviour for unresolved placeholders: `error`, `keep`, or `empty` |

## Placeholder Syntax

Placeholders follow the `${KEY}` pattern where `KEY` must start with a
letter or underscore and contain only letters, digits, and underscores.

```
# .env.tpl
DATABASE_URL=postgres://${DB_HOST}:${DB_PORT}/${DB_NAME}
APP_SECRET=${SECRET_KEY}
```

## Variable Resolution Order

1. Values from `--vars` files (later files override earlier ones).
2. Process environment variables (do **not** override var-file values).

## Missing Key Behaviour

| Mode | Result |
|---|---|
| `error` | Command exits non-zero with a descriptive message |
| `keep` | Placeholder is left verbatim in the output |
| `empty` | Placeholder is replaced with an empty string |

## Example

```sh
# Pull secrets from Vault first, then render a template
vaultpipe sync --path secret/data/app --output .env.secrets
vaultpipe template --src .env.tpl --dst .env --vars .env.secrets
```
