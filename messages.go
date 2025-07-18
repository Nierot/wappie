package main

import (
	"fmt"

	"go.mau.fi/whatsmeow/types/events"
)

type Message struct {
	PlainMessage string
	DisplayName string
}

type MessageHandler struct {
	Name string

	Logs
}

func (m *MessageHandler) Handle(e *events.Message) {
	plainMsg := e.Message.GetConversation()
	displayName := e.Message.Chat.DisplayName

	m.log("Plain message", plainMsg)
	m.log("Display name")
}

func (m *MessageHandler) log(name string, arg interface{}) {
	fmt.Print("New event:\n")
	for _, v := range args {
		fmt.Print("")
	}
}

func (m *MessageHandler) print() {

}