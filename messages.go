package main

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type MessageType int

const (
	MessagePlain MessageType = iota
	MessageUnknown
)

type Message struct {
	Type MessageType
	PlainMessage string
	DisplayName string
	ID string
	ReceivedAt string
	IsGroup bool
}

func (m *Message) String() string {
	str := "*received message:*\n"

	mt := reflect.ValueOf(m).Elem()
	t := mt.Type()

	for i := 0; i < mt.NumField(); i++ {
		field := t.Field(i)
		value := mt.Field(i)
		val := value.Interface()

		if (val == nil) {
			val = "unset"
		}
		str += fmt.Sprintf("%s: %v", field.Name, val)
		if (i != mt.NumField()-1) {
			str += "\n"
		}
	}

	return str
}

func (m *Message) Log() {
	if (m == nil) {
		panic("message is nil but ::Log is called")
	}

	fmt.Println(m.String())
}

type MessageHandler struct {
	Client *whatsmeow.Client
}

func (m *MessageHandler) Handle(e *events.Message) {
	if (e == nil) {
		return
	}

	// niet op eigen berichten reageren
	if e.Info.IsFromMe {
		return
	}

	msg := m.MakeMessage(e)

	// Selecteer commando

	// stuur alle logs in een DM
	m.Send(e.Info.Chat, msg.String())
	// En in de console
	msg.Log()
}

func (m *MessageHandler) MakeMessage(e *events.Message) Message {
	ctx := context.Background()

	str := ""
	msgType := m.GetMessageType(e)

	switch msgType {
	case MessagePlain:
		str = m.GetMessagePlaintext(e)
	}

	// In groepen en DMs is the JID anders, haal de goede op
	jid, err := m.Client.Store.LIDs.GetLIDForPN(ctx, e.Info.Sender.ToNonAD())

	if err != nil {
		jid = e.Info.Sender.ToNonAD()
	}

	msg := Message{
		Type: msgType,
		PlainMessage: str,
		DisplayName: e.Info.PushName,
		ID: jid.String(),
		ReceivedAt: e.Info.Timestamp.Local().Format(time.DateTime),
		IsGroup: e.Info.IsGroup,
	}

	return msg
}

func (m *MessageHandler) GetMessagePlaintext(e *events.Message) string {
	// als reply op ander bericht: GetConversation werkt niet

	fmt.Println(e.Info.PushName)

	mt := reflect.TypeOf(e.Message)

	fmt.Println(mt)

	str := e.Message.GetConversation()

	if (str == "") {

	}

	return e.Message.GetConversation()
}

func (m *MessageHandler) GetMessageType(e *events.Message) MessageType {
	if e.Message.GetConversation() != "" {
		return MessagePlain
	}

	return MessageUnknown
}

func (m *MessageHandler) Send(jid types.JID, msgString string) {
	ctx := context.Background()

	msg := &waE2E.Message{
		Conversation: proto.String(msgString),
	}

	m.Client.SendMessage(ctx, jid, msg)
}