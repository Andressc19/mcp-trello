package tools

import (
	"context"
	"encoding/json"

	"github.com/Andressc19/mcp-trello-go/trello"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetListsInput struct {
	BoardID string `json:"boardId" jsonschema:"ID of the board"`
}

type CreateListInput struct {
	BoardID string `json:"boardId" jsonschema:"ID of the board"`
	Name    string `json:"name" jsonschema:"Name of the list"`
}

type UpdateListInput struct {
	ListID string `json:"listId" jsonschema:"ID of the list"`
	Name   string `json:"name,omitempty" jsonschema:"New name for the list"`
	Closed bool   `json:"closed,omitempty" jsonschema:"Whether the list is closed"`
}

type ArchiveListInput struct {
	ListID string `json:"listId" jsonschema:"ID of the list to archive"`
}

func RegisterLists(server *mcp.Server, client *trello.TrelloClient) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_lists",
		Description: "Get all lists from a board",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in GetListsInput) (*mcp.CallToolResult, any, error) {
		lists, err := client.GetLists(in.BoardID)
		if err != nil {
			return nil, nil, err
		}
		data, err := json.MarshalIndent(lists, "", "  ")
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_list",
		Description: "Create a new list on a board",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in CreateListInput) (*mcp.CallToolResult, any, error) {
		list, err := client.CreateList(in.BoardID, in.Name)
		if err != nil {
			return nil, nil, err
		}
		data, err := json.MarshalIndent(list, "", "  ")
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "update_list",
		Description: "Update a list (rename)",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in UpdateListInput) (*mcp.CallToolResult, any, error) {
		params := map[string]string{}
		if in.Name != "" {
			params["name"] = in.Name
		}
		if in.Closed {
			params["closed"] = "true"
		}
		list, err := client.UpdateList(in.ListID, params)
		if err != nil {
			return nil, nil, err
		}
		data, err := json.MarshalIndent(list, "", "  ")
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "archive_list",
		Description: "Archive a list",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in ArchiveListInput) (*mcp.CallToolResult, any, error) {
		list, err := client.ArchiveList(in.ListID)
		if err != nil {
			return nil, nil, err
		}
		data, err := json.MarshalIndent(list, "", "  ")
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})
}
