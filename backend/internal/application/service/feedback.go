package service

import (
	"context"

	"github.com/diegoclair/leaderpro/internal/domain"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

type feedbackService struct {
	infra domain.Infrastructure
}

func newFeedbackService(infra domain.Infrastructure) contract.FeedbackApp {
	return &feedbackService{
		infra: infra,
	}
}

func (s *feedbackService) CreateFeedback(ctx context.Context, feedback entity.Feedback) error {
	// TODO: Implement
	return nil
}

func (s *feedbackService) GetFeedbackByUUID(ctx context.Context, feedbackUUID string) (entity.Feedback, error) {
	// TODO: Implement
	return entity.Feedback{}, nil
}

func (s *feedbackService) GetPersonFeedback(ctx context.Context, personUUID string, take, skip int64) ([]entity.Feedback, int64, error) {
	// TODO: Implement
	return []entity.Feedback{}, 0, nil
}

func (s *feedbackService) UpdateFeedback(ctx context.Context, feedbackUUID string, feedback entity.Feedback) error {
	// TODO: Implement
	return nil
}

func (s *feedbackService) DeleteFeedback(ctx context.Context, feedbackUUID string) error {
	// TODO: Implement
	return nil
}

func (s *feedbackService) GetFeedbackSummary(ctx context.Context, personUUID string, period string) (entity.FeedbackSummary, error) {
	// TODO: Implement
	return entity.FeedbackSummary{}, nil
}