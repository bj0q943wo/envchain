# envchain

Lightweight utility to chain and merge `.env` files with override precedence for local dev workflows.

---

## Installation

```bash
go install github.com/yourusername/envchain@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envchain.git && cd envchain && go build -o envchain .
```

---

## Usage

Pass one or more `.env` files in order of increasing precedence. Later files override earlier ones.

```bash
envchain .env .env.local .env.dev -- your-command --flag
```

**Example:**

```bash
# .env         → BASE_URL=http://localhost:3000, DEBUG=false
# .env.local   → DEBUG=true

envchain .env .env.local -- go run main.go
# Result: BASE_URL=http://localhost:3000, DEBUG=true
```

You can also print the merged environment without running a command:

```bash
envchain --print .env .env.local
```

### Options

| Flag | Description |
|------|-------------|
| `--print` | Print merged env vars to stdout instead of executing a command |
| `--export` | Prefix each line with `export` for shell sourcing |
| `--strict` | Exit with an error if any specified file does not exist |

---

## How It Works

`envchain` reads each `.env` file left to right, merging keys into a single environment map. Duplicate keys are overwritten by the rightmost file, giving you predictable, layered configuration — no surprises.

---

## License

MIT © [yourusername](https://github.com/yourusername)