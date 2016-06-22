package main

import (
	"log"
	"time"
	"errors"
	"encoding/json"
	ara "github.com/solher/arangolite"
)

// Save the user in database
func (user *User) Save() error {
	if user.RegistrationDate.IsZero() {
		user.RegistrationDate = time.Now()
	}
	var q *ara.Query
	if user.Key == nil {
		rd, _ := user.RegistrationDate.MarshalJSON()
		q = ara.NewQuery(`INSERT {
				Username: %q,
				Email: %q,
				Password: %q,
				Age: %d,
				Gender: %q,
				Likes: %q,
				Meets: %q,
				RegistrationDate: %s
			} IN users`,
			user.Username,
			user.Email,
			user.Password,
			user.Age,
			user.Gender,
			user.Likes,
			user.Meets,
			rd,
			)

	} else {
		q = ara.NewQuery(`UPDATE %q WITH {
				Username: %q,
				Email: %q,
				Password: %q,
				Age: %d,
				Gender: %q,
				Likes: %q,
				Meets: %q
			} IN users`,
			*user.Key,
			user.Username,
			user.Email,
			user.Password,
			user.Age,
			user.Gender,
			user.Likes,
			user.Meets,
			)
	}
	log.Println(q)
	_, err := db.Run(q)
	if err != nil {
		log.Println(err)
		return err
	}
	var users []User
	q = ara.NewQuery(`FOR user IN users FILTER user.Username == %q RETURN user`, user.Username).Cache(true).BatchSize(500)
	resp, err := db.Run(q)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(string(resp))
	err = json.Unmarshal(resp, &users)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(users)
	if len(users) > 0 {
		*user = users[0]
		return nil
	}
	return errors.New("End of process...")
}

// Log the user in app
func (user *User) Login() error {
	q := ara.NewQuery(`FOR user IN users FILTER user.Username == %q RETURN user`, user.Username).Cache(true).BatchSize(500)
	resp, err := db.Run(q)
	if err != nil {
		log.Println(err)
		return err
	}
	var users []User
	err = json.Unmarshal(resp, &users)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(users)
	if len(users) > 0 {
		userTmp := users[0]
		if (userTmp.Password != user.Password) {
			return errors.New("wrong password")
		}
		*user = userTmp
		return nil
	}
	return errors.New("unknown username")
}

// Get all Users from collection
func (user *User) GetAll() (*[]User, error) {
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
	return &users, nil
}
