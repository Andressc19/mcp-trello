package main

import (
	"context"
	"log"
	"os"

	"github.com/Andressc19/mcp-trello-go/tools"
	"github.com/Andressc19/mcp-trello-go/trello"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	log.SetOutput(os.Stderr)

	creds, err := LoadCredentials()
	if err != nil {
		log.Fatal(err)
	}

	client := trello.NewTrelloClient(creds.APIKey, creds.Token)

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "trello",
		Version: "1.0.0",
	}, nil)

	tools.RegisterBoards(server, client)
	tools.RegisterLists(server, client)
	tools.RegisterCards(server, client)
	tools.RegisterLabels(server, client)
	tools.RegisterChecklists(server, client)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
