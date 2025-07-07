package service

import (
	"context"

	"github.com/diegoclair/leaderpro/internal/domain"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

type oneOnOneService struct {
	infra domain.Infrastructure
}

func newOneOnOneService(infra domain.Infrastructure) contract.OneOnOneApp {
	return &oneOnOneService{
		infra: infra,
	}
}

func (s *oneOnOneService) CreateOneOnOne(ctx context.Context, oneOnOne entity.OneOnOne) error {
	// TODO: Implement
	return nil
}

func (s *oneOnOneService) GetOneOnOneByUUID(ctx context.Context, oneOnOneUUID string) (entity.OneOnOne, error) {
	// TODO: Implement
	return entity.OneOnOne{}, nil
}

func (s *oneOnOneService) GetPersonOneOnOnes(ctx context.Context, personUUID string, take, skip int64) ([]entity.OneOnOne, int64, error) {
	// TODO: Implement
	return []entity.OneOnOne{}, 0, nil
}

func (s *oneOnOneService) GetManagerOneOnOnes(ctx context.Context, take, skip int64) ([]entity.OneOnOne, int64, error) {
	// TODO: Implement
	return []entity.OneOnOne{}, 0, nil
}

func (s *oneOnOneService) UpdateOneOnOne(ctx context.Context, oneOnOneUUID string, oneOnOne entity.OneOnOne) error {
	// TODO: Implement
	return nil
}

func (s *oneOnOneService) DeleteOneOnOne(ctx context.Context, oneOnOneUUID string) error {
	// TODO: Implement
	return nil
}

func (s *oneOnOneService) GetUpcomingOneOnOnes(ctx context.Context) ([]entity.OneOnOne, error) {
	// TODO: Implement
	return []entity.OneOnOne{}, nil
}

func (s *oneOnOneService) GetOverdueOneOnOnes(ctx context.Context) ([]entity.OneOnOne, error) {
	// TODO: Implement
	return []entity.OneOnOne{}, nil
}