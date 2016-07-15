package chats

import (
	"github.com/Gaiidenn/gowa-backend/users"
	"github.com/Gaiidenn/gowa-backend/database"
	"time"
	"github.com/satori/go.uuid"
)


func OpenPrivateChat(user1 *users.User, user2 *users.User) (*Chat, error) {

	var chat Chat;
	chat.Private = true
	chat.Users = make([]*users.User, 0)
	chat.Users = append(chat.Users, user1, user2)

	db := database.GetDB()

	// CREATING / GETTING Chat Node
	u1 := uuid.NewV4()
	chat.ID = u1.String()

	cd := time.Now().String()

	query := `
		MERGE (u:User {id:{0}, username:{1}, token: {2}})
		MERGE (v:User {id:{3}, username:{4}, token: {5}})
		MERGE (u)-[:HAS_CHAT]->(chat:Chat {private:true})<-[:HAS_CHAT]-(v)
		ON CREATE SET chat.id = {6} SET chat.createdAt = {7}
		RETURN chat.id, chat.createdAt
		`
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(
		user1.ID,
		user1.Username,
		user1.Token,
		user2.ID,
		user2.Username,
		user2.Token,
		chat.ID,
		cd,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&chat.ID,
			&chat.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
	}

	// GETTING existing Messages & SETTING conversation
	chat.Conversation = make([]*Message, 0)

	query = `
		MATCH (m:Message)<-[:CONTAINS]-(c:Chat {id:{0}})
		MATCH (u:User)-[:SENT]->(m)
		RETURN m.msg, m.createdAt, u.id, u.username
		ORDER BY m.createdAt ASC
		`
	stmt, err = db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err = stmt.Query(chat.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var msg Message
		var u users.User
		msg.User = &u
		err := rows.Scan(
			&msg.Msg,
			&msg.CreatedAt,
			&msg.User.ID,
			&msg.User.Username,
		)
		if err != nil {
			return nil, err
		}
		chat.Conversation = append(chat.Conversation, &msg)
	}

	return &chat, nil
}

func GetByID(id string) (*Chat, error) {
	var chat Chat;
	chat.ID = id

	db := database.GetDB()

	// GETTING chat
	query := `
		MATCH (c:Chat {id:{0}})
		RETURN c.private, c.createdAt
		`
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&chat.Private,
			&chat.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
	}

	// GETTING chat's users
	query = `
		MATCH (u:User)-[:HAS_CHAT]->(:Chat {id:{0}})
		RETURN u.id, u.username, u.token
		`
	stmt, err = db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err = stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u users.User
		err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.Token,
		)
		if err != nil {
			return nil, err
		}
		chat.Users = append(chat.Users, &u)
	}

	// GETTING chat's messages
	query = `
		MATCH (m:Message)<-[:CONTAINS]-(:Chat {id: {0}})
		MATCH (u:User)-[:SENT]->(m)
		RETURN m.msg, m.createdAt, u.id, u.username
		ORDER BY m.createdAt ASC
		`
	stmt, err = db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err = stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m Message
		var u users.User
		m.ChatID = id
		m.User = &u
		err := rows.Scan(
			&m.Msg,
			&m.CreatedAt,
			&m.User.ID,
			&m.User.Username,
		)
		if err != nil {
			return nil, err
		}
		chat.Conversation = append(chat.Conversation, &m)
	}

	return &chat, nil
}

func (m *Message) Save() error {
	db := database.GetDB()

	cd := time.Now().String()
	if len(m.CreatedAt) == 0 {
		m.CreatedAt = cd
	}

	query := `
		MATCH (c:Chat {id:{0}})
		MERGE (u:User {id:{1}, username:{2}, token:{3}})
		CREATE (c)-[:CONTAINS]->(m:Message {msg:{4}, createdAt:{5}})
		CREATE (u)-[:SENT]->(m)
		`
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(
		m.ChatID,
		m.User.ID,
		m.User.Username,
		m.User.Token,
		m.Msg,
		m.CreatedAt,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
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