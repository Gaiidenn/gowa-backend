package chats

import "github.com/Gaiidenn/gowa-backend/users"

func (chat *Chat) NewChat(users []*users.User) (*Chat, error) {
	return chat.newChat(users)
}

func (chat *Chat) Update() error {
	return chat.update()
}