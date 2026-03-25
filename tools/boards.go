package tools

import (
	"context"
	"encoding/json"

	"github.com/Andressc19/mcp-trello/trello"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetBoardInput struct {
	BoardID string `json:"boardId" jsonschema:"ID of the board"`
}

func RegisterBoards(server *mcp.Server, client *trello.TrelloClient) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_boards",
		Description: "List all boards the user has access to",
	}, func(_ context.Context, _ *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
		boards, err := client.ListBoards()
		if err != nil {
			return nil, nil, err
		}
		data, err := json.MarshalIndent(boards, "", "  ")
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_board",
		Description: "Get information about a specific board",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in GetBoardInput) (*mcp.CallToolResult, any, error) {
		board, err := client.GetBoard(in.BoardID)
		if err != nil {
			return nil, nil, err
		}
		data, err := json.MarshalIndent(board, "", "  ")
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_board_labels",
		Description: "Get all labels from a board",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in GetBoardInput) (*mcp.CallToolResult, any, error) {
		labels, err := client.GetBoardLabels(in.BoardID)
		if err != nil {
			return nil, nil, err
		}
		data, err := json.MarshalIndent(labels, "", "  ")
		if err != nil {
			return nil, nil, err
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	})
}
