package chats

import (
	"github.com/Gaiidenn/gowa-backend/users"
)

type Chat struct {
	ID 			 string
	Users        []*users.User
	Conversation []Message
	CreatedAt    string
	Private 	 bool
}

type Message struct {
	ChatID    string
	User      *users.User
	Msg       string
	CreatedAt string
}
