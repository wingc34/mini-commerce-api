package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect initialises a GORM database connection using the provided DSN.
// Connection pool settings are tuned for a small production workload.
func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Show SQL queries in development - set to Silent in production
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// Get the underlying *sql.DB to configure the connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Maximum number of connections kept open in the pool
	sqlDB.SetMaxOpenConns(25)

	// Maximum number of idle connections kept in the pool
	sqlDB.SetMaxIdleConns(10)

	log.Println("Database connected successfully")
	return db, nil
}
