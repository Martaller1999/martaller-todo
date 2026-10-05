package users_service

import (
	"context"
	"fmt"

	"github.com/Martaller1999/martaller-todo/internal/core/domain"
)

func (s *UsersService) GetUser(
	ctx context.Context,
	id int,
) (domain.User, error) {
	user, err := s.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("получение пользователя из репозитория: %w", err)
	}
	return user, nil
}
