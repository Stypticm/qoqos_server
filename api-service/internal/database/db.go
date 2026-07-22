package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"api-service/internal/config"
	"api-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	// Настройка логгера GORM для отладки
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info, // Info покажет все SQL запросы
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      false,
			Colorful:                  true,
		},
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	log.Println("Database connection established")
	
	// Auto-migrate models
	err = DB.AutoMigrate(
		&models.User{},
		&models.Skupka{},
		&models.MarketplaceLot{},
		&models.BlogPost{},
		&models.Order{},
		&models.OrderItem{},
		&models.RepairRequest{},
		&models.OperatorChat{},
		&models.OperatorMessage{},
		&models.QuickLead{},
		&models.TradeInEvaluation{},
		&models.Master{},
		&models.Point{},
		&models.DeviceInspection{},
		&models.Device{},
		&models.MarketPrice{},
		&models.CartItem{},
		&models.FavoriteItem{},
		&models.AuthRequest{},
		&models.PushSubscription{},
		&models.BotAccess{},
		&models.AgentAuditLog{},
		&models.IdempotencyKey{},
	)
	if err != nil {
		log.Printf("Failed to auto-migrate: %v", err)
	} else {
		log.Println("Database auto-migration completed")
	}
}