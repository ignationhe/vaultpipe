# `vaultpipe merge`

Merge two `.env` files with configurable conflict resolution.

## Usage

```
vaultpipe merge <base-file> <incoming-file> [flags]
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--strategy`, `-s` | `keep` | Conflict resolution strategy: `keep`, `overwrite`, or `both` |
| `--suffix` | `_NEW` | Suffix appended to the incoming key when `--strategy=both` |
| `--output`, `-o` | *(base file)* | Destination file; defaults to overwriting the base file |

## Strategies

### `keep` (default)

Existing keys in the base file are preserved. Incoming keys that do not exist
in the base are added.

```
base:     API_KEY=old
incoming: API_KEY=new  →  API_KEY=old   (unchanged)
```

### `overwrite`

Incoming values replace existing ones on conflict.

```
base:     API_KEY=old
incoming: API_KEY=new  →  API_KEY=new
```

### `both`

Both values are kept. The incoming key is renamed with the configured suffix.

```
base:     API_KEY=old
incoming: API_KEY=new  →  API_KEY=old  +  API_KEY_NEW=new
```

## Examples

```bash
# Merge vault-fetched secrets into .env without overwriting local overrides
vaultpipe merge .env .env.vault --strategy keep

# Overwrite local .env with fresh vault secrets, save to .env.merged
vaultpipe merge .env .env.vault --strategy overwrite --output .env.merged

# Inspect both versions side-by-side
vaultpipe merge .env .env.vault --strategy both --suffix _VAULT
```
