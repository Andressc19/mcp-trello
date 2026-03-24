package tools

import (
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/Andressc19/mcp-trello-go/trello"
)

type CreateLabelInput struct {
	BoardID string `json:"boardId" jsonschema:"ID of the board"`
	Name    string `json:"name" jsonschema:"Name of the label"`
	Color   string `json:"color" jsonschema:"Color of the label (yellow, green, blue, red, purple, orange, sky, pink, black, null)"`
}

type AddLabelToCardInput struct {
	CardID  string `json:"cardId" jsonschema:"ID of the card"`
	LabelID string `json:"labelId" jsonschema:"ID of the label to add"`
}

type RemoveLabelFromCardInput struct {
	CardID  string `json:"cardId" jsonschema:"ID of the card"`
	LabelID string `json:"labelId" jsonschema:"ID of the label to remove"`
}

func RegisterLabels(server *mcp.Server, client *trello.TrelloClient) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_label",
		Description: "Create a new label on a board",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in CreateLabelInput) (*mcp.CallToolResult, any, error) {
		label, err := client.CreateLabel(in.BoardID, in.Name, in.Color)
		if err != nil {
			return nil, nil, err
		}
		data, err := json.MarshalIndent(label, "", "  ")
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "add_label_to_card",
		Description: "Add a label to a card",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in AddLabelToCardInput) (*mcp.CallToolResult, any, error) {
		err := client.AddLabelToCard(in.CardID, in.LabelID)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: "Label added to card"}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "remove_label_from_card",
		Description: "Remove a label from a card",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in RemoveLabelFromCardInput) (*mcp.CallToolResult, any, error) {
		err := client.RemoveLabelFromCard(in.CardID, in.LabelID)
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: "Label removed from card"}},
		}, nil, nil
	})
}
