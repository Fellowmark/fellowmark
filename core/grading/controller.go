package grading

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

type GradingController struct {
	DB *gorm.DB
}

func (controller GradingController) CreateRouters(route *mux.Router) {
	controller.CreatePrivilegedRouter(route.NewRoute().Subrouter())
	controller.GetGradesForStudent(route.PathPrefix("/student").Subrouter())
	controller.GetGradesForMarker(route.PathPrefix("/marker").Subrouter())
	controller.GetGradesForStaff(route.NewRoute().Subrouter())
}

func (controller GradingController) CreatePrivilegedRouter(route *mux.Router) {
	route.Use(utils.AuthenticationMiddleware())
	route.Use(utils.EnrollmentCheckMiddleware(controller.DB, func(r *http.Request) string { return mux.Vars(r)["moduleId"] }))

	controller.CreateGradeRouter(route.NewRoute().Subrouter())
}

func (controller GradingController) CreateGradeRouter(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Grade{}))
	route.Use(utils.MarkerCheckMiddleware(controller.DB))
	// TODO: check if controllerade is valid (i.e., between min and max mark)

	route.HandleFunc("", utils.DBCreateHandleFunc(controller.DB, true)).Methods(http.MethodPost)
}

func (controller GradingController) GetGradesForStudent(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Grade{}))
	route.Use(utils.AuthenticationMiddleware())
	route.Use(utils.EnrollmentCheckMiddleware(controller.DB, func(r *http.Request) string { return mux.Vars(r)["moduleId"] }))
	route.Use(utils.MarkeeCheckMiddleware(controller.DB))

	route.HandleFunc("", utils.DBGetFromDataParams(controller.DB, &models.Grade{}, &[]models.Grade{})).Methods(http.MethodGet)
}

func (controller GradingController) GetGradesForMarker(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Grade{}))
	route.Use(utils.AuthenticationMiddleware())
	route.Use(utils.EnrollmentCheckMiddleware(controller.DB, func(r *http.Request) string { return mux.Vars(r)["moduleId"] }))
	route.Use(utils.MarkerCheckMiddleware(controller.DB))

	route.HandleFunc("", utils.DBGetFromDataParams(controller.DB, &models.Grade{}, &[]models.Grade{})).Methods(http.MethodGet)
}

func (controller GradingController) GetGradesForStaff(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Grade{}))
	route.Use(utils.AuthenticationMiddleware())
	route.Use(utils.SupervisionCheckMiddleware(controller.DB, func(r *http.Request) string { return mux.Vars(r)["moduleId"] }))

	route.HandleFunc("", utils.DBGetFromDataBody(controller.DB, &models.Grade{}, &[]models.Grade{})).Methods(http.MethodGet)
}
