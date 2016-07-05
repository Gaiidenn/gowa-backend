package users

import (
	"log"
	"encoding/json"
	ara "github.com/solher/arangolite"
	"github.com/Gaiidenn/gowa-backend/database"
)

// Get all Users from collection
func (user *User) GetAll() (*[]User, error) {
	db := database.GetDB()
	q := ara.NewQuery(`FOR user IN users RETURN user`).Cache(true).BatchSize(500)
	log.Println(q)
	resp, err := db.Run(q)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var users []User
	err = json.Unmarshal(resp, &users)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(users)

	// Adding connected users
	for _, cUser := range h.connectedUsers {
		exists := false
		for i, user := range users {
			if cUser.Username == user.Username {
				exists = true
				users[i].Connected = true
				continue
			}
		}
		if !exists {
			users = append(users, *cUser)
		}
	}

	return &users, nil
}

func (user *User) availableUsername() (bool, error) {
	db := database.GetDB()
	q := ara.NewQuery(`FOR user IN users FILTER user.Username == %q RETURN user`, user.Username).Cache(true).BatchSize(500)
	resp, err := db.Run(q)
	if err != nil {
		log.Println(err)
		return false, err
	}
	var users []User
	err = json.Unmarshal(resp, &users)
	if err != nil {
		log.Println(err)
		return false, err
	}
	if len(users) > 0 {
		var key string
		if user.Document.Key != nil {
			key = *user.Document.Key
		} else {
			key = ""
		}
		for _, u := range users {
			if u.Username == user.Username && *u.Document.Key != key {
				return false, nil
			}
		}
	}
	return true, nil
}

func (user *User) readyForSave() bool {
	log.Println(len(user.Username), len(user.Password), len(user.Email))
	if len(user.Username) < 4 {
		return false
	}
	if len(user.Password) < 4 {
		return false
	}
	if len(user.Email) < 4 {
		return false
	}
	return true
}