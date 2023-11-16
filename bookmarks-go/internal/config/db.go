package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func GetGormDb(config AppConfig, logger *Logger) *gorm.DB {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Db.Host, config.Db.Port, config.Db.UserName, config.Db.Password, config.Db.Database)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: glogger.Default.LogMode(glogger.Info),
	})

	if err != nil {
		logger.Fatal(err)
	}
	return db
}
