package aibot

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

const (
	stopSequence = "stop_sequence"
)

type ClaudeBot struct {
	Cookie           string
	OrganizationUuid string
	client           *http.Client
}

func NewClaudeBot(cookie, organization string) *ClaudeBot {
	return &ClaudeBot{
		Cookie:           cookie,
		OrganizationUuid: organization,
		client:           &http.Client{},
	}
}

type ClaudeBotRequest struct {
	ConversationUuid string `json:"conversation_uuid"`
	OrganizationUuid string `json:"organization_uuid"`
	Text             string `json:"text"`
	Completion       struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
	} `json:"completion"`
}

type ClaudeComplationResponse struct {
	Completion string `json:"completion,omitempty"`
	StopReason string `json:"stop_reason,omitempty"`
	Stop       string `json:"stop,omitempty"`
	Model      string `json:"model,omitempty"`
	Truncated  bool   `json:"truncated,omitempty"`
}
type ClaudeReader struct {
	*StreamDecoder
}

func (r *ClaudeReader) Recive() string {
	resp := &ClaudeComplationResponse{}
	for r.Next() {
		switch r.Field() {
		case "data":
			json.Unmarshal([]byte(r.Value()), resp)
			if resp.StopReason == stopSequence {
				break
			}
		}
	}
	return resp.Completion
}

func NewClaudeDecoder(stream *StreamDecoder) *ClaudeReader {
	return &ClaudeReader{
		StreamDecoder: stream,
	}
}

func (c *ClaudeBot) SendMessage(conversationId, text string) (*ClaudeReader, error) {
	if len(conversationId) == 0 {
		conversationId = uuid.New().String()
	}
	body := &ClaudeBotRequest{
		ConversationUuid: conversationId,
		OrganizationUuid: c.OrganizationUuid,
		Text:             text,
		Completion: struct {
			Model  string `json:"model"`
			Prompt string `json:"prompt"`
		}{
			Model:  "claude-2",
			Prompt: text,
		},
	}
	data, _ := json.Marshal(body)
	log.Println(string(data))
	uri := "https://claude.ai/api/append_message"
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cookie", c.Cookie)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("acccpt", "text/event-stream")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return &ClaudeReader{NewDecoder(resp.Body)}, err
}
