package users_service

import (
	"context"
	"fmt"

	"github.com/Martaller1999/martaller-todo/internal/core/domain"
	core_errors "github.com/Martaller1999/martaller-todo/internal/core/errors"
)

func (s *UsersService) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.User, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf(
			"лимит не должен быть отрицательным: %w",
			core_errors.ErrInvalidArgument)
	}
	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf(
			"offset must be non-negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	users, err := s.usersRepository.GetUsers(
		ctx,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("get users from repository")
	}

	return users, nil
}
