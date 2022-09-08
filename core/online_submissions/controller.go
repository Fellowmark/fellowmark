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

	controller.CreateOnlineSubmissionRoute(route.PathPrefix("/create").Subrouter())
	controller.UpdateOnlineSubmissionRoute(route.PathPrefix("/update").Subrouter())
}

func (controller OnlineSubmissionController) CreateOnlineSubmissionRoute(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.OnlineSubmission{}))
	route.Use(controller.CreateOnlineSubmissionPermissionCheck())
	route.Use(controller.CreateOnlineSubmissionHandleFunc())
	route.HandleFunc("", utils.DBCreateHandleFunc(controller.DB, true)).Methods(http.MethodPost)
}

func (controller OnlineSubmissionController) UpdateOnlineSubmissionRoute(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.OnlineSubmission{}))
	route.Use(controller.UpdateOnlineSubmissionPermissionCheck())
	route.Use(controller.UpdateOnlineSubmissionHandleFunc())
	route.HandleFunc("", controller.SaveContentInDB(controller.DB)).Methods(http.MethodPut)
}
