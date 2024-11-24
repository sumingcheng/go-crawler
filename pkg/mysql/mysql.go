package mysql

import (
	"crawler/pkg/config"
	"crawler/pkg/logger"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewDB(cfg config.MySQLConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	logger.Info("Connecting to database",
		"host", cfg.Host,
		"port", cfg.Port,
		"database", cfg.Database,
	)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		logger.Error("Failed to connect to database",
			"error", err,
			"host", cfg.Host,
			"port", cfg.Port,
		)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool parameters
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	logger.Info("Successfully connected to database",
		"host", cfg.Host,
		"port", cfg.Port,
		"maxOpenConns", 20,
		"maxIdleConns", 10,
	)

	return db, nil
}
