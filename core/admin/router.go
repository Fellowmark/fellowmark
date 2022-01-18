package admin

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

type AdminRoute struct {
	DB *gorm.DB
}

func (ur AdminRoute) CreateRouters(route *mux.Router) {
	ur.CreateAuthRouter(route.PathPrefix("/auth").Subrouter())
}

func (ur AdminRoute) CreateAuthRouter(route *mux.Router) {
	loginRoute := route.NewRoute().Subrouter()
	loginRoute.Use(utils.DecodeParamsMiddleware(&models.Admin{}))
	loginRoute.Use(ur.AdminLoginMiddleware)
	loginRoute.HandleFunc("/login", utils.LoginHandleFunc(ur.DB, "Admin")).Methods(http.MethodGet)
}
