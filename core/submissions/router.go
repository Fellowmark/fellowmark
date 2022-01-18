package submissions

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

type FileserverRoute struct {
	DB            *gorm.DB
	UploadPath    string
	MaxUploadSize int64
}

func (fr FileserverRoute) CreateRouters(route *mux.Router) {
	fr.CreatePrivilegedRoute(route.NewRoute().Subrouter())
}

func (fr FileserverRoute) CreatePrivilegedRoute(route *mux.Router) {
	route.Use(utils.ValidateJWTMiddleware("Student", &models.Student{}))
	route.Use(utils.ModulePermCheckMiddleware(fr.DB, "moduleId"))

	fr.CreateDownloadRoute(route.NewRoute().Subrouter())
	fr.CreateUploadRoute(route.NewRoute().Subrouter())
}

func (fr FileserverRoute) CreateDownloadRoute(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Submission{}))
	route.Use(fr.FileAuthMiddleware(false))
	route.Use(fr.UpdateFilePathMiddleware())
	route.HandleFunc("", fr.DownloadHandlerFunc()).Methods(http.MethodGet)
}

func (fr FileserverRoute) CreateUploadRoute(route *mux.Router) {
	route.Use(utils.ValidateJWTMiddleware("Student", &models.Student{}))
	route.Use(utils.EnrollmentCheckMiddleware(fr.DB, func(r *http.Request) string { return mux.Vars(r)["moduleId"] }))
	route.Use(utils.DecodeParamsMiddleware(&models.Submission{}))
	route.Use(fr.FileAuthMiddleware(true))
	route.Use(fr.GetSubmissionMiddleware())
	route.Use(fr.UploadMiddleware(fr.UploadPath, fr.MaxUploadSize))
	route.Use(fr.UpdateSubmissionContentFile())
	route.HandleFunc("", fr.StoreUploadLocationInDB(fr.DB)).Methods(http.MethodPost)
}
