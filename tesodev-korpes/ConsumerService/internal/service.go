package internal

import (
	"context"
	_ "tesodev-korpes/OrderService/client"
	"tesodev-korpes/OrderService/internal/types"
	_ "tesodev-korpes/pkg"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetByID(ctx context.Context, id string) (*types.Order, error) {
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// if order nil (404 error)

	return order, nil
}
