package student

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

type StudentRoute struct {
	DB *gorm.DB
}

func (ur StudentRoute) CreateRouters(route *mux.Router) {
	ur.CreateAuthRouter(route.PathPrefix("/auth").Subrouter())
}

func (ur StudentRoute) CreateAuthRouter(route *mux.Router) {
	signUpRoute := route.NewRoute().Subrouter()
	signUpRoute.Use(utils.DecodeBodyMiddleware(&models.User{}))
	signUpRoute.Use(utils.SanitizeDataMiddleware())
	signUpRoute.Use(utils.UserPasswordHashMiddleware)
	signUpRoute.HandleFunc("/signup", utils.UserCreateHandleFunc(ur.DB, &models.Student{})).Methods(http.MethodPost)

	loginRoute := route.NewRoute().Subrouter()
	loginRoute.HandleFunc("/login", utils.LoginHandleFunc(ur.DB, utils.ModelDBScope(&models.Student{}))).Methods(http.MethodGet)
}
