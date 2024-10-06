package repository

import (
	"errors"
	"grpc-hexa/internal/core/dto"
)

var (
	DuplicateUser = errors.New("duplicate user")
)

type UserRepository interface {
	Insert(user dto.UserDTO) error
}
