package online_submissions

import (
	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
	"net/http"
)

type OnlineSubmissionController struct {
	DB *gorm.DB
}

func (controller OnlineSubmissionController) CreateRouters(route *mux.Router) {
	controller.CreatePrivilegedRoute(route.NewRoute().Subrouter())
}

func (controller OnlineSubmissionController) CreatePrivilegedRoute(route *mux.Router) {
	route.Use(utils.AuthenticationMiddleware())

	controller.CreateStartOnlineSubmissionRoute(route)
	controller.CreateUpdateOnlineSubmissionRoute(route)
}

func (controller OnlineSubmissionController) CreateStartOnlineSubmissionRoute(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.OnlineSubmission{}))
	route.Use(controller.StartOnlineSubmissionPermissionCheck())
	route.Use(controller.GetOnlineSubmissionMiddleware())
	route.Use(controller.UpdateOnlineContent())
	route.HandleFunc("", controller.SaveContentInDB(controller.DB)).Methods(http.MethodPost)
}

func (controller OnlineSubmissionController) CreateUpdateOnlineSubmissionRoute(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.OnlineSubmission{}))
	route.Use(controller.UpdateOnlineSubmissionPermissionCheck())
	route.Use(controller.GetOnlineSubmissionMiddleware())
	route.Use(controller.UpdateOnlineContent())
	route.HandleFunc("", controller.SaveContentInDB(controller.DB)).Methods(http.MethodPut)
}
