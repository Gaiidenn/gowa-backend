package users

import (
	//"log"
	//valid "github.com/asaskevich/govalidator"
)

// User struct
type User struct {
	Token            string `json:"token"`
	ID               string `json:"id"`
	Username         string `json:"username"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	RegistrationDate string `json:"registrationDate"`
	Age              int `json:"age"`
	Gender           string `json:"gender"`
	Description      string `json:"description"`
	Connected        bool `json:"connected"`
}

type Meet struct {
	UserID string `json:"userID"`
	ChatID string `json:"chatID"`
}

type Like struct {
	UserID   string `json:"userID"`
	Positive bool `json:"positive"`
}