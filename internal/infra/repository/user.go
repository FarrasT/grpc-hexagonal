package repository

import (
	"errors"
	"grpc-hexa/internal/core/dto"
	"grpc-hexa/internal/core/port/repository"
	"strings"
)

const (
	duplicateEntryMsg = "Duplicate entry"
	numberRowInserted = 1
)

const (
	insertUserStatement = "INSERT INTO tbl_user ( " +
		"`username`, " +
		"`password`, " +
		"`display_name`, " +
		"`created_at`," +
		"`updated_at`) " +
		"VALUES (?, ?, ?, ?, ?)"
)

type userRepository struct {
	db repository.Database
}

func NewUserRepository(db repository.Database) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u userRepository) Insert(user dto.UserDTO) error {
	result, err := u.db.GetDB().Exec(
		insertUserStatement,
		user.UserName,
		user.Password,
		user.DisplayName,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), duplicateEntryMsg) {
			return repository.DuplicateUser
		}
		return err
	}

	numRow, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if numRow != numberRowInserted {
		return errors.New("failed to insert user")
	}

	return nil
}
