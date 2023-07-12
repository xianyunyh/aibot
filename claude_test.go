package aibot

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func Test_ClaudeSendMessage(t *testing.T) {
	cookie := ""
	conversationId := uuid.New().String()
	organizationId := uuid.New().String()
	bot := NewClaudeBot(cookie, organizationId)
	reader, err := bot.SendMessage(conversationId, "今天天气怎么样")
	if err != nil {
		t.Fatal(err)
	}
	text := reader.Recive()
	fmt.Println(text)
}
