package users

import (
	"log"
	"errors"
	"github.com/Gaiidenn/gowa-backend/database"
	//ara "github.com/solher/arangolite"
	"time"
	//"encoding/json"
	"github.com/satori/go.uuid"
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

	user.ID = uuid.NewV4().String()

	if user.RegistrationDate == "" {
		user.RegistrationDate = time.Now().String()
	}
	stmt, err := db.Prepare(`
		MERGE (u:User {username:{0}, token:{1}})
		ON CREATE SET
			u.id = {2},
			u.email = {3},
			u.password = {4},
			u.registrationDate = {5},
			u.age = {6},
			u.gender = {7},
			u.description = {8}
		ON MATCH SET
			u.email = {9},
			u.password = {10},
			u.registrationDate = {11},
			u.age = {12},
			u.gender = {13},
			u.description = {14}
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	log.Println("Trying to create user : ", user)
	rows, err := stmt.Query(
		user.Username,
		user.Token,
		user.ID,
		user.Email,
		user.Password,
		user.RegistrationDate,
		user.Age,
		user.Gender,
		user.Description,
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
		var tmp interface{}
		err := rows.Scan(&tmp)
		if err != nil {
			return err
		}
		log.Println(tmp)
	}
	return nil
}

func (user *User) UpdateToken() error {
	db := database.GetDB()

	stmt, err := db.Prepare(`
		MATCH (u:User {id: {0}})
		SET u.token = {1}
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	log.Println(user.Token)
	rows, err := stmt.Query(
		user.ID,
		user.Token,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func (user *User) Update() error {
	if user.ID == "" {
		return errors.New("No such ID")
	}

	db := database.GetDB()

	stmt, err := db.Prepare(`
		MATCH (u:User {id: {0}})
		SET u.username = {1}
		SET u.email = {2}
		SET u.age = {3}
		SET u.gender = {4}
		SET u.description = {5}
		SET u.token = {6}
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
		user.Token,
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
		MATCH (u:User {id: {0}})
		SET u.password: {1}
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
		WHERE exists(u.registrationDate)
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
		WHERE exists(u.registrationDate)
		RETURN
			u.id,
			u.token,
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
			&u.Token,
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
