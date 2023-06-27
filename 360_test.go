package aibot

import (
	"fmt"
	"testing"
)

func Test_SendMessage(t *testing.T) {
	cookie := ""
	bot := NewChat360(cookie)
	conversationId := bot.generateMessageId()
	decoder, err := bot.SendMessage(conversationId, "hello")
	if err != nil {
		t.Fatal(err)
	}
	for decoder.Next() {
		switch decoder.Field() {
		case "data":
			fmt.Printf("data:%s\n", decoder.Value())
		}
	}
}
