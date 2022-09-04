package main

import (
	"context"
	"flag"
	"github.com/nus-utils/nus-peer-review/online_submissions"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/nus-utils/nus-peer-review/admin"
	"github.com/nus-utils/nus-peer-review/assignment"
	DB "github.com/nus-utils/nus-peer-review/db"
	"github.com/nus-utils/nus-peer-review/grading"
	"github.com/nus-utils/nus-peer-review/loggers"
	"github.com/nus-utils/nus-peer-review/module"
	"github.com/nus-utils/nus-peer-review/staff"
	"github.com/nus-utils/nus-peer-review/student"
	"github.com/nus-utils/nus-peer-review/submissions"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"

	"github.com/gorilla/mux"
)

func main() {
	print(os.Getenv("DATABASE_URL"))
	loggers.InitLoggers(os.Getenv("RUN_ENV"))
	db := DB.InitDB(os.Getenv("DATABASE_URL"))
	InitServer(db)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	utils.HandleResponseWithObject(w, "Server is healthy", http.StatusOK)
}

func InitServer(pool *gorm.DB) {
	var wait time.Duration
	flag.DurationVar(
		&wait,
		"graceful-timeout",
		time.Second*15,
		"the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m",
	)
	flag.Parse()

	loggers.InfoLogger.Println("Starting server")

	utils.SchemaDecoder.IgnoreUnknownKeys(true)
	router := mux.NewRouter()

	studentController := student.StudentController{
		DB: pool,
	}

	staffController := staff.StaffController{
		DB: pool,
	}

	adminController := admin.AdminController{
		DB: pool,
	}

	moduleController := module.ModuleController{
		DB: pool,
	}

	assignmentController := assignment.AssignmentController{
		DB: pool,
	}

	submissionController := submissions.FileserverController{
		DB:            pool,
		UploadPath:    "/tmp",
		MaxUploadSize: 30 * 1024 * 1024,
	}

	gradingController := grading.GradingController{
		DB: pool,
	}

	onlineSubmissionController := online_submissions.OnlineSubmissionController{
		DB: pool,
	}

	// gradingController := grading.GradingRoute{
	// 	DB: pool,
	// }

	studentController.CreateRouters(router.PathPrefix("/student").Subrouter())
	staffController.CreateRouters(router.PathPrefix("/staff").Subrouter())
	adminController.CreateRouters(router.PathPrefix("/admin").Subrouter())
	moduleController.CreateRouters(router.PathPrefix("/module").Subrouter())
	assignmentController.CreateRouters(router.PathPrefix("/assignment").Subrouter())
	submissionController.CreateRouters(router.PathPrefix("/submission").Subrouter())
	gradingController.CreateRouters(router.PathPrefix("/grade").Subrouter())
	onlineSubmissionController.CreateRouters(router.PathPrefix("/online_submission").Subrouter())

	router.HandleFunc("/health", healthCheck).Methods(http.MethodGet)
	mux.CORSMethodMiddleware(router)

	srv := &http.Server{
		Addr:         ":5000",
		Handler:      utils.SetHeaders(router),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		ErrorLog:     loggers.ErrorLogger,
	}

	go func() {
		loggers.ErrorLogger.Println(srv.ListenAndServe())
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Immediately release DB connections
	DB.CloseDB(pool)

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// TODO: you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	loggers.InfoLogger.Println("shutting down")
}
