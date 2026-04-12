# `vaultpipe interpolate`

Resolve variable references within a `.env` file, expanding `${VAR}` and `$VAR`
syntax using values from the same file or the OS environment.

## Usage

```
vaultpipe interpolate [flags]
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-i, --input` | `.env` | Source `.env` file to read |
| `-o, --output` | *(same as input)* | Destination file to write resolved values |
| `--fallback-env` | `true` | Fall back to OS environment variables when a key is not found in the file |
| `--error-missing` | `false` | Return a non-zero exit code when a referenced variable cannot be resolved |
| `--default` | `""` | Replacement value for unresolved variables (when `--error-missing` is false) |

## Examples

### Basic interpolation

```dotenv
# .env
BASE_URL=https://example.com
API_URL=${BASE_URL}/api
HEALTH_URL=${BASE_URL}/health
```

```bash
vaultpipe interpolate --input .env --output .env.resolved
```

Resulting `.env.resolved`:

```dotenv
BASE_URL=https://example.com
API_URL=https://example.com/api
HEALTH_URL=https://example.com/health
```

### Fail on missing variables

```bash
vaultpipe interpolate --error-missing
```

### Supply a default for missing variables

```bash
vaultpipe interpolate --default UNKNOWN
```

## Notes

- Both `${VAR}` and `$VAR` syntax are supported.
- Self-referential or circular references are not detected and will produce
  unexpected output; avoid them.
- When `--fallback-env` is enabled (default), the tool checks the OS environment
  after the file's own keys, which is useful for CI pipelines where some values
  are injected externally.
