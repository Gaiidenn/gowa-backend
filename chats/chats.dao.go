package chats

import (
	"github.com/Gaiidenn/gowa-backend/users"
	"github.com/Gaiidenn/gowa-backend/database"
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

	query := `
		MERGE (u:User {id:{0}})
		SET u.username = {1}
		SET u.token = {2}
		MERGE (v:User {id:{3}})
		SET v.username = {4}
		SET v.token = {5}
		MERGE (u)-[:HAS_CHAT]->(chat:Chat {private:true})<-[:HAS_CHAT]-(v)
		ON CREATE SET chat.id = {6}
		ON CREATE SET chat.createdAt = timestamp()
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

	query := `
		MATCH (c:Chat {id:{0}})
		MERGE (u:User {id:{1}})
		SET u.username = {2}
		SET u.token = {3}
		CREATE (c)-[:CONTAINS]->(m:Message {msg:{4}, createdAt:timestamp()})
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
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}
