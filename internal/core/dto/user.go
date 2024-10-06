package dto

import "time"

type UserDTO struct {
	UserName    string
	Password    string
	DisplayName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
