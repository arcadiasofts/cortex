package database

import (
	"backend/config"
	"context"
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(lc fx.Lifecycle, cfg *config.AppConfig, logger *zap.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s",
			cfg.DB.Host,
			cfg.DB.User,
			cfg.DB.Pass,
			cfg.DB.Name,
			cfg.DB.Port,
			cfg.DB.TimeZone,
		)))

	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			sqlDB, _ := db.DB()
			return sqlDB.Close()
		},
	})
	return db, nil
}
