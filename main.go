package main

import (
	"context"
	"log"
	"os"

	"github.com/Andressc19/mcp-trello/tools"
	"github.com/Andressc19/mcp-trello/trello"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Check for setup command
	if len(os.Args) > 1 && os.Args[1] == "setup" {
		if err := RunSetup(); err != nil {
			log.Fatal(err)
		}
		return
	}

	log.SetOutput(os.Stderr)

	creds, err := LoadCredentials()
	if err != nil {
		log.Fatal(err, "\n\nRun 'mcp-trello setup' to configure credentials")
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
