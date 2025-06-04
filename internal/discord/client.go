package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	webhookURL string
}

type WebhookMessage struct {
	Content string   `json:"content"`
	Embeds  []Embed  `json:"embeds,omitempty"`
}

type Embed struct {
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	Color       int     `json:"color,omitempty"`
	Fields      []Field `json:"fields,omitempty"`
}

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

func NewClient(webhookURL string) *Client {
	return &Client{
		webhookURL: webhookURL,
	}
}

func (c *Client) SendMessage(content string) error {
	msg := WebhookMessage{
		Content: content,
	}
	return c.send(msg)
}

func (c *Client) SendEmbed(embed Embed) error {
	msg := WebhookMessage{
		Embeds: []Embed{embed},
	}
	return c.send(msg)
}

func (c *Client) send(msg WebhookMessage) error {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	resp, err := http.Post(c.webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}