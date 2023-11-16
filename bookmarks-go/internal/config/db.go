package config

import (
	"fmt"
	"github.com/sivaprasadreddy/bookmarks/db"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func GetGormDb(config AppConfig, logger *Logger) *gorm.DB {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Db.Host, config.Db.Port, config.Db.UserName, config.Db.Password, config.Db.Database)
	dbConn, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: glogger.Default.LogMode(glogger.Info),
	})

	if err != nil {
		logger.Fatalf("Error connecting to database: %v", err)
	}
	applyDbMigrations(config, logger)
	return dbConn
}

func applyDbMigrations(config AppConfig, logger *Logger) {
	d, err := iofs.New(db.MigrationsFS, "migrations")
	if err != nil {
		logger.Fatalf("Error while loading db migrations from sources: %v", err)
	}
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.Db.UserName, config.Db.Password, config.Db.Host, config.Db.Port, config.Db.Database)
	m, err := migrate.NewWithSourceInstance("iofs", d, databaseURL)
	if err != nil {
		logger.Fatalf("Error while loading db migrations: %v", err)
	}
	err = m.Up()
	if err != nil {
		logger.Fatalf("Error while applying db migrations: %v", err)
	}
	logger.Infof("Database migrations applied successfully")
}
