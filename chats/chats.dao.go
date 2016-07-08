package chats

import (
	"github.com/Gaiidenn/gowa-backend/users"
	"github.com/Gaiidenn/gowa-backend/database"
	"time"
	"log"
	ara "github.com/solher/arangolite"
	"encoding/json"
)

func (chat *Chat) newChat(users []*users.User) (*Chat, error) {
	db := database.GetDB()
	var conv []Message
	ca, err := time.Now().MarshalJSON()
	if err != nil {
		return nil, err
	}
	var q *ara.Query
	q = ara.NewQuery(`INSERT {
			users: %q,
			conversation: %q,
			createdAt: %s
		} IN chats RETURN NEW`,
		users,
		conv,
		ca,
	)

	log.Println(q)
	resp, err := db.Run(q)
	if err != nil {
		return nil, err
	}
	var tmpChat *Chat
	err = json.Unmarshal(resp, tmpChat)
	if err != nil {
		return nil, err
	}
	log.Println(tmpChat)
	return tmpChat, nil
}

func (chat *Chat) update() error {
	db := database.GetDB()
	var q *ara.Query
	q = ara.NewQuery(`UPDATE %q WITH {
			users: %q
			conversation: %q
		} IN chats`,
		*chat.Key,
		chat.users,
		chat.conversation,
	)
	log.Println(q)
	_, err := db.Run(q)
	if err != nil {
		return err
	}
	return nil
}

func (chat *Chat) getByKey(key string) (*Chat, error) {
	db := database.GetDB()
	var q *ara.Query
	q = ara.NewQuery(`FOR chat IN chats FILTER chat._key == %q RETURN chat`, key)
	log.Println(q)
	resp, err := db.Run(q)
	if err != nil {
		return nil, err
	}
	var tmpChat *Chat
	err = json.Unmarshal(resp, tmpChat)
	if err != nil {
		return nil, err
	}
	return tmpChat, nil
}