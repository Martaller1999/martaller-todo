package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/Martaller1999/martaller-todo/internal/core/domain"
)

func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperTimeout())
	defer cancel()

	query := `
	INSERT INTO martallertodo.users (full_name, phone_number)
	VALUES ($1, $2)
	RETURNING id, version, full_name, phone_number;
	`

	row := r.pool.QueryRow(ctx, query, user.FullName, user.PhoneNumber)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("ошибка сканирования: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
	)

	return userDomain, nil
}
