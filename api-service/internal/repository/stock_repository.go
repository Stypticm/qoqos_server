package repository

import (
	"api-service/internal/database"
	"api-service/internal/models"
	"gorm.io/gorm"
)

var StockRepo *StockRepository

func init() {
	StockRepo = NewStockRepository()
}

type StockRepository struct {
	db *gorm.DB
}

func NewStockRepository() *StockRepository {
	return &StockRepository{db: database.DB}
}

func (r *StockRepository) GetStock(partName, partModel string) (*models.StockMovement, error) {
	var stock models.StockMovement
	err := r.db.Model(&models.StockMovement{}).
		Where("part_name = ? AND part_model = ?", partName, partModel).
		Order("created_at DESC").
		First(&stock).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *StockRepository) CreateMovement(movement *models.StockMovement) error {
	return r.db.Create(movement).Error
}

func (r *StockRepository) GetMovements(partName, partModel string) ([]models.StockMovement, error) {
	var movements []models.StockMovement
	db := r.db
	if partName != "" {
		db = db.Where(`"partName" = ?`, partName)
	}
	if partModel != "" {
		db = db.Where(`"partModel" = ?`, partModel)
	}
	err := db.Order(`"createdAt" DESC`).Find(&movements).Error
	return movements, err
}
