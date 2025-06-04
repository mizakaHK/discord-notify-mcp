package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hk/discord-notify-mcp/internal/discord"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type DiscordNotifyServer struct {
	client    *discord.Client
	mcpServer *server.MCPServer
}

func NewDiscordNotifyServer(webhookURL string) *DiscordNotifyServer {
	s := &DiscordNotifyServer{
		client: discord.NewClient(webhookURL),
	}
	
	// Create MCP server
	s.mcpServer = server.NewMCPServer(
		"discord-notify-mcp",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)
	
	// Register tools
	s.registerTools()
	
	return s
}

func (s *DiscordNotifyServer) registerTools() {
	// Simple message tool
	sendMessageTool := mcp.NewTool("discord_send_message",
		mcp.WithDescription("Send a simple text message to Discord"),
		mcp.WithString("content",
			mcp.Required(),
			mcp.Description("The message content to send"),
		),
	)
	
	s.mcpServer.AddTool(sendMessageTool, s.handleSendMessage)
	
	// Embed message tool - simplified version
	sendEmbedTool := mcp.NewTool("discord_send_embed",
		mcp.WithDescription("Send an embedded message to Discord with rich formatting"),
		mcp.WithString("title",
			mcp.Description("The embed title"),
		),
		mcp.WithString("description",
			mcp.Description("The embed description"),
		),
		mcp.WithNumber("color",
			mcp.Description("The embed color (decimal color value)"),
		),
		mcp.WithString("fields_json",
			mcp.Description("JSON array of fields. Each field should have 'name', 'value', and optional 'inline' properties"),
		),
	)
	
	s.mcpServer.AddTool(sendEmbedTool, s.handleSendEmbed)
}

func (s *DiscordNotifyServer) GetMCPServer() *server.MCPServer {
	return s.mcpServer
}

func (s *DiscordNotifyServer) handleSendMessage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	content, ok := args["content"].(string)
	if !ok {
		return mcp.NewToolResultError("content must be a string"), nil
	}

	err := s.client.SendMessage(content)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to send message: %v", err)), nil
	}

	return mcp.NewToolResultText("Message sent successfully"), nil
}

func (s *DiscordNotifyServer) handleSendEmbed(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	embed := discord.Embed{}

	if title, ok := args["title"].(string); ok {
		embed.Title = title
	}

	if description, ok := args["description"].(string); ok {
		embed.Description = description
	}

	if color, ok := args["color"].(float64); ok {
		embed.Color = int(color)
	}

	// Parse fields from JSON string
	if fieldsJSON, ok := args["fields_json"].(string); ok && fieldsJSON != "" {
		var fields []discord.Field
		if err := json.Unmarshal([]byte(fieldsJSON), &fields); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid fields_json format: %v", err)), nil
		}
		embed.Fields = fields
	}

	err := s.client.SendEmbed(embed)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to send embed: %v", err)), nil
	}

	return mcp.NewToolResultText("Embed sent successfully"), nil
}