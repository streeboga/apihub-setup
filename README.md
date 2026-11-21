# apihub-setup / Автонастройщик APIHub

> **EN:** A cross-platform CLI tool that automatically configures AI coding tools (OpenCode, Claude Code, Codex, Cursor, Continue, Aider, Cline) to work through [APIHub](https://hub.gnzs.pro) — a marketplace of AI models and APIs.
>
> **RU:** Кроссплатформенный CLI-инструмент, который автоматически настраивает AI-кодинговые инструменты (OpenCode, Claude Code, Codex, Cursor, Continue, Aider, Cline) для работы через [APIHub](https://hub.gnzs.pro) — биржу моделей ИИ и API.

---

## EN: Quick Start

### 1. Download

Download the latest release for your OS from [Releases](https://github.com/streeboga/apihub-setup/releases):

| OS | File |
|---|---|
| macOS (Apple Silicon) | `apihub-setup-darwin-arm64` |
| macOS (Intel) | `apihub-setup-darwin-amd64` |
| Linux (x64) | `apihub-setup-linux-amd64` |
| Windows (x64) | `apihub-setup-windows-amd64.exe` |

Or build from source:

```bash
git clone https://github.com/streeboga/apihub-setup.git
cd apihub-setup
go build -o apihub-setup ./cmd/apihub-setup
```

### 2. Get your APIHub token

1. Go to [hub.gnzs.pro](https://hub.gnzs.pro) and log in (or register)
2. Open **Settings → API Keys** (or `/app/tokens`)
3. Click **Create token** — copy the token (starts with a number, e.g. `1234|abc...`)
4. The token has these abilities by default: `gateway:call`, `wallet:read`, `usage:read`

> ⚠️ The token is shown only once. Save it securely.

### 3. Run

```bash
# Configure all detected tools
apihub-setup --token=YOUR_TOKEN

# Configure a specific tool
apihub-setup --token=YOUR_TOKEN --tool=opencode

# Preview changes without writing
apihub-setup --token=YOUR_TOKEN --dry-run

# List detected tools
apihub-setup --list-tools

# Restore from backup
apihub-setup --restore --tool=opencode
```

Or use environment variable:

```bash
export APIHUB_TOKEN=YOUR_TOKEN
apihub-setup
```

### Supported tools

| Tool | Config file | Status |
|---|---|---|
| OpenCode | `~/.config/opencode/opencode.json` | ✅ |
| Claude Code | `~/.claude/settings.json` | ✅ |
| Codex | `~/.codex/config.toml` | ✅ |
| Cursor | `~/Library/Application Support/Cursor/User/settings.json` | ✅ |
| Continue | `~/.continue/config.json` | ✅ |
| Aider | `~/.aider.conf.yml` | ✅ |
| Cline | VS Code `settings.json` | ✅ |

### What it does

1. Fetches the full model catalog from APIHub (300+ chat models)
2. Detects which AI coding tools are installed on your machine
3. For each tool: creates a backup, patches the config with the APIHub provider and all models, writes atomically
4. Shows a summary of what was changed

No proxy server is started — tools connect directly to APIHub's gateway.

---

## RU: Быстрый старт

### 1. Скачайте

Скачайте последний релиз для вашей ОС из [Releases](https://github.com/streeboga/apihub-setup/releases):

| ОС | Файл |
|---|---|
| macOS (Apple Silicon) | `apihub-setup-darwin-arm64` |
| macOS (Intel) | `apihub-setup-darwin-amd64` |
| Linux (x64) | `apihub-setup-linux-amd64` |
| Windows (x64) | `apihub-setup-windows-amd64.exe` |

Или соберите из исходников:

```bash
git clone https://github.com/streeboga/apihub-setup.git
cd apihub-setup
go build -o apihub-setup ./cmd/apihub-setup
```

### 2. Получите токен APIHub

1. Зайдите на [hub.gnzs.pro](https://hub.gnzs.pro) и войдите (или зарегистрируйтесь)
2. Откройте **Настройки → API-ключи** (или `/app/tokens`)
3. Нажмите **Создать ключ** — скопируйте токен (начинается с числа, например `1234|abc...`)
4. По умолчанию ключ имеет права: `gateway:call`, `wallet:read`, `usage:read`

> ⚠️ Токен показывается только один раз. Сохраните его в безопасном месте.

### 3. Запустите

```bash
# Настроить все обнаруженные инструменты
apihub-setup --token=ВАШ_ТОКЕН

# Настроить конкретный инструмент
apihub-setup --token=ВАШ_ТОКЕН --tool=opencode

# Показать что будет изменено без записи
apihub-setup --token=ВАШ_ТОКЕН --dry-run

# Список обнаруженных инструментов
apihub-setup --list-tools

# Восстановить из бэкапа
apihub-setup --restore --tool=opencode
```

Или через переменную окружения:

```bash
export APIHUB_TOKEN=ВАШ_ТОКЕН
apihub-setup
```

### Поддерживаемые инструменты

| Инструмент | Файл конфига | Статус |
|---|---|---|
| OpenCode | `~/.config/opencode/opencode.json` | ✅ |
| Claude Code | `~/.claude/settings.json` | ✅ |
| Codex | `~/.codex/config.toml` | ✅ |
| Cursor | `~/Library/Application Support/Cursor/User/settings.json` | ✅ |
| Continue | `~/.continue/config.json` | ✅ |
| Aider | `~/.aider.conf.yml` | ✅ |
| Cline | VS Code `settings.json` | ✅ |

### Что делает инструмент

1. Загружает полный каталог моделей из APIHub (300+ chat-моделей)
2. Обнаруживает, какие AI-кодинговые инструменты установлены на вашей машине
3. Для каждого инструмента: создаёт бэкап, патчит конфиг провайдером APIHub и всеми моделями, записывает атомарно
4. Показывает сводку изменений

Локальный прокси-сервер не запускается — инструменты подключаются напрямую к шлюзу APIHub.

---

## Безопасность / Security

- **EN:** A timestamped backup is created before any config change. Restore with `--restore`.
- **RU:** Перед изменением конфига создаётся бэкап с timestamp. Восстановление: `--restore`.

## Architecture / Архитектура

```
cmd/apihub-setup/main.go          — CLI entry point / точка входа
internal/apihub/client.go         — APIHub API client / клиент API
internal/backup/backup.go          — Backup/restore / бэкап и восстановление
internal/detector/detector.go     — Tool detector / детектор инструментов
internal/config/                  — Config patchers / патчеры конфигов
  opencode.go, claudecode.go, codex.go, cursor.go,
  continue.go, aider.go, cline.go
```

## License

MIT
