package module

import (
	"net/http"
	"os"

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
	mr.GetEnrollmentsRoute(route.PathPrefix("/enroll").Subrouter())
	mr.GetSupervisionsRoute(route.PathPrefix("/supervise").Subrouter())

	gr := grading.GradingRoute{DB: mr.DB}
	gr.CreateRouters(route.PathPrefix("/{moduleId}/grade").Subrouter())

	sr := submissions.FileserverRoute{DB: mr.DB, UploadPath: "/app/submissions", MaxUploadSize: 10 * 1024 * 1024}
	sr.CreateRouters(route.PathPrefix("/{moduleId}/submit").Subrouter())
}
func (mr ModuleRoute) CreatePrivilegedRouter(route *mux.Router) {
	if os.Getenv("RUN_ENV") == "production" {
		route.Use(utils.ValidateJWTMiddleware("Admin", "claims"))
	}

	mr.CreateModuleRouter(route.NewRoute().Subrouter())
	mr.CreateEnrollmentRoute(route.PathPrefix("/enroll").Subrouter())
	mr.CreateSupervisionRoute(route.PathPrefix("/supervise").Subrouter())
}

func (mr ModuleRoute) CreateModuleRouter(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Module{}, "module"))
	route.Use(utils.SanitizeDataMiddleware("module"))
	route.HandleFunc("", utils.DBCreateHandleFunc(mr.DB, &models.Module{}, "module", true)).Methods(http.MethodPost)
}

func (mr ModuleRoute) CreateEnrollmentRoute(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Enrollment{}, "enrollment"))
	route.HandleFunc("", utils.DBCreateHandleFunc(mr.DB, &models.Enrollment{}, "enrollment", true)).Methods(http.MethodPost)
}

func (mr ModuleRoute) CreateSupervisionRoute(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Supervision{}, "supervision"))
	route.HandleFunc("", utils.DBCreateHandleFunc(mr.DB, &models.Supervision{}, "supervision", true)).Methods(http.MethodPost)
}

func (mr ModuleRoute) GetModulesRoute(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Module{}, "module"))
	route.HandleFunc("", utils.DBGetFromData(mr.DB, &models.Module{}, "module", &[]models.Module{})).Methods(http.MethodGet)
}

func (mr ModuleRoute) GetEnrollmentsRoute(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Enrollment{}, "enrollment"))
	route.HandleFunc("", utils.DBGetFromData(mr.DB, &models.Enrollment{}, "enrollment", &[]models.Enrollment{})).Methods(http.MethodGet)
}

func (mr ModuleRoute) GetSupervisionsRoute(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Supervision{}, "supervision"))
	route.HandleFunc("", utils.DBGetFromData(mr.DB, &models.Supervision{}, "supervision", &[]models.Supervision{})).Methods(http.MethodGet)
}
