package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hk/discord-notify-mcp/internal/server"
	"github.com/joho/godotenv"
	mcpServer "github.com/mark3labs/mcp-go/server"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if webhookURL == "" {
		log.Fatal("DISCORD_WEBHOOK_URL is not set")
	}

	// Create Discord notify server
	discordServer := server.NewDiscordNotifyServer(webhookURL)

	// Start the server
	if err := mcpServer.ServeStdio(discordServer.GetMCPServer()); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}