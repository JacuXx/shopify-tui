# @jacuxx/shopify-tui

Interactive TUI (Terminal UI) for managing multiple Shopify development servers — Vim-style navigation, background servers, real-time logs, and more.

## Installation

**macOS / Linux:**
```sh
curl -fsSL https://raw.githubusercontent.com/JacuXx/shopify-tui/main/install.sh | sh
```

**Windows (PowerShell):**
```powershell
iwr https://raw.githubusercontent.com/JacuXx/shopify-tui/main/install.ps1 | iex
```

## Usage

```bash
sho
```

> **Requirement:** [Shopify CLI](https://shopify.dev/docs/api/shopify-cli) must be installed: `npm install -g @shopify/cli`

## Features

- **Multi-store management** — Save and switch between multiple Shopify stores
- **Background servers** — Run multiple dev servers simultaneously
- **Real-time logs** — Interactive log viewer with scroll and selection
- **Vim-style navigation** — `j`/`k` to move, `l`/`Enter` to select
- **Quick actions popup** — `space`/`m` for Pull, Push, VS Code, Terminal
- **Update notifications** — Notified when a new version is available
- **Nerd Font icons** — With automatic ASCII fallback

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `l` / `Enter` | Select |
| `space` / `m` | Open actions popup |
| `Ctrl+Q` | Quit |

Full documentation: [github.com/JacuXx/shopify-tui](https://github.com/JacuXx/shopify-tui)

## License

MIT © [JacuXx](https://github.com/JacuXx)
