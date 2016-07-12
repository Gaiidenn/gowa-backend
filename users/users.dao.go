package users

import (
	"log"
	"errors"
	"github.com/Gaiidenn/gowa-backend/database"
	//ara "github.com/solher/arangolite"
	"time"
	//"encoding/json"
	"github.com/satori/go.uuid"
	"fmt"
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


	log.Println("User Save() : ", user)
	if user.ID == "" {
		user.Create()
	} else {
		user.Update()
	}

	user.Connected = connected

	return nil

	/*
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
	return errors.New("User.Save: db query returned empty")*/
}

func (user *User) Create() error {
	db := database.GetDB()

	u1 := uuid.NewV4()
	id, err := fmt.Printf("%s", u1)
	if err != nil {
		return err
	}
	user.ID = u1.String()

	if user.RegistrationDate == "" {
		user.RegistrationDate = time.Now().String()
	}
	stmt, err := db.Prepare(`
		CREATE (u:User {
			ID: {0},
			Username: {1},
			Email: {2},
			Password: {3},
			RegistrationDate: {4},
			Age: {5},
			Gender: {6},
			Description: {7}
		})
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	log.Println("Trying to create user : ", user)
	rows, err := stmt.Query(
		user.ID,
		user.Username,
		user.Email,
		user.Password,
		user.RegistrationDate,
		user.Age,
		user.Gender,
		user.Description,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (user *User) Update() error {
	if user.ID == "" {
		return errors.New("No such ID")
	}

	db := database.GetDB()

	stmt, err := db.Prepare(`
		MERGE (u:User {
			Username: {0},
			Email: {1},
			Age: {2},
			Gender: {3},
			Description: {4}
		}) WHERE u.ID = {5}
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	log.Println("Trying to update user : ", user)
	rows, err := stmt.Query(
		user.Username,
		user.Email,
		user.Age,
		user.Gender,
		user.Description,
		user.ID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (user *User) GetAll() ([]User, error) {
	db := database.GetDB()

	stmt, err := db.Prepare(`
		MATCH (u:User)
		RETURN
			u.Username,
			u.ID,
			u.Age,
			u.Gender,
			u.Description,
			u.RegistrationDate
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0, 0)
	for rows.Next() {
		var u User
		err := rows.Scan(&u.Username, &u.ID, &u.Age, &u.Gender, &u.Description, &u.RegistrationDate)
		if err != nil {
			log.Println(err, u)
			return nil, err
		}
		log.Println(u)
		users = append(users, u)
	}

	return users, nil
	/*
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

	return users, nil*/
}

func (user *User) getByUsername(username string) (*User, error) {
	db := database.GetDB()

	stmt, err := db.Prepare(`
		MATCH (u:User)
		WHERE u.Username = {0}
		RETURN
			u.ID,
			u.Username,
			u.Password,
			u.Email,
			u.Age,
			u.Gender,
			u.Description,
			u.RegistrationDate
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(user.Username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var u User
	for rows.Next() {
		err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.Password,
			&u.Email,
			&u.Age,
			&u.Gender,
			&u.Description,
			&u.RegistrationDate,
		)
		if err != nil {
			return nil, err
		}
	}
	log.Println(u)
	return &u, nil

	/*
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
	return nil, nil*/
}
