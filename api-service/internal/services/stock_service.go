package services

import (
	"fmt"
	"time"

	"api-service/internal/models"
	"api-service/internal/repository"
)

type StockService struct {
	repo *repository.StockRepository
}

func NewStockService() *StockService {
	return &StockService{
		repo: repository.NewStockRepository(),
	}
}

func (s *StockService) AdjustStock(partName, partModel string, delta int, reason string, referenceId *string) error {
	movement := &models.StockMovement{
		PartName:    partName,
		PartModel:   partModel,
		Delta:       delta,
		Reason:      reason,
		ReferenceId: referenceId,
		CreatedAt:   time.Now(),
	}
	return s.repo.CreateMovement(movement)
}

func (s *StockService) GetCurrentStock(partName, partModel string) (int, error) {
	movements, err := s.repo.GetMovements(partName, partModel)
	if err != nil {
		return 0, err
	}
	total := 0
	for _, m := range movements {
		total += m.Delta
	}
	return total, nil
}

func (s *StockService) GetMovements(partName, partModel string) ([]models.StockMovement, error) {
	return s.repo.GetMovements(partName, partModel)
}

func (s *StockService) IsAvailable(partName, partModel string, requiredQty int) (bool, int, error) {
	current, err := s.GetCurrentStock(partName, partModel)
	if err != nil {
		return false, 0, err
	}
	return current >= requiredQty, current, nil
}

func (s *StockService) ReserveForRepair(repairRequestID, partName, partModel string) error {
	return s.AdjustStock(partName, partModel, -1, "repair_reserved", &repairRequestID)
}

func (s *StockService) CompleteRepair(repairRequestID, partName, partModel string) error {
	return s.AdjustStock(partName, partModel, 1, "repair_returned", &repairRequestID)
}

func (s *StockService) ConsumeForRepair(repairRequestID, partName, partModel string) error {
	return s.AdjustStock(partName, partModel, -1, "repair_completed", &repairRequestID)
}

func (s *StockService) GetStockReport() (map[string]int, error) {
	report := make(map[string]int)
	
	movements, err := repository.NewStockRepository().GetMovements("", "")
	if err != nil {
		return nil, err
	}
	
	for _, m := range movements {
		key := fmt.Sprintf("%s|%s", m.PartName, m.PartModel)
		report[key] += m.Delta
	}
	
	return report, nil
}
