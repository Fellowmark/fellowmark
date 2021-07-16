package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/routes"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

const (
	DBType = "postgres"
)

func InitLoggers() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger.Println("Logging has begun")
}

func main() {
	InitLoggers()
	InitDB(DBType, os.Getenv("DATABASE_URL"))
	InitServer()
}

func InitServer() {
	route := mux.NewRouter()
	route.HandleFunc("/signup", routes.Login).Methods(http.MethodPost)
	route.HandleFunc("/login", routes.SignUp).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:         "127.0.0.1:8000",
		Handler:      route,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	ErrorLogger.Fatal(srv.ListenAndServe())
}

func InitDB(database string, databaseUrl string) {
	connection := GetDatabase(database, databaseUrl)
	InitialMigration(connection)
	InfoLogger.Println("starting server")
}

func GetDatabase(database string, databaseUrl string) *gorm.DB {
	connection, err := gorm.Open(database, databaseUrl)
	if err != nil {
		ErrorLogger.Fatalln(err)
	}

	db := connection.DB()

	err = db.Ping()
	if err != nil {
		ErrorLogger.Fatal("database not connected")
	}

	InfoLogger.Println("connected to database")
	return connection
}

func InitialMigration(connection *gorm.DB) {
	defer CloseDB(connection)
	connection.AutoMigrate(models.Student{})
	connection.AutoMigrate(models.Staff{})
	connection.AutoMigrate(models.Module{})
	connection.AutoMigrate(models.Assignment{})
	connection.AutoMigrate(models.Question{})
	connection.AutoMigrate(models.Rubric{})
	connection.AutoMigrate(models.Pairing{})
	connection.AutoMigrate(models.Submission{})
	connection.AutoMigrate(models.Grade{})
}

func CloseDB(connection *gorm.DB) {
	db := connection.DB()
	db.Close()
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
		ErrorLogger.Fatal("student entry failed")
	}

	staff := models.Staff{
		Email:    "e1000000@u.nus.edu",
		Name:     "Staff A",
		Password: "password",
	}
	result = db.Create(&staff)
	if result.Error != nil {
		ErrorLogger.Fatal("staff entry failed")
	}

	module := models.Module{
		Code:     "CS2113T",
		Semester: "2122-1",
		Name:     "Software Engineering & OOP",
		Staff:    staff,
	}
	result = db.Create(&module)
	if result.Error != nil {
		ErrorLogger.Fatal("module entry failed")
	}

	assignment := models.Assignment{
		Name:   "Lecture 1 Quiz",
		Module: module,
	}
	result = db.Create(&assignment)
	if result.Error != nil {
		ErrorLogger.Fatal("assignment entry failed")
	}

	question := models.Question{
		QuestionNumber: 1,
		QuestionText:   "What is 2+2?",
		Assignment:     assignment,
	}
	result = db.Create(&question)
	if result.Error != nil {
		ErrorLogger.Fatal("question entry failed")
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
		ErrorLogger.Fatal("rubric entry failed")
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
		ErrorLogger.Fatal("pairing entry failed")
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
		ErrorLogger.Fatal("submission entry failed")
	}
}
