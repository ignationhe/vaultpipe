# vaultpipe

> A CLI tool to sync secrets from HashiCorp Vault into local `.env` files with diff-aware updates.

---

## Installation

```bash
go install github.com/yourusername/vaultpipe@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultpipe/releases).

---

## Usage

Set your Vault address and token, then run `vaultpipe` pointing at a Vault secret path and a target `.env` file:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.xxxxxxxx"

vaultpipe sync \
  --path "secret/data/myapp/production" \
  --output .env
```

vaultpipe will fetch the secrets from Vault, compare them against your existing `.env` file, and apply only the changed or missing keys — leaving unrelated entries untouched.

### Additional Commands

```bash
# Preview changes without writing to disk
vaultpipe sync --path "secret/data/myapp/production" --output .env --dry-run

# Show current diff between Vault secrets and local .env
vaultpipe diff --path "secret/data/myapp/production" --output .env
```

### Example Output

```
~ DB_PASSWORD   [updated]
+ REDIS_URL     [added]
  API_KEY       [unchanged]

2 changes applied to .env
```

---

## Configuration

| Flag | Env Var | Description |
|------|---------|-------------|
| `--path` | `VAULTPIPE_PATH` | Vault secret path |
| `--output` | `VAULTPIPE_OUTPUT` | Target `.env` file path |
| `--dry-run` | — | Preview changes without writing |

---

## License

[MIT](LICENSE) © 2024 yourusername