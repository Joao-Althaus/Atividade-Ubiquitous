package user

import "time"

// User is the public representation returned by the API.
type User struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Birthday   string    `json:"birthday"`
	CPF        string    `json:"cpf"`
	InsertedAt time.Time `json:"inserted_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Input contains the fields accepted when creating or updating a user.
type Input struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Birthday string `json:"birthday"`
	CPF      string `json:"cpf"`
}
