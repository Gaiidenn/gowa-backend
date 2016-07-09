package chats

import (
	"github.com/Gaiidenn/gowa-backend/users"
	"time"
	ara "github.com/solher/arangolite"
)

type Chat struct {
	ara.Document // Must include arango Document in every struct you want to save id, key, rev after saving it
	Users        []*users.User
	Conversation []Message
	CreatedAt    time.Time
}

type Message struct {
	ChatKey   string
	User      *users.User
	Msg       string
	CreatedAt time.Time
}
