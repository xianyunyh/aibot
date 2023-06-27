package aibot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

const (
	// 360智脑
	RoleId = "00000001"
)

type Chat360 struct {
	Cookie string
	client *http.Client
}

func NewChat360(cookie string) *Chat360 {
	return &Chat360{
		Cookie: cookie,
		client: &http.Client{},
	}
}

func (c *Chat360) GetCookie() string {
	return c.Cookie
}

func (c *Chat360) SetCookie(cookie string) {
	c.Cookie = cookie
}
func (c *Chat360) SetClient(client *http.Client) {
	c.client = client
}

// CreateConversation 创建一个新的会话
func (c *Chat360) CreateConversation() (string, error) {
	return "", nil
}

// SendMessage 发送会话消息
func (c *Chat360) SendMessage(conversationId, message string) (*StreamDecoder, error) {
	body := &CreateConversationRequest{
		Prompt:       message,
		MessageId:    c.generateMessageId(),
		Role:         RoleId,
		SourceType:   "prophet_web",
		IsRegenerate: false,
		IsSo:         false,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	uri := "https://chat.360.cn/backend-api/api/common/chat"
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
	return NewDecoder(resp.Body), err
}

// generateMessageId 生成消息ID
func (c *Chat360) generateMessageId() string {
	return fmt.Sprintf("msg-%s", c.generateUUID())
}

// GenerateConversationId 生成会话ID
func (c *Chat360) GenerateConversationId() string {
	return fmt.Sprintf("con-%s", c.generateUUID())
}

// generateUUID 生成UUID
func (c *Chat360) generateUUID() string {
	return uuid.New().String()
}
