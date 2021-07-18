package main

import (
	"net/http"
	"os"
	"time"

	"github.com/nus-utils/nus-peer-review/db"
	"github.com/nus-utils/nus-peer-review/loggers"
	"github.com/nus-utils/nus-peer-review/routes"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	DBType = "postgres"
)

func main() {
	loggers.InitLoggers(os.Getenv("RUN_ENV"))
	db.InitDB(DBType, os.Getenv("DATABASE_URL"))
	InitServer()
}

func InitServer() {
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
