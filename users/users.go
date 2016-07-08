package users

import (
	//"log"
	"time"
	ara "github.com/solher/arangolite"
	//valid "github.com/asaskevich/govalidator"
)

// User struct
type User struct {
	Token string
	ara.Document // Must include arango Document in every struct you want to save id, key, rev after saving it
	Username string `unique:"users" required:"-"`
	Email string `unique:"users" required:"-"`
	Password string `required:"-"`
	RegistrationDate time.Time
	Profile UserProfile
	Connected bool
	Likes []string // TODO: Change for an array of key => value for "UserId" => bool (dislike option)
	Meets []string // Users already met
}

type UserProfile struct {
	Age int
	Gender string
	Description string
}
