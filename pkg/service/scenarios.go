package service

import (
	"context"

	"github.com/maxoov1/faq-api/pkg/models"
	"github.com/maxoov1/faq-api/pkg/repository"
)

type ScenariosService struct {
	repo repository.Scenarios
}

func NewScenariosService(repo repository.Scenarios) *ScenariosService {
	return &ScenariosService{repo: repo}
}

func (s *ScenariosService) Create(ctx context.Context, scenarioInput ScenariosInput) error {
	scenario := models.Scenario{
		BotID:    scenarioInput.BotID,
		Name:     scenarioInput.Name,
		Triggers: scenarioInput.Triggers,
		Dialogs:  scenarioInput.Dialogs,
	}

	return s.repo.Create(ctx, scenario)
}

func (s *ScenariosService) GetByID(ctx context.Context, id string) (models.Scenario, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ScenariosService) GetAll(ctx context.Context, botID string) ([]models.Scenario, error) {
	return s.repo.GetAll(ctx, botID)
}

func (s *ScenariosService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *ScenariosService) Update(ctx context.Context, scenarioInput ScenariosInput) error {
	scenario := models.Scenario{
		BotID:    scenarioInput.BotID,
		Name:     scenarioInput.Name,
		Triggers: scenarioInput.Triggers,
		Dialogs:  scenarioInput.Dialogs,
	}

	return s.repo.Update(ctx, scenario)
}
