package module

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

type ModuleRoute struct {
	DB *gorm.DB
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

func (mr ModuleRoute) GetEnrollments(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Enrollment{}, "enrollment"))
	route.HandleFunc("", utils.DBGetFromData(mr.DB, "enrollments", "enrollment", &[]models.Enrollment{})).Methods(http.MethodGet)
}

func (mr ModuleRoute) GetSupervisions(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Enrollment{}, "enrollment"))
	route.HandleFunc("", utils.DBGetFromData(mr.DB, "enrollments", "enrollment", &[]models.Enrollment{})).Methods(http.MethodGet)
}
