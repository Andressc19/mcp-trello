package tools

import (
	"context"
	"encoding/json"

	"github.com/Andressc19/mcp-trello-go/trello"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetCardsByListInput struct {
	ListID string `json:"listId" jsonschema:"ID of the list"`
}

type GetCardInput struct {
	CardID string `json:"cardId" jsonschema:"ID of the card"`
}

type AddCardInput struct {
	ListID      string   `json:"listId" jsonschema:"ID of the list"`
	Name        string   `json:"name" jsonschema:"Name of the card"`
	Description string   `json:"description,omitempty" jsonschema:"Card description"`
	DueDate     string   `json:"dueDate,omitempty" jsonschema:"Due date (ISO 8601)"`
	Labels      []string `json:"labels,omitempty" jsonschema:"Array of label IDs"`
}

type UpdateCardInput struct {
	CardID      string `json:"cardId" jsonschema:"ID of the card"`
	Name        string `json:"name,omitempty" jsonschema:"New name"`
	Description string `json:"description,omitempty" jsonschema:"New description"`
	DueDate     string `json:"dueDate,omitempty" jsonschema:"New due date"`
	DueComplete bool   `json:"dueComplete,omitempty" jsonschema:"Mark due as complete"`
	ListID      string `json:"listId,omitempty" jsonschema:"New list ID"`
}

type MoveCardInput struct {
	CardID string `json:"cardId" jsonschema:"ID of the card"`
	ListID string `json:"listId" jsonschema:"ID of the target list"`
}

type ArchiveCardInput struct {
	CardID string `json:"cardId" jsonschema:"ID of the card to archive"`
}

type DeleteCardInput struct {
	CardID string `json:"cardId" jsonschema:"ID of the card to delete"`
}

func RegisterCards(server *mcp.Server, client *trello.TrelloClient) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_cards_by_list_id",
		Description: "Get all cards from a specific list",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in GetCardsByListInput) (*mcp.CallToolResult, any, error) {
		cards, err := client.GetCardsByList(in.ListID)
		if err != nil {
			return nil, nil, err
		}
		data, err := json.MarshalIndent(cards, "", "  ")
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_card",
		Description: "Get detailed information about a specific card",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in GetCardInput) (*mcp.CallToolResult, any, error) {
		card, err := client.GetCard(in.CardID)
		if err != nil {
			return nil, nil, err
		}
		data, err := json.MarshalIndent(card, "", "  ")
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "add_card_to_list",
		Description: "Add a new card to a list",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in AddCardInput) (*mcp.CallToolResult, any, error) {
		card, err := client.AddCard(trello.CardInput{
			ListID:      in.ListID,
			Name:        in.Name,
			Description: in.Description,
			DueDate:     in.DueDate,
			Labels:      in.Labels,
		})
		if err != nil {
			return nil, nil, err
		}
		data, err := json.MarshalIndent(card, "", "  ")
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_card",
		Description: "Update card details",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in UpdateCardInput) (*mcp.CallToolResult, any, error) {
		card, err := client.UpdateCard(trello.CardUpdate{
			CardID:      in.CardID,
			Name:        in.Name,
			Description: in.Description,
			DueDate:     in.DueDate,
			DueComplete: in.DueComplete,
			ListID:      in.ListID,
		})
		if err != nil {
			return nil, nil, err
		}
		data, err := json.MarshalIndent(card, "", "  ")
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "move_card",
		Description: "Move a card to another list",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in MoveCardInput) (*mcp.CallToolResult, any, error) {
		card, err := client.MoveCard(in.CardID, in.ListID)
		if err != nil {
			return nil, nil, err
		}
		data, err := json.MarshalIndent(card, "", "  ")
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "archive_card",
		Description: "Archive a card",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in ArchiveCardInput) (*mcp.CallToolResult, any, error) {
		card, err := client.ArchiveCard(in.CardID)
		if err != nil {
			return nil, nil, err
		}
		data, err := json.MarshalIndent(card, "", "  ")
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_card",
		Description: "Permanently delete a card",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in DeleteCardInput) (*mcp.CallToolResult, any, error) {
		err := client.DeleteCard(in.CardID)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: "Card deleted"}},
		}, nil, nil
	})
}
