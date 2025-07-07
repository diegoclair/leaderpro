package service

import (
	"context"

	"github.com/diegoclair/leaderpro/internal/domain"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

type userService struct {
	infra domain.Infrastructure
}

func newUserService(infra domain.Infrastructure) contract.UserApp {
	return &userService{
		infra: infra,
	}
}

func (s *userService) CreateUser(ctx context.Context, user entity.User) error {
	// TODO: Implement
	return nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	// TODO: Implement
	return entity.User{}, nil
}

func (s *userService) GetUserByUUID(ctx context.Context, userUUID string) (entity.User, error) {
	// TODO: Implement
	return entity.User{}, nil
}

func (s *userService) GetLoggedUser(ctx context.Context) (entity.User, error) {
	// TODO: Implement
	return entity.User{}, nil
}

func (s *userService) GetLoggedUserID(ctx context.Context) (int64, error) {
	// TODO: Implement
	return 0, nil
}

func (s *userService) UpdateUser(ctx context.Context, userUUID string, user entity.User) error {
	// TODO: Implement
	return nil
}