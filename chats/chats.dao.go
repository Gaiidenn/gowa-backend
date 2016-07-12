package chats

import (
	"github.com/Gaiidenn/gowa-backend/users"
	"github.com/Gaiidenn/gowa-backend/database"
	"time"
	"github.com/satori/go.uuid"
	"fmt"
	"log"
)


func OpenChat(users []*users.User) (*Chat, error) {

	var chat Chat;
	if (len(users) == 2) {
		chat.Private = true
	} else {
		chat.Private = false
	}
	db := database.GetDB()

	u1 := uuid.NewV4()
	chat.ID = u1.String()

	cd := time.Now().String()

	query := "MERGE "
	params := make([]interface{}, 0)
	index := 0
	for _, user := range users {
		log.Println(index)
		query += " (chat:Chat)-[:HAS_CHAT]-(:User {id:"+fmt.Printf("%s", index)+"}),"
		params = append(params, user.ID)
		index++
	}
	query = query[0:len(query)-1]
	query += " ON CREATE SET chat.createdAt = {"+fmt.Fprintf("%s", index)+"}"
	index++
	params = append(params, cd)
	query += ", chat.id = {"+string(index)+"}"
	index++
	params = append(params, chat.ID)
	query += ", chat.private = {"+string(index)+"}"
	index++
	params = append(params, chat.Private)
	query += " RETURN chat.id, chat.createdAt, chat.private"
	log.Println(query)
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&chat.ID,
			&chat.CreatedAt,
			&chat.Private,
		)
		if err != nil {
			return nil, err
		}
	}
	return &chat, nil
}

// ---------------- OLD CODE --------------------
/*
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
			Conversation: %s,
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
			Users: %s,
			Conversation: %s
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
*/