package submissions

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

type FileserverController struct {
	DB            *gorm.DB
	UploadPath    string
	MaxUploadSize int64
}

func (controller FileserverController) CreateRouters(route *mux.Router) {
	controller.CreatePrivilegedRoute(route.NewRoute().Subrouter())
}

func (controller FileserverController) CreatePrivilegedRoute(route *mux.Router) {
	route.Use(utils.AuthenticationMiddleware())

	controller.CreateDownloadRoute(route.NewRoute().Subrouter())
	controller.CreateUploadRoute(route.NewRoute().Subrouter())
}

func (controller FileserverController) CreateDownloadRoute(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Submission{}))
	route.Use(controller.DownloadPermissionCheck())
	route.Use(controller.UpdateFilePathMiddleware())
	route.HandleFunc("", controller.DownloadHandlerFunc()).Methods(http.MethodGet)
}

func (controller FileserverController) CreateUploadRoute(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Submission{}))
	route.Use(controller.UploadPermissionCheck())
	route.Use(controller.GetSubmissionMiddleware())
	route.Use(controller.UploadMiddleware(controller.UploadPath, controller.MaxUploadSize))
	route.Use(controller.UpdateSubmissionContentFile())
	route.HandleFunc("", controller.StoreUploadLocationInDB(controller.DB)).Methods(http.MethodPost)
}
