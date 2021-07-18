package db

import (
	"github.com/jinzhu/gorm"
	"github.com/nus-utils/nus-peer-review/loggers"
	"github.com/nus-utils/nus-peer-review/models"
)

func InitDB(database string, databaseUrl string) {
	connection := GetDatabase(database, databaseUrl)
	InitialMigration(connection)
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
	defer CloseDB(connection)
	connection.AutoMigrate(models.User{})
}

func CloseDB(connection *gorm.DB) {
	db := connection.DB()
	db.Close()
}
