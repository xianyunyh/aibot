package aibot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

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
	req.Header.Set("Acccpt", "text/event-stream")
	c.client.Transport = &http.Transport{
		Proxy: func(r *http.Request) (*url.URL, error) {
			return url.Parse("http://127.0.0.1:7890")
		},
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.Header.Get("content-type") != "text/event-stream" {
		data, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		log.Println(string(data))
		return nil, fmt.Errorf("content-type is not text/event-stream")
	}
	return &ClaudeReader{NewDecoder(resp.Body)}, err
}
