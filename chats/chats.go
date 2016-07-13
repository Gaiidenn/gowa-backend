package chats

import (
	"github.com/Gaiidenn/gowa-backend/users"
)

type Chat struct {
	ID 			 string `json:"id"`
	Users        []*users.User `json:"users"`
	Conversation []Message `json:"conversation"`
	CreatedAt    string `json:"createdAt"`
	Private 	 bool `json:"private"`
}

type Message struct {
	ChatID    string `json:"chatID"`
	UserID    string `json:"userID"`
	Msg       string `json:"msg"`
	CreatedAt string `json:"createdAt"`
}
