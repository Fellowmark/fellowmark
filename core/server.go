package main

import (
	"net/http"
	"os"
	"time"

	DB "github.com/nus-utils/nus-peer-review/db"
	"github.com/nus-utils/nus-peer-review/loggers"
	"github.com/nus-utils/nus-peer-review/routes"
	"gorm.io/gorm"

	"github.com/gorilla/mux"
)

const (
	DBType = "postgres"
)

func main() {
	loggers.InitLoggers(os.Getenv("RUN_ENV"))
	db := DB.InitDB(DBType, os.Getenv("DATABASE_URL"))
	DB.ResetDatabase(db)
	DB.InsertDummyData(db)
	DB.CloseDB(db)
	InitServer(db)
}

func InitServer(db *gorm.DB) {
	loggers.InfoLogger.Println("Starting server")
	route := mux.NewRouter()
	route.HandleFunc("/signup", routes.Login).Methods(http.MethodPost)
	route.HandleFunc("/login", routes.SignUp).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:         "127.0.0.1:8000",
		Handler:      route,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	loggers.ErrorLogger.Fatal(srv.ListenAndServe())
}
