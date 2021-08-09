package grading

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

type GradingRoute struct {
	DB *gorm.DB
}

func (gr GradingRoute) CreateRouters(route *mux.Router) {
	gr.CreatePrivilegedRouter(route.NewRoute().Subrouter())
	gr.GetGradesForMarkee(route.PathPrefix("/marker").Subrouter())
	gr.GetGradesForMarkee(route.PathPrefix("/markee").Subrouter())
	gr.GetGradesForMarkee(route.NewRoute().Subrouter())
}

func (gr GradingRoute) CreatePrivilegedRouter(route *mux.Router) {
	if os.Getenv("RUN_ENV") == "production" {
		route.Use(utils.ValidateJWTMiddleware("Student", "claims"))
		route.Use(utils.EnrollmentCheckMiddleware(gr.DB, "claims", "moduleId"))
	}

	gr.CreateGradeRouter(route.NewRoute().Subrouter())
}

func (gr GradingRoute) CreateGradeRouter(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Grade{}, "grade"))
	if os.Getenv("RUN_ENV") == "production" {
		route.Use(utils.MarkerCheckMiddleware(gr.DB, "grade", "claims"))
	}
	route.HandleFunc("", utils.DBCreateHandleFunc(gr.DB, &models.Grade{}, "grade", true)).Methods(http.MethodPost)
}

func (gr GradingRoute) GetGradesForMarkee(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Grade{}, "grade"))
	if os.Getenv("RUN_ENV") == "production" {
		route.Use(utils.ValidateJWTMiddleware("Student", "claims"))
		route.Use(utils.EnrollmentCheckMiddleware(gr.DB, "claims", "moduleId"))
		route.Use(utils.MarkeeCheckMiddleware(gr.DB, "grade", "claims"))
	}

	route.HandleFunc("", utils.DBGetFromData(gr.DB, &models.Grade{}, "grade", &[]models.Grade{})).Methods(http.MethodGet)
}

func (gr GradingRoute) GetGradesForMarker(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Grade{}, "grade"))
	if os.Getenv("RUN_ENV") == "production" {
		route.Use(utils.ValidateJWTMiddleware("Student", "claims"))
		route.Use(utils.EnrollmentCheckMiddleware(gr.DB, "claims", "moduleId"))
		route.Use(utils.MarkerCheckMiddleware(gr.DB, "grade", "claims"))
	}

	route.HandleFunc("", utils.DBGetFromData(gr.DB, &models.Grade{}, "grade", &[]models.Grade{})).Methods(http.MethodGet)
}

func (gr GradingRoute) GetGradesForStaff(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Grade{}, "grade"))
	if os.Getenv("RUN_ENV") == "production" {
		route.Use(utils.ValidateJWTMiddleware("Staff", "claims"))
		route.Use(utils.SupervisionCheckMiddleware(gr.DB, "claims", "moduleId"))
	}

	route.HandleFunc("", utils.DBGetFromData(gr.DB, &models.Grade{}, "grade", &[]models.Grade{})).Methods(http.MethodGet)
}
