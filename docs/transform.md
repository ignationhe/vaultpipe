# `vaultpipe transform`

Apply bulk value transformations to an env file: uppercase, lowercase, or trim whitespace.

## Usage

```
vaultpipe transform [flags]
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--input`, `-i` | `.env` | Source env file to read |
| `--output`, `-o` | _(same as input)_ | Destination file (overwrites input if omitted) |
| `--op` | `trim` | Operation: `uppercase`, `lowercase`, `trim` |
| `--keys` | _(all)_ | Comma-separated list of keys to transform; omit to apply to all |
| `--skip-errors` | `false` | Continue processing if a transform fails |

## Examples

### Trim all values

```bash
vaultpipe transform --input .env --op trim
```

### Uppercase specific keys

```bash
vaultpipe transform --input .env --op uppercase --keys DB_PASSWORD,API_SECRET
```

### Lowercase and write to a new file

```bash
vaultpipe transform --input .env --op lowercase --output .env.lower
```

## Behaviour

- **Exact key rules** take priority over the wildcard (`*`) rule.
- The source file is **not modified** until the write step succeeds.
- Use `--skip-errors` when running bulk transforms where individual failures are acceptable.

## Programmatic API

```go
import "github.com/your-org/vaultpipe/internal/envfile"

opts := envfile.DefaultTransformOptions()
opts.Rules["SECRET"] = envfile.UppercaseValues()
opts.Rules["*"]      = envfile.TrimSpaceValues()

result, err := envfile.Transform(src, opts)
```
