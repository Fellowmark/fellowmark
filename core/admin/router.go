package admin

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

type AdminRoute struct {
	DB *gorm.DB
}

func (ur AdminRoute) CreateRouters(route *mux.Router) {
	ur.CreatePrivilegedRouter(route.NewRoute().Subrouter())
	ur.CreateAuthRouter(route.PathPrefix("/auth").Subrouter())
}

func (ur AdminRoute) CreateAuthRouter(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Admin{}, "user"))

	loginRoute := route.NewRoute().Subrouter()
	loginRoute.Use(ur.EmailCheck)
	loginRoute.Use(ur.PasswordCheck)
	loginRoute.HandleFunc("/login", utils.LoginHandleFunc(ur.DB, "Admin", "user")).Methods(http.MethodGet)
}

func (ur AdminRoute) CreatePrivilegedRouter(route *mux.Router) {
	if os.Getenv("RUN_ENV") == "production" {
		route.Use(utils.ValidateJWTMiddleware("Admin", "claims", &models.Admin{}))
	}

	ur.CreateStaffOptsRouter(route.PathPrefix("/staff").Subrouter())
}

func (ur AdminRoute) CreateStaffOptsRouter(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Staff{}, "user"))
	route.Use(utils.SanitizeDataMiddleware("user"))
	route.Use(ur.PasswordHash)
	route.HandleFunc("", utils.DBCreateHandleFunc(ur.DB, &models.Staff{}, "user", true)).Methods(http.MethodPost)
}
