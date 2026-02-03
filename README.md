# Lane

**The fastest way to generate Stripe invoices from your terminal.**

```
lane 500 --client "Apple" --desc "Consulting" --send
```

---

## Features

- **Instant invoices** — Create and send Stripe invoices in seconds
- **Auto-copy** — Payment links copied to clipboard automatically
- **Email delivery** — Send invoices directly to clients
- **Multi-currency** — USD, EUR, GBP, and 100+ currencies
- **Beautiful output** — Clean, styled terminal interface
- **Cross-platform** — macOS, Linux, and Windows

---

## Installation

### From source

```bash
# Clone and build
git clone https://github.com/forrestcai35/lane.git
cd lane/cli/src
make build

# Or install to GOPATH
make install
```

### Pre-built binaries

Download from [Releases](https://github.com/forrestcai35/lane/releases) for your platform.

---

## Quick Start

```bash
# 1. Authenticate with your Lane account
lane login

# 2. Create your first invoice
lane 100 --client "Acme Corp" --desc "Consulting services"
```

---

## Usage

```
lane <amount> [flags]
```

### Examples

```bash
# Basic invoice
lane 100 --client "Acme Corp" --desc "Consulting"

# With email delivery
lane 500 --client "Apple" --desc "Web Design" --email "tim@apple.com" --send

# Different currency
lane 2500 --desc "Logo Design" --currency eur

# Without clipboard copy
lane 750 --client "Startup Inc" --desc "API Development" --no-copy
```

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--client` | `-c` | Client name |
| `--email` | `-e` | Client email address |
| `--desc` | `-d` | Invoice description (required) |
| `--currency` | | Currency code (default: `usd`) |
| `--send` | | Send invoice via email (requires `--email`) |
| `--no-copy` | | Don't copy payment link to clipboard |

### Commands

| Command | Description |
|---------|-------------|
| `lane login` | Authenticate with Lane |
| `lane logout` | Remove stored credentials |

---

## Authentication

Lane uses browser-based authentication:

```bash
lane login
```

This opens your browser to authenticate. Once complete, your CLI is automatically connected. Credentials are stored locally in `~/.lane/`.

---

## Development

```bash
cd cli/src

# Build
make build

# Run with arguments
make run ARGS="500 --client Apple --desc Test"

# Run tests
make test

# Build for all platforms
make release
```

---

## License

MIT
