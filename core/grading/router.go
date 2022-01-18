package grading

import (
	"net/http"

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
	gr.GetGradesForStudent(route.PathPrefix("/student").Subrouter())
	gr.GetGradesForMarker(route.PathPrefix("/marker").Subrouter())
	gr.GetGradesForStaff(route.NewRoute().Subrouter())
}

func (gr GradingRoute) CreatePrivilegedRouter(route *mux.Router) {
	route.Use(utils.ValidateJWTMiddleware("Student", &models.Student{}))
	route.Use(utils.EnrollmentCheckMiddleware(gr.DB, func(r *http.Request) string { return mux.Vars(r)["moduleId"] }))

	gr.CreateGradeRouter(route.NewRoute().Subrouter())
}

func (gr GradingRoute) CreateGradeRouter(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Grade{}))
	route.Use(utils.MarkerCheckMiddleware(gr.DB))
	// TODO check if grade is valid (i.e., between min and max mark)

	route.HandleFunc("", utils.DBCreateHandleFunc(gr.DB, &models.Grade{}, true)).Methods(http.MethodPost)
}

func (gr GradingRoute) GetGradesForStudent(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Grade{}))
	route.Use(utils.ValidateJWTMiddleware("Student", &models.Student{}))
	route.Use(utils.EnrollmentCheckMiddleware(gr.DB, func(r *http.Request) string { return mux.Vars(r)["moduleId"] }))
	route.Use(utils.MarkeeCheckMiddleware(gr.DB))

	route.HandleFunc("", utils.DBGetFromDataParams(gr.DB, &models.Grade{}, &[]models.Grade{})).Methods(http.MethodGet)
}

func (gr GradingRoute) GetGradesForMarker(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Grade{}))
	route.Use(utils.ValidateJWTMiddleware("Student", &models.Student{}))
	route.Use(utils.EnrollmentCheckMiddleware(gr.DB, func(r *http.Request) string { return mux.Vars(r)["moduleId"] }))
	route.Use(utils.MarkerCheckMiddleware(gr.DB))

	route.HandleFunc("", utils.DBGetFromDataParams(gr.DB, &models.Grade{}, &[]models.Grade{})).Methods(http.MethodGet)
}

func (gr GradingRoute) GetGradesForStaff(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Grade{}))
	route.Use(utils.ValidateJWTMiddleware("Staff", &models.Staff{}))
	route.Use(utils.SupervisionCheckMiddleware(gr.DB, func(r *http.Request) string { return mux.Vars(r)["moduleId"] }))

	route.HandleFunc("", utils.DBGetFromDataBody(gr.DB, &models.Grade{}, &[]models.Grade{})).Methods(http.MethodGet)
}
