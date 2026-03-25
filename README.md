# mcp-trello

MCP server for Trello written in Go. Provides 21 tools for interacting with Trello boards, lists, cards, labels, and checklists via the Model Context Protocol.

## Prerequisites

- Go 1.25+
- Trello API key + token ([get yours here](https://trello.com/power-ups/admin))

## Quick Setup (Interactive TUI)

The fastest way to configure mcp-trello is using the interactive TUI installer:

```bash
mcp-trello setup
```

This launches a beautiful terminal interface (powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea)) that will:

1. **Guide you through credential setup** — Enter your Trello API Key and Token with secure input masking
2. **Save credentials** to `~/.config/opencode/mcp-trello.json`
3. **Automatically configure OpenCode** — Adds the MCP server to your `opencode.json` config
4. **Verify configuration** — Check that everything is set up correctly

### TUI Features

- **Keyboard navigation**: Use arrow keys (`↑/↓`) or `j/k` to navigate, `Enter` to select
- **Input masking**: API keys and tokens are automatically masked for security
- **Browser integration**: Press `o` to open Trello's admin page in your browser
- **Back navigation**: Press `b` to go back to the previous step
- **Verify config**: Check current configuration from the main menu
- **Uninstall option**: Remove credentials and OpenCode config completely

## Manual Installation

```bash
go install github.com/Andressc19/mcp-trello@latest
```

The binary will be installed to `$GOPATH/bin/mcp-trello`. Make sure `$GOPATH/bin` is in your `PATH`.

## Configuration

### Option 1: Interactive Setup (Recommended)

```bash
mcp-trello setup
```

### Option 2: Config File

Create the config file at `~/.config/opencode/mcp-trello.json`:

```json
{
  "apiKey": "your_trello_api_key",
  "token": "your_trello_token"
}
```

### Option 3: Environment Variables

```bash
export TRELLO_API_KEY=your_api_key
export TRELLO_TOKEN=your_token
```

**Priority**: environment variables > config file

## OpenCode Integration

Add to your `~/.config/opencode/opencode.json`:

```json
{
  "mcp": {
    "trello": {
      "command": ["mcp-trello"],
      "enabled": true,
      "type": "local"
    }
  }
}
```

Or let the TUI installer do it automatically with `mcp-trello setup`.

## Development

### Dependencies

For TUI development, the following dependencies are required:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — The framework for TUI applications
- [Bubbles](https://github.com/charmbracelet/bubbles) — UI components for Bubble Tea
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Style definitions for TUI

These are included in `go.mod` and will be automatically fetched when building.

### Build

```bash
go build -o mcp-trello .
```

### Run Tests

```bash
go test ./...
```

## Tools

### Boards (3)

| Tool | Description |
|------|-------------|
| `list_boards` | List all boards the user has access to |
| `get_board` | Get information about a specific board |
| `get_board_labels` | Get all labels from a board |

### Lists (4)

| Tool | Description |
|------|-------------|
| `get_lists` | Get all lists from a board |
| `create_list` | Create a new list on a board |
| `update_list` | Update a list (rename) |
| `archive_list` | Archive a list |

### Cards (7)

| Tool | Description |
|------|-------------|
| `get_cards_by_list_id` | Get all cards from a specific list |
| `get_card` | Get detailed information about a specific card |
| `add_card_to_list` | Add a new card to a list |
| `update_card` | Update card details |
| `move_card` | Move a card to another list |
| `archive_card` | Archive a card |
| `delete_card` | Permanently delete a card |

### Labels (3)

| Tool | Description |
|------|-------------|
| `create_label` | Create a new label on a board |
| `add_label_to_card` | Add a label to a card |
| `remove_label_from_card` | Remove a label from a card |

### Checklists (4)

| Tool | Description |
|------|-------------|
| `create_checklist` | Create a new checklist on a card |
| `add_checklist_item` | Add an item to a checklist |
| `complete_checkitem` | Mark a checklist item as complete |
| `uncomplete_checkitem` | Mark a checklist item as incomplete |

**Total: 21 tools**

## Usage Examples

### List all your boards

```json
{
  "tool": "list_boards"
}
```

### Create a new list on a board

```json
{
  "tool": "create_list",
  "arguments": {
    "boardId": "60a7f5e3b9c8d1234567890a",
    "name": "My New List"
  }
}
```

### Add a card to a list

```json
{
  "tool": "add_card_to_list",
  "arguments": {
    "listId": "60a7f5e3b9c8d1234567890b",
    "name": "New Task",
    "description": "This is a new task",
    "dueDate": "2024-12-31T23:59:59Z",
    "labels": ["60a7f5e3b9c8d1234567890c"]
  }
}
```

### Move a card to another list

```json
{
  "tool": "move_card",
  "arguments": {
    "cardId": "60a7f5e3b9c8d1234567890d",
    "listId": "60a7f5e3b9c8d1234567890e"
  }
}
```

### Create a checklist on a card

```json
{
  "tool": "create_checklist",
  "arguments": {
    "cardId": "60a7f5e3b9c8d1234567890d",
    "name": "To Do Items"
  }
}
```

## License

MIT