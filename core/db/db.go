package db

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/alexedwards/argon2id"
	"github.com/nus-utils/nus-peer-review/loggers"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
)

func InitDB(databaseUrl string) *gorm.DB {
	connection := GetDatabase(databaseUrl)
	db, _ := connection.DB()
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)
	InitialMigration(connection)
	// ResetDatabase(connection)
	SetupAdmin(connection, &models.Admin{
		Name:     os.Getenv("ADMIN_NAME"),
		Email:    os.Getenv("ADMIN_EMAIL"),
		Password: os.Getenv("ADMIN_PASSWORD"),
	})
	if os.Getenv("RUN_ENV") != "production" {
		InsertDummyData(connection)
	}
	// LogPairings(connection)
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
	pool.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&admin)
}

func InsertDummyData(db *gorm.DB) {
	var result *gorm.DB
	passwordHash := utils.HashString("password")
	students := []models.Student{
		{
			Email:    "e0000000@u.nus.edu",
			Name:     "Student A",
			Password: passwordHash,
		},
		{
			Email:    "e0000001@u.nus.edu",
			Name:     "Student B",
			Password: passwordHash,
		},
		{
			Email:    "e0000002@u.nus.edu",
			Name:     "Student C",
			Password: passwordHash,
		},
		{
			Email:    "e0000003@u.nus.edu",
			Name:     "Student D",
			Password: passwordHash,
		},
		{
			Email:    "e0000004@u.nus.edu",
			Name:     "Student E",
			Password: passwordHash,
		},
		{
			Email:    "e0000005@u.nus.edu",
			Name:     "Student F",
			Password: passwordHash,
		},
		{
			Email:    "e0000006@u.nus.edu",
			Name:     "Student G",
			Password: passwordHash,
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
		Name:      "Lecture 1 Quiz",
		Module:    module,
		GroupSize: 3,
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

	utils.InitializePairings(db, assignment)
	utils.SetNewPairings(db, assignment)

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

// (for testing only) resets DB data
func ResetDatabase(db *gorm.DB) {
	db.Exec("DROP TABLE IF EXISTS grades")
	db.Exec("DROP TABLE IF EXISTS submissions")
	db.Exec("DROP TABLE IF EXISTS pairings")
	db.Exec("DROP TABLE IF EXISTS rubrics")
	db.Exec("DROP TABLE IF EXISTS questions")
	db.Exec("DROP TABLE IF EXISTS assignments")
	db.Exec("DROP TABLE IF EXISTS supervisions")
	db.Exec("DROP TABLE IF EXISTS enrollments")
	db.Exec("DROP TABLE IF EXISTS modules")
	db.Exec("DROP TABLE IF EXISTS students")
	db.Exec("DROP TABLE IF EXISTS staffs")
	db.Exec("DROP TABLE IF EXISTS admins")
	InitialMigration(db)
}
