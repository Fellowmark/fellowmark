package module

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

type ModuleController struct {
	DB *gorm.DB
}

type BatchEnrollment struct {
	ModuleID  uint `json:"moduleId"`
	StudentID uint `json:"studentId"`
	StudentIDs []uint `json:"studentIds"`
	StudentEmails []string `json:"studentEmails"`
}

type BatchSupervision struct {
	ModuleID  uint `json:"moduleId"`
	StaffID uint `json:"staffId"`
	StaffIDs []uint `json:"staffIds"`
	StaffEmails []string `json:"staffEmails"`
}

func (controller ModuleController) CreateRouters(route *mux.Router) {
	controller.CreatePrivilegedRouter(route.NewRoute().Subrouter())
	controller.GetModulesRoute(route.NewRoute().Subrouter())
	controller.GetEnrollmentsRoute(route.PathPrefix("/enrolls").Subrouter())
	controller.GetSupervisionsRoute(route.PathPrefix("/supervises").Subrouter())

}
func (controller ModuleController) CreatePrivilegedRouter(route *mux.Router) {
	route.Use(utils.AuthenticationMiddleware())

	controller.CreateModuleRouter(route.NewRoute().Subrouter())
	controller.CreateEnrollmentRoute(route.PathPrefix("/enroll").Subrouter())
	controller.CreateSupervisionRoute(route.PathPrefix("/supervise").Subrouter())
	controller.GetStudentEnrolledModules(route.PathPrefix("/enroll").Subrouter())
	controller.GetStaffSupervisions(route.PathPrefix("/supervise").Subrouter())
}

func (controller ModuleController) CreateModuleRouter(route *mux.Router) {
	route.Use(utils.IsStaffMiddleware(controller.DB))
	route.Use(utils.DecodeBodyMiddleware(&models.Module{}))
	route.Use(utils.SanitizeDataMiddleware())
	route.HandleFunc("", controller.ModuleCreateHandleFunc()).Methods(http.MethodPost)
}

func (controller ModuleController) CreateEnrollmentRoute(route *mux.Router) {
	route.Use(utils.IsStaffMiddleware(controller.DB))
	route.Use(utils.DecodeBodyMiddleware(&BatchEnrollment{}))
	route.Use(controller.CheckStaffSupervision())
	route.Use(controller.EnrollmentDataPrepare())
	route.HandleFunc("", controller.EnrollmentCreateHandleFunc()).Methods(http.MethodPost)
}

func (controller ModuleController) CreateSupervisionRoute(route *mux.Router) {
	route.Use(utils.IsStaffMiddleware(controller.DB))
	route.Use(utils.DecodeBodyMiddleware(&BatchSupervision{}))
	route.Use(controller.CheckStaffSupervision())
	route.Use(controller.SupervisionDataPrepare())
	route.HandleFunc("", controller.SupervisionCreateHandleFunc()).Methods(http.MethodPost)
}

func (controller ModuleController) GetModulesRoute(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Module{}))
	route.HandleFunc("", utils.DBGetFromDataParams(controller.DB, &models.Module{}, &[]models.Module{})).Methods(http.MethodGet)
}

func (controller ModuleController) GetEnrollmentsRoute(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Enrollment{}))
	route.HandleFunc("", utils.DBGetFromDataParams(controller.DB, &models.Enrollment{}, &[]models.Enrollment{})).Methods(http.MethodGet)
}

func (controller ModuleController) GetSupervisionsRoute(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Supervision{}))
	route.HandleFunc("", utils.DBGetFromDataParams(controller.DB, &models.Supervision{}, &[]models.Supervision{})).Methods(http.MethodGet)
}

func (controller ModuleController) GetStudentEnrolledModules(route *mux.Router) {
	route.HandleFunc("", controller.GetStudentEnrollmentsHandleFunc()).Methods(http.MethodGet)
}

func (controller ModuleController) GetStaffSupervisions(route *mux.Router) {
	route.HandleFunc("", controller.GetStaffSupervisionsHandleFunc()).Methods(http.MethodGet)
}
