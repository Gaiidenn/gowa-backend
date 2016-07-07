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
	var q *ara.Query
	if user.Key == nil {
		rd, _ := user.RegistrationDate.MarshalJSON()
		q = ara.NewQuery(`INSERT {
				Username: %q,
				Email: %q,
				Password: %q,
				Profile: {
					Age: %d,
					Gender: %q,
					Description: %q
				},
				RegistrationDate: %s,
				Likes: %q,
				Meets: %q
			} IN users`,
			user.Username,
			user.Email,
			user.Password,
			user.Profile.Age,
			user.Profile.Gender,
			user.Profile.Description,
			rd,
			user.Likes,
			user.Meets,
		)

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
				Likes: %q,
				Meets: %q
			} IN users`,
			*user.Key,
			user.Username,
			user.Email,
			user.Password,
			user.Profile.Age,
			user.Profile.Gender,
			user.Profile.Description,
			user.Likes,
			user.Meets,
		)
	}
	log.Println(q)
	_, err = db.Run(q)
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
		user.Connected = connected
		return nil
	}
	return errors.New("End of process...")
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
