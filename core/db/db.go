package db

import (
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/alexedwards/argon2id"
	"github.com/nus-utils/nus-peer-review/loggers"
	"github.com/nus-utils/nus-peer-review/models"
)

func InitDB(databaseUrl string) *gorm.DB {
	connection := GetDatabase(databaseUrl)
	db, _ := connection.DB()
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)
	InitialMigration(connection)
	SetupAdmin(connection, &models.Admin{
		Name:     os.Getenv("ADMIN_NAME"),
		Email:    os.Getenv("ADMIN_EMAIL"),
		Password: os.Getenv("ADMIN_PASSOWRD"),
	})
	return connection
}

func GetDatabase(databaseUrl string) *gorm.DB {
	connection, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  databaseUrl,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: false,
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
	connection.AutoMigrate(&models.Module{})
	connection.AutoMigrate(&models.Assignment{})
	connection.AutoMigrate(&models.Question{})
	connection.AutoMigrate(&models.Rubric{})
	connection.AutoMigrate(&models.Pairing{})
	connection.AutoMigrate(&models.Submission{})
	connection.AutoMigrate(&models.Grade{})
	connection.AutoMigrate(&models.Supervision{})
	connection.AutoMigrate(&models.Enrollment{})
	connection.AutoMigrate(&models.Admin{})
}

func CloseDB(connection *gorm.DB) {
	db, _ := connection.DB()
	db.Close()
}

func SetupAdmin(pool *gorm.DB, admin *models.Admin) {
	hash, _ := argon2id.CreateHash(admin.Password, argon2id.DefaultParams)
	admin.Password = hash
	pool.Create(&admin)
}

func InsertDummyData(db *gorm.DB) {
	defer CloseDB(db)
	var result *gorm.DB
	students := []models.Student{
		{
			Email:    "e0000000@u.nus.edu",
			Name:     "Student A",
			Password: "password",
		},
		{
			Email:    "e0000001@u.nus.edu",
			Name:     "Student B",
			Password: "password",
		},
		{
			Email:    "e00000002@u.nus.edu",
			Name:     "Student C",
			Password: "password",
		},
		{
			Email:    "e0000003@u.nus.edu",
			Name:     "Student D",
			Password: "password",
		},
		{
			Email:    "e0000004@u.nus.edu",
			Name:     "Student E",
			Password: "password",
		},
		{
			Email:    "e0000005@u.nus.edu",
			Name:     "Student F",
			Password: "password",
		},
	}
	result = db.Create(&students)
	if result.Error != nil {
		loggers.ErrorLogger.Fatal("student entry failed")
	}

	staff := models.Staff{
		Email:    "e1000000@u.nus.edu",
		Name:     "Staff A",
		Password: "password",
	}
	result = db.Create(&staff)
	if result.Error != nil {
		loggers.ErrorLogger.Fatal("staff entry failed")
	}

	module := models.Module{
		Code:     "CS2113T",
		Semester: "2122-1",
		Name:     "Software Engineering & OOP",
	}
	result = db.Create(&module)
	if result.Error != nil {
		loggers.ErrorLogger.Fatal("module entry failed")
	}

	enrollments := []models.Enrollment{}
	for _, student := range students {
		enrollments = append(enrollments, models.Enrollment{
			Module:  module,
			Student: student,
		})
	}
	result = db.Create(&enrollments)
	if result.Error != nil {
		loggers.ErrorLogger.Fatal("enrollments entry failed")
	}

	supervision := models.Supervision{
		Module: module,
		Staff:  staff,
	}
	result = db.Create(&supervision)
	if result.Error != nil {
		loggers.ErrorLogger.Fatal("supervision entry failed")
	}

	assignment := models.Assignment{
		Name:   "Lecture 1 Quiz",
		Module: module,
	}
	result = db.Create(&assignment)
	if result.Error != nil {
		loggers.ErrorLogger.Fatal("assignment entry failed")
	}

	question := models.Question{
		QuestionNumber: 1,
		QuestionText:   "What is 2+2?",
		Assignment:     assignment,
	}
	result = db.Create(&question)
	if result.Error != nil {
		loggers.ErrorLogger.Fatal("question entry failed")
	}

	rubrics := []models.Rubric{
		{
			Question:    question,
			Criteria:    "Type",
			Description: "1/1: Answer is an integer. 0/1 otherwise",
			MinMark:     0,
			MaxMark:     1,
		},
		{
			Question:    question,
			Criteria:    "Correctness",
			Description: "1/1: Answer is 4. 0/1 otherwise",
			MinMark:     0,
			MaxMark:     1,
		},
	}
	result = db.Create(&rubrics)
	if result.Error != nil {
		loggers.ErrorLogger.Fatal("rubric entry failed")
	}

	pairings := []models.Pairing{}
	// trivial cartesian product of students with each other
	// TODO replace with proper pair assignment
	for idx1, student := range students {
		for idx2, marker := range students {
			if idx1 != idx2 {
				pairings = append(pairings, models.Pairing{
					Assignment: assignment,
					Student:    student,
					Marker:     marker,
				})
			}
		}
	}
	result = db.Create(&pairings)
	if result.Error != nil {
		loggers.ErrorLogger.Fatal("pairing entry failed")
	}

	submission := []models.Submission{
		{
			SubmittedBy: students[0],
			Question:    question,
			Content:     "4",
		},
	}
	result = db.Create(&submission)
	if result.Error != nil {
		loggers.ErrorLogger.Fatal("submission entry failed")
	}
}
