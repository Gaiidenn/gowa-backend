package chats
/*
import (
	"time"

)

func NewMessage(msg *Message) (*Chat, error) {
	//log.Println("NewMessage : ", msg)
	chat, err := GetByKey(msg.ChatKey)
	if err != nil {
		return nil, err
	}
	msg.CreatedAt = time.Now()
	conv := append(chat.Conversation, *msg)
	chat.Conversation = conv
	err = chat.update()
	if err != nil {
		return nil, err
	}
	return chat, err
}*/