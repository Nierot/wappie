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

//go:generate stringer -type=MessageType --linecomment
type MessageType int

const (
	MessagePlain MessageType = iota
	MessageExtended
	MessageUnknown
)

type Message struct {
	DisplayName string
	MessageID string
	SenderID string
	ChatID string
	ReceivedAt string
	IsGroup bool
	Type MessageType
	PlainMessage string
}

func (m *Message) String() string {
	// str := "*received message:*\n"
	str := ""

	mt := reflect.ValueOf(m).Elem()
	t := mt.Type()

	for i := 0; i < mt.NumField(); i++ {
		field := t.Field(i)
		value := mt.Field(i)
		val := value.Interface()

		if (val == nil || val == "") {
			val = "unset"
		}

		if valType := reflect.TypeOf(val); valType.String() == "MessageType" {
			val = valType.Elem().String()
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
	m.MarkRead(&msg, e)
	// En in de console
	msg.Log()
}

func (m *MessageHandler) MarkRead(msg *Message, e *events.Message) {
	m.Client.MarkRead([]string{msg.MessageID}, time.Now(), e.Info.Chat, e.Info.Sender, types.ReceiptTypeRead)
}

func (m *MessageHandler) MakeMessage(e *events.Message) Message {
	ctx := context.Background()

	msgType := m.GetMessageType(e)
	str := m.GetMessagePlaintext(e, msgType)

	// In groepen en DMs is the JID anders, haal de goede op
	jid, err := m.Client.Store.LIDs.GetLIDForPN(ctx, e.Info.Sender.ToNonAD())

	if err != nil {
		jid = e.Info.Sender.ToNonAD()
	}

	msg := Message{
		Type: msgType,
		PlainMessage: str,
		DisplayName: e.Info.PushName,
		MessageID: e.Info.ID,
		SenderID: jid.String(),
		ChatID: e.Info.Chat.String(),
		ReceivedAt: e.Info.Timestamp.Local().Format(time.DateTime),
		IsGroup: e.Info.IsGroup,
	}

	return msg
}

func (m *MessageHandler) GetMessagePlaintext(e *events.Message, msgType MessageType) string {
	// als reply op ander bericht: GetConversation werkt niet

	switch msgType {
	case MessagePlain:
		return e.Message.GetConversation()
	case MessageExtended:
		return e.Message.ExtendedTextMessage.GetText()
	}

	return ""
}

func (m *MessageHandler) GetMessageType(e *events.Message) MessageType {
	if len(e.Message.GetConversation()) != 0 {
		return MessagePlain
	}

	fmt.Println(e.Message)
	fmt.Println(e.Message.ExtendedTextMessage.GetText())

	if len(e.Message.ExtendedTextMessage.GetText()) != 0 {
		return MessageExtended
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