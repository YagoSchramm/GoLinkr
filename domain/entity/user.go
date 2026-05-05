package entity

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID   `json:"id,omitempty"`
	Name     string      `json:"name,omitempty"`
	Email    string      `json:"email,omitempty"`
	Password string      `json:"password,omitempty"`
	Links    []Link      `json:"links,omitempty"`
	Usage    []Analytics `json:"usage,omitempty"`
}

type UserCredentials struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}
