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
				Token: %q,
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
			user.Token,
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
				Token: %q,
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
			user.Token,
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
