package model

import "time"

type User struct {
	ID        int       `json:"id"`
	NIK       string    `json:"nik"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	NIK      string `json:"nik"`
	FullName string `json:"full_name"`
}

type CreateUserRequest struct {
	NIK      string `json:"nik"`
	FullName string `json:"full_name"`
}

type DeleteUserRequest struct {
	ID string `json:"id"`
}
