# mcp-trello-go

MCP server for Trello written in Go. Provides 21 tools for interacting with Trello boards, lists, cards, labels, and checklists via the Model Context Protocol.

## Prerequisites

- Go 1.21+
- Trello API key + token ([get yours here](https://trello.com/power-ups/admin))

## Installation

```bash
go install github.com/Andressc19/mcp-trello-go@latest
```

The binary will be installed to `$GOPATH/bin/mcp-trello`. Make sure `$GOPATH/bin` is in your `PATH`.

## Configuration

Create the config file at `~/.config/opencode/mcp-trello.json`:

```json
{
  "apiKey": "your_trello_api_key",
  "token": "your_trello_token"
}
```

Alternatively, set environment variables:

```bash
export TRELLO_API_KEY=your_api_key
export TRELLO_TOKEN=your_token
```

Priority: environment variables > config file.

## OpenCode Integration

Add to your `~/.config/opencode/opencode.json`:

```json
{
  "mcp": {
    "trello": {
      "command": ["mcp-trello"]
    }
  }
}
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
