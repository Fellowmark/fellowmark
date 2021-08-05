package module

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
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
}
func (mr ModuleRoute) CreatePrivilegedRouter(route *mux.Router) {
	if os.Getenv("RUN_ENV") == "production" {
		route.Use(utils.ValidateJWTMiddleware("Admin", "claims"))
	}

	mr.CreateModuleRouter(route)
	mr.CreateEnrollmentRoute(route.PathPrefix("/enroll").Subrouter())
	mr.CreateSupervisionRoute(route.PathPrefix("/supervise").Subrouter())
}

func (mr ModuleRoute) CreateModuleRouter(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Module{}, "module"))
	route.Use(utils.SanitizeDataMiddleware("module"))
	route.HandleFunc("", utils.DBCreateHandleFunc(mr.DB, "modules", "module")).Methods(http.MethodPost)
}

func (mr ModuleRoute) CreateEnrollmentRoute(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Enrollment{}, "enrollment"))
	route.HandleFunc("", utils.DBCreateHandleFunc(mr.DB, "enrollments", "enrollment")).Methods(http.MethodPost)
}

func (mr ModuleRoute) CreateSupervisionRoute(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Supervision{}, "supervision"))
	route.HandleFunc("", utils.DBCreateHandleFunc(mr.DB, "supervisions", "supervision")).Methods(http.MethodPost)
}

func (mr ModuleRoute) GetModulesRoute(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Module{}, "module"))
	route.HandleFunc("", utils.DBGetFromData(mr.DB, "modules", "module", &[]models.Module{})).Methods(http.MethodGet)
}

func (mr ModuleRoute) GetEnrollmentsRoute(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Enrollment{}, "enrollment"))
	route.HandleFunc("", utils.DBGetFromData(mr.DB, "enrollments", "enrollment", &[]models.Enrollment{})).Methods(http.MethodGet)
}

func (mr ModuleRoute) GetSupervisionsRoute(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Supervision{}, "supervision"))
	route.HandleFunc("", utils.DBGetFromData(mr.DB, "supervisions", "supervision", &[]models.Supervision{})).Methods(http.MethodGet)
}
