package main

import (
	//"log"
	"time"
	ara "github.com/solher/arangolite"
	//valid "github.com/asaskevich/govalidator"
)

// User struct
type User struct {
	ara.Document // Must include arango Document in every struct you want to save id, key, rev after saving it
	Username string `unique:"users" required:"-"`
	Email string `unique:"users" required:"-"`
	Password string `required:"-"`
	Age int
	Gender string `enum:"M,F"`
	Likes []string // TODO: Change for an array of key => value for "UserId" => bool (dislike option)
	Meets []string // Users already met
	RegistrationDate time.Time
}

/**
 *  Methods needed for Model
 */
/*
// GetKey return the key of user in db
func (user *User) GetKey() string {
	return user.Key
}

// GetCollection returns the collection name
func (user *User) GetCollection() string {
	return "users"
}

// GetError returns error
func (user *User) GetError() (string, bool) {
	// default error bool and messages. Could be any kind of error
	return user.Message, user.Error
}
*/
/**
 *  Pre & Post save hooks
 */
/*
// PreSave func for extra validation
func (user *User) PreSave(c *ara.Context) {
	user.Validate(c)
	return
}

func (user *User) PreUpdate(c *ara.Context) {
	user.Validate(c)
	return
}

func (user *User) Validate() {
	if len(user.Username) < 3 {
		c.Err["Username"] = "too short"
	}
	if valid.IsEmail(user.Email) == false  {
		c.Err["Email"] = "invalid"
	}
	if len(user.Password) < 3 {
		c.Err["Password"] = "too short"
	}
}*/
