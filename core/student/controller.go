package student

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

type StudentController struct {
	DB *gorm.DB
}

func (controller StudentController) CreateRouters(route *mux.Router) {
	controller.CreateAuthRouter(route.PathPrefix("/auth").Subrouter())
}

func (controller StudentController) CreateAuthRouter(route *mux.Router) {
	signUpRoute := route.NewRoute().Subrouter()
	signUpRoute.Use(utils.DecodeBodyMiddleware(&models.User{}))
	signUpRoute.Use(utils.SanitizeDataMiddleware())
	signUpRoute.Use(utils.UserPasswordHashMiddleware)
	signUpRoute.HandleFunc("/signup", utils.UserCreateHandleFunc(controller.DB, &models.Student{})).Methods(http.MethodPost)

	loginRoute := route.NewRoute().Subrouter()
	loginRoute.HandleFunc("/login", utils.LoginHandleFunc(controller.DB, utils.ModelDBScope(&models.Student{}))).Methods(http.MethodGet)
}
