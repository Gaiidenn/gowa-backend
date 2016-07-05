package users

import (
	"log"
	"time"
	"errors"
	"encoding/json"
	ara "github.com/solher/arangolite"
	"github.com/Gaiidenn/gowa-backend/database"
)

// Save the user in database
func (user *User) Save() error {
	ok, err := user.availableUsername();
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
		return nil
	}
	return errors.New("End of process...")
}

// Log the user in app
func (user *User) Login() error {
	db := database.GetDB()
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