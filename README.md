# envchain-export

> Utility to safely export and import scoped environment variable sets with encryption support

---

## Installation

```bash
go install github.com/yourusername/envchain-export@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envchain-export.git
cd envchain-export && go build -o envchain-export .
```

---

## Usage

**Export a scoped set of environment variables:**

```bash
envchain-export export --scope production --output env.enc
```

**Import and decrypt into your current shell:**

```bash
eval $(envchain-export import --scope production --input env.enc)
```

**List available scopes:**

```bash
envchain-export list
```

Encryption is handled automatically using AES-256-GCM. You will be prompted for a passphrase on export and import unless `--passphrase` is provided via flag or the `ENVCHAIN_PASSPHRASE` environment variable is set.

---

## Options

| Flag | Description |
|------|-------------|
| `--scope` | Named scope for the variable set |
| `--input` | Path to encrypted input file |
| `--output` | Path for encrypted output file |
| `--passphrase` | Encryption passphrase (use env var in CI) |

---

## License

MIT © [yourusername](https://github.com/yourusername)