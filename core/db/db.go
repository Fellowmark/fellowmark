package db

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/alexedwards/argon2id"
	"github.com/nus-utils/nus-peer-review/loggers"
	"github.com/nus-utils/nus-peer-review/models"
)

func InitDB(databaseUrl string) *gorm.DB {
	dbUser := os.Getenv("POSTGRES_USER")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DATABASE")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)
	connection := GetDatabase(dsn)
	db, _ := connection.DB()
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)
	InitialMigration(connection)
	SetupAdmin(connection, &models.Admin{
		Name:     os.Getenv("ADMIN_NAME"),
		Email:    os.Getenv("ADMIN_EMAIL"),
		Password: os.Getenv("ADMIN_PASSWORD"),
	})
	return connection
}

func GetDatabase(databaseUrl string) *gorm.DB {
	connection, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  databaseUrl,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: false,
		AllowGlobalUpdate:                        true,
		FullSaveAssociations:                     true,
	})
	if err != nil {
		loggers.ErrorLogger.Fatalln(err)
	}

	db, _ := connection.DB()

	err = db.Ping()
	if err != nil {
		loggers.ErrorLogger.Fatal("Database not connected")
	}

	loggers.InfoLogger.Println("Connected to database")
	return connection
}

func InitialMigration(connection *gorm.DB) {
	connection.AutoMigrate(&models.Student{})
	connection.AutoMigrate(&models.Staff{})
	connection.AutoMigrate(&models.PendingStaff{})
	connection.AutoMigrate(&models.Module{})
	connection.AutoMigrate(&models.Assignment{})
	connection.AutoMigrate(&models.Question{})
	connection.AutoMigrate(&models.Rubric{})
	connection.AutoMigrate(&models.Pairing{})
	connection.AutoMigrate(&models.Submission{})
	connection.AutoMigrate(&models.Grade{})
	connection.AutoMigrate(&models.Supervision{})
	connection.AutoMigrate(&models.Enrollment{})
	connection.AutoMigrate(&models.Assistance{})
	connection.AutoMigrate(&models.Admin{})
}

func CloseDB(connection *gorm.DB) {
	db, _ := connection.DB()
	db.Close()
}

func SetupAdmin(db *gorm.DB, admin *models.Admin) {
	hash, _ := argon2id.CreateHash(admin.Password, argon2id.DefaultParams)
	admin.Password = hash
	db.FirstOrCreate(admin)
	db.Where(models.Admin{Email: admin.Email}).Updates(admin)
}

func LogPairings(db *gorm.DB) {
	result := []models.Pairing{}
	db.Find(&result)
	for _, pairing := range result {
		loggers.InfoLogger.Println(
			fmt.Sprintf("submitter: %v, marker: %v, active: %v", pairing.StudentID, pairing.MarkerID, pairing.Active),
		)
	}
}

func LogStudents(db *gorm.DB) {
	result := []models.Student{}
	db.Find(&result)
	for _, student := range result {
		loggers.InfoLogger.Println(
			fmt.Sprintf("email: %v, name: %v", student.Email, student.Name),
		)
	}
}
