package tools

import (
	"context"
	"encoding/json"

	"github.com/Andressc19/mcp-trello-go/trello"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type CreateChecklistInput struct {
	CardID string `json:"cardId" jsonschema:"ID of the card"`
	Name   string `json:"name,omitempty" jsonschema:"Name of the checklist (default: Checklist)"`
}

type AddChecklistItemInput struct {
	ChecklistID string `json:"checklistId" jsonschema:"ID of the checklist"`
	Text        string `json:"text" jsonschema:"Text of the item"`
}

type CompleteCheckitemInput struct {
	CardID      string `json:"cardId" jsonschema:"ID of the card"`
	ChecklistID string `json:"checklistId" jsonschema:"ID of the checklist"`
	ItemID      string `json:"itemId" jsonschema:"ID of the item to complete"`
}

type UncompleteCheckitemInput struct {
	CardID      string `json:"cardId" jsonschema:"ID of the card"`
	ChecklistID string `json:"checklistId" jsonschema:"ID of the checklist"`
	ItemID      string `json:"itemId" jsonschema:"ID of the item to uncomplete"`
}

func RegisterChecklists(server *mcp.Server, client *trello.TrelloClient) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_checklist",
		Description: "Create a new checklist on a card",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in CreateChecklistInput) (*mcp.CallToolResult, any, error) {
		checklist, err := client.CreateChecklist(in.CardID, in.Name)
		if err != nil {
			return nil, nil, err
		}
		data, err := json.MarshalIndent(checklist, "", "  ")
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "add_checklist_item",
		Description: "Add an item to a checklist",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in AddChecklistItemInput) (*mcp.CallToolResult, any, error) {
		item, err := client.AddChecklistItem(in.ChecklistID, in.Text)
		if err != nil {
			return nil, nil, err
		}
		data, err := json.MarshalIndent(item, "", "  ")
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "complete_checkitem",
		Description: "Mark a checklist item as complete",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in CompleteCheckitemInput) (*mcp.CallToolResult, any, error) {
		err := client.CompleteCheckItem(in.CardID, in.ChecklistID, in.ItemID)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: "Checkitem marked as complete"}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "uncomplete_checkitem",
		Description: "Mark a checklist item as incomplete",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in UncompleteCheckitemInput) (*mcp.CallToolResult, any, error) {
		err := client.UncompleteCheckItem(in.CardID, in.ChecklistID, in.ItemID)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: "Checkitem marked as incomplete"}},
		}, nil, nil
	})
}
