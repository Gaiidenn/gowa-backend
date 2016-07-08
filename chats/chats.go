package chats

import (
	"github.com/Gaiidenn/gowa-backend/users"
	"time"
	ara "github.com/solher/arangolite"
)

type Chat struct {
	ara.Document // Must include arango Document in every struct you want to save id, key, rev after saving it
	users        []*users.User
	msg          chan *Message
	conversation []Message
	createdAt 	 time.Time
}

type Message struct {
	user      *users.User
	msg       string
	createdAt time.Time
}
