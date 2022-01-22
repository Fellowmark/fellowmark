package admin

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

type AdminController struct {
	DB *gorm.DB
}

func (controller AdminController) CreateRouters(route *mux.Router) {
	controller.CreateAuthRouter(route.PathPrefix("/auth").Subrouter())
}

func (controller AdminController) CreateAuthRouter(route *mux.Router) {
	loginRoute := route.NewRoute().Subrouter()
	loginRoute.HandleFunc("/login", utils.LoginHandleFunc(controller.DB, utils.ModelDBScope(&models.Admin{}))).Methods(http.MethodGet)
}
