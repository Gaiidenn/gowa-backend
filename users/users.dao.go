package users

import (
	"errors"
	"github.com/Gaiidenn/gowa-backend/database"
	//ara "github.com/solher/arangolite"
	"time"
	//"encoding/json"
)

// Save the user in database
func (user *User) Save() error {
	connected := user.Connected
	ok, err := user.AvailableUsername();
	if err != nil {
		//log.Println(err)
		return err
	}
	if !ok {
		return errors.New("username already exists")
	}

	if user.ID == "" {
		return errors.New("No such ID")
	}

	if user.RegistrationDate == "" {
		user.RegistrationDate = time.Now().String()
	}

	db := database.GetDB()

	stmt, err := db.Prepare(`
		MERGE (u:User {id: {0}})
		SET u.username = {1}
		SET u.password = {2}
		SET u.email = {3}
		SET u.age = {4}
		SET u.gender = {5}
		SET u.description = {6}
		SET u.token = {7}
		SET u.registrationDate = {8}
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	//log.Println("Trying to update user : ", user)
	rows, err := stmt.Query(
		user.ID,
		user.Username,
		user.Password,
		user.Email,
		user.Age,
		user.Gender,
		user.Description,
		user.Token,
		user.RegistrationDate,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	user.Connected = connected

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
	//log.Println(user.Token)
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
	//log.Println("Trying to update user : ", user)
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
			//log.Println(err, u)
			return nil, err
		}
		//log.Println(u)
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
	//log.Println(u)
	return &u, nil
}

func CountRegistered() (int, error) {
	db := database.GetDB()
	stmt, err := db.Prepare(`
		MATCH (:User)
		WHERE exists(u.registrationDate) AND u.registrationDate <> ""
		RETURN count(*) AS total
	`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var total int
	for rows.Next() {
		err := rows.Scan(&total)
		if err != nil {
			return 0, err
		}
	}

	return total, nil
}

func (user *User) GetPeopleMet() ([]User, error) {
	db := database.GetDB()
	users := make([]User, 0)
	stmt, err := db.Prepare(`
		MATCH (:User {id:{0}})-[:HAS_CHAT]->(c:Chat {private:true})<-[:HAS_CHAT]-(u:User)
		RETURN
			u.id,
			u.username
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Username)
		if err != nil {
			//log.Println(err, u)
			return nil, err
		}
		//log.Println(u)
		users = append(users, u)
	}

	return users, nil
}
