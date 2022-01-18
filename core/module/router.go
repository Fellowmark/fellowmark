package module

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/grading"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/submissions"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

type ModuleRoute struct {
	DB *gorm.DB
}

func (mr ModuleRoute) CreateRouters(route *mux.Router) {
	mr.CreatePrivilegedRouter(route.NewRoute().Subrouter())
	mr.GetModulesRoute(route.NewRoute().Subrouter())
	mr.GetEnrollmentsRoute(route.PathPrefix("/enrolls").Subrouter())
	mr.GetSupervisionsRoute(route.PathPrefix("/supervises").Subrouter())
	mr.GetStudentEnrolledModules(route.PathPrefix("/enroll").Subrouter())
	mr.GetStaffSupervisions(route.PathPrefix("/supervise").Subrouter())

	gr := grading.GradingRoute{DB: mr.DB}
	gr.CreateRouters(route.PathPrefix("/{moduleId}/grade").Subrouter())

	sr := submissions.FileserverRoute{DB: mr.DB, UploadPath: "/app", MaxUploadSize: 30 * 1024 * 1024}
	sr.CreateRouters(route.PathPrefix("/{moduleId}/submit").Subrouter())
}
func (mr ModuleRoute) CreatePrivilegedRouter(route *mux.Router) {
	route.Use(utils.ValidateJWTMiddleware("Admin", &models.Admin{}))

	mr.CreateModuleRouter(route.NewRoute().Subrouter())
	mr.CreateEnrollmentRoute(route.PathPrefix("/enroll").Subrouter())
	mr.CreateSupervisionRoute(route.PathPrefix("/supervise").Subrouter())
}

func (mr ModuleRoute) CreateModuleRouter(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Module{}))
	route.Use(utils.SanitizeDataMiddleware())
	route.HandleFunc("", utils.DBCreateHandleFunc(mr.DB, &models.Module{}, true)).Methods(http.MethodPost)
}

func (mr ModuleRoute) CreateEnrollmentRoute(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Enrollment{}))
	route.HandleFunc("", utils.DBCreateHandleFunc(mr.DB, &models.Enrollment{}, true)).Methods(http.MethodPost)
}

func (mr ModuleRoute) CreateSupervisionRoute(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Supervision{}))
	route.HandleFunc("", utils.DBCreateHandleFunc(mr.DB, &models.Supervision{}, true)).Methods(http.MethodPost)
}

func (mr ModuleRoute) GetModulesRoute(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Module{}))
	route.HandleFunc("", utils.DBGetFromDataParams(mr.DB, &models.Module{}, &[]models.Module{})).Methods(http.MethodGet)
}

func (mr ModuleRoute) GetEnrollmentsRoute(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Enrollment{}))
	route.HandleFunc("", utils.DBGetFromDataParams(mr.DB, &models.Enrollment{}, &[]models.Enrollment{})).Methods(http.MethodGet)
}

func (mr ModuleRoute) GetSupervisionsRoute(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Supervision{}))
	route.HandleFunc("", utils.DBGetFromDataParams(mr.DB, &models.Supervision{}, &[]models.Supervision{})).Methods(http.MethodGet)
}

func (mr ModuleRoute) GetStudentEnrolledModules(route *mux.Router) {
	route.Use(utils.ValidateJWTMiddleware("Student", &models.Student{}))
	route.HandleFunc("", utils.GetStudentEnrollments(mr.DB)).Methods(http.MethodGet)
}

func (mr ModuleRoute) GetStaffSupervisions(route *mux.Router) {
	route.Use(utils.ValidateJWTMiddleware("Staff", &models.Staff{}))
	route.HandleFunc("", utils.GetStaffSupervisions(mr.DB)).Methods(http.MethodGet)
}
