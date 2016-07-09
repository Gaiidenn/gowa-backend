package chats

import (
	"github.com/Gaiidenn/gowa-backend/users"
	"github.com/Gaiidenn/gowa-backend/database"
	"time"
	"log"
	ara "github.com/solher/arangolite"
	"encoding/json"
	"errors"
)

func NewChat(users []*users.User) (*Chat, error) {
	db := database.GetDB()
	var conv []Message
	ca, err := time.Now().MarshalJSON()
	u, err := json.Marshal(users)
	if err != nil {
		return nil, err
	}
	var q *ara.Query
	q = ara.NewQuery(`INSERT {
			Users: %s,
			Conversation: %q,
			CreatedAt: %s
		} IN chats RETURN NEW`,
		u,
		conv,
		ca,
	).Cache(true).BatchSize(500)

	log.Println(q)
	resp, err := db.Run(q)
	if err != nil {
		return nil, err
	}
	var tmpChat []Chat
	err = json.Unmarshal(resp, &tmpChat)
	if err != nil {
		return nil, err
	}
	if (len(tmpChat) > 0) {
		return &tmpChat[0], nil
	}
	return nil, errors.New("NewChat returned empty array")
}

func (chat *Chat) update() error {
	users, err := json.Marshal(chat.Users)
	if err != nil {
		return err
	}
	conv, err := json.Marshal(chat.Conversation)
	if err != nil {
		return err
	}
	db := database.GetDB()
	var q *ara.Query
	q = ara.NewQuery(`UPDATE %q WITH {
			Users: %q,
			Conversation: %q
		} IN chats`,
		*chat.Key,
		users,
		conv,
	).Cache(true).BatchSize(500)
	log.Println(q)
	_, err = db.Run(q)
	if err != nil {
		return err
	}
	return nil
}

func GetByKey(key string) (*Chat, error) {
	db := database.GetDB()
	log.Println("GetByKey : ", key)
	var q *ara.Query
	q = ara.NewQuery(`FOR chat IN chats FILTER chat._key == %q RETURN chat`, key).Cache(true).BatchSize(500)
	log.Println(q)
	resp, err := db.Run(q)
	if err != nil {
		return nil, err
	}
	log.Println("Chat Response : ", resp)
	var tmpChat []Chat
	err = json.Unmarshal(resp, &tmpChat)
	if err != nil {
		return nil, err
	}
	if (len(tmpChat) > 0) {
		return &tmpChat[0], nil
	}
	return nil, errors.New("Chat.GetByKey: return empty array")
}