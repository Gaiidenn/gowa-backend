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
			id: {0},
			username: {1},
			email: {2},
			password: {3},
			registrationDate: {4},
			age: {5},
			gender: {6},
			description: {7}
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
		MATCH (u:User {
			id = {0}
		}) SET
			username: {1},
			email: {2},
			age: {3},
			gender: {4},
			description: {5}
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	log.Println("Trying to update user : ", user)
	rows, err := stmt.Query(
		user.ID,
		user.Username,
		user.Email,
		user.Age,
		user.Gender,
		user.Description,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (user *User) UpdatePassword() error {
	if user.ID == "" {
		return errors.New("No such ID")
	}

	db := database.GetDB()

	stmt, err := db.Prepare(`
		MATCH (u:User {
			id = {0}
		}) SET
			password: {1}
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	log.Println("Trying to update user : ", user)
	rows, err := stmt.Query(
		user.ID,
		user.Password,
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
			u.username,
			u.id,
			u.age,
			u.gender,
			u.description,
			u.registrationDate
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
		err := rows.Scan(
			&u.Username,
			&u.ID,
			&u.Age,
			&u.Gender,
			&u.Description,
			&u.RegistrationDate,
		)
		if err != nil {
			log.Println(err, u)
			return nil, err
		}
		log.Println(u)
		users = append(users, u)
	}

	return users, nil
}

func (user *User) getByUsername(username string) (*User, error) {
	db := database.GetDB()

	stmt, err := db.Prepare(`
		MATCH (u:User {username: {0}})
		RETURN
			u.id,
			u.username,
			u.password,
			u.email,
			u.age,
			u.gender,
			u.description,
			u.registrationDate
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
}
