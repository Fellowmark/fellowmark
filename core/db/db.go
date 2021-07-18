package db

import (
	"database/sql"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nus-utils/nus-peer-review/loggers"
	"github.com/nus-utils/nus-peer-review/models"
)

func InitDB(database string, databaseUrl string) *sql.DB {
	connection := GetDatabase(database, databaseUrl)
	db := connection.DB()
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)
	InitialMigration(connection)
	return db
}

func GetDatabase(database string, databaseUrl string) *gorm.DB {
	connection, err := gorm.Open(database, databaseUrl)
	if err != nil {
		loggers.ErrorLogger.Fatalln(err)
	}

	db := connection.DB()

	err = db.Ping()
	if err != nil {
		loggers.ErrorLogger.Fatal("Database not connected")
	}

	loggers.InfoLogger.Println("Connected to database")
	return connection
}

func InitialMigration(connection *gorm.DB) {
	connection.AutoMigrate(models.User{})
}
