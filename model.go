package aibot

type CreateConversationRequest struct {
	ConversationId string `json:"conversation_id"`
	MessageId      string `json:"message_id"`
	Role           string `json:"role"`
	Prompt         string `json:"prompt"`
	SourceType     string `json:"source_type"`
	IsRegenerate   bool   `json:"is_regenerate"`
	IsSo           bool   `json:"is_so"`
}

type CreateConversationResponse struct {
}

type CreateMessageRequest struct {
}
