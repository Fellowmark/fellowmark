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
	connection.AutoMigrate(models.User{})
}

func CloseDB(connection *gorm.DB) {
	db := connection.DB()
	db.Close()
}
