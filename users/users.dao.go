package users

import (
	"log"
	"errors"
	"github.com/Gaiidenn/gowa-backend/database"
	ara "github.com/solher/arangolite"
	"time"
	"encoding/json"
)

// Save the user in database
func (user *User) Save() error {
	connected := user.Connected
	ok, err := user.AvailableUsername();
	if err != nil {
		log.Println(err)
		return err
	}
	if !ok {
		return errors.New("username already exists")
	}

	db := database.GetDB()
	if user.RegistrationDate.IsZero() {
		user.RegistrationDate = time.Now()
	}
	l, err := json.Marshal(user.Likes)
	if err != nil {
		return err
	}
	m, err := json.Marshal(user.Meets)
	if err != nil {
		return err
	}
	var q *ara.Query
	if user.Key == nil {
		rd, _ := user.RegistrationDate.MarshalJSON()
		q = ara.NewQuery(`FOR i IN 1..1 INSERT {
				Username: %q,
				Email: %q,
				Password: %q,
				Profile: {
					Age: %d,
					Gender: %q,
					Description: %q
				},
				RegistrationDate: %s,
				Likes: %s,
				Meets: %s
			} IN users RETURN NEW`,
			user.Username,
			user.Email,
			user.Password,
			user.Profile.Age,
			user.Profile.Gender,
			user.Profile.Description,
			rd,
			l,
			m,
		).Cache(true).BatchSize(500)

	} else {
		q = ara.NewQuery(`UPDATE %q WITH {
				Username: %q,
				Email: %q,
				Password: %q,
				Profile: {
					Age: %d,
					Gender: %q,
					Description: %q
				},
				Likes: %s,
				Meets: %s
			} IN users RETURN NEW`,
			*user.Key,
			user.Username,
			user.Email,
			user.Password,
			user.Profile.Age,
			user.Profile.Gender,
			user.Profile.Description,
			l,
			m,
		).Cache(true).BatchSize(500)
	}
	log.Println(q)
	resp, err := db.Run(q)
	if err != nil {
		log.Println(err)
		return err
	}
	var tmpUser []User
	err = json.Unmarshal(resp, &tmpUser)
	if err != nil {
		return err
	}
	if (len(tmpUser) > 0) {
		tmpUser[0].Connected = connected
		*user = tmpUser[0]
		return nil
	}
	return errors.New("User.Save: db query returned empty")
}

func (user *User) GetAll() ([]User, error) {
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

	return users, nil
}

func (user *User) getByUsername(username string) (*User, error) {
	db := database.GetDB()
	q := ara.NewQuery(`FOR user IN users FILTER user.Username == %q RETURN user`, user.Username).Cache(true).BatchSize(500)
	resp, err := db.Run(q)
	if err != nil {
		return nil, err
	}
	var users []User
	err = json.Unmarshal(resp, &users)
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return &users[0], nil
	}
	return nil, nil
}
