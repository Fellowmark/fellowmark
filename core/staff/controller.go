package staff

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

type StaffController struct {
	DB *gorm.DB
}

func (controller StaffController) CreateRouters(route *mux.Router) {
	controller.CreateAuthRouter(route.PathPrefix("/auth").Subrouter())
	controller.CreatePrivilegedRouter(route.PathPrefix("/module/{moduleId}").Subrouter())
}

func (controller StaffController) CreateAuthRouter(route *mux.Router) {
	loginRoute := route.NewRoute().Subrouter()
	loginRoute.HandleFunc("/login", utils.LoginHandleFunc(controller.DB, utils.ModelDBScope(&models.Staff{}))).Methods(http.MethodGet)

	signUpRoute := route.NewRoute().Subrouter()
	signUpRoute.Use(utils.DecodeBodyMiddleware(&models.User{}))
	signUpRoute.Use(utils.SanitizeDataMiddleware())
	signUpRoute.Use(utils.UserPasswordHashMiddleware)
	signUpRoute.HandleFunc("/signup", utils.UserCreateHandleFunc(controller.DB, &models.Staff{})).Methods(http.MethodPost)
}

func (controller StaffController) CreatePrivilegedRouter(route *mux.Router) {
	route.Use(utils.AuthenticationMiddleware())
	controller.GetPairingsRoute(route.PathPrefix("/pairing").Subrouter())
}

func (controller StaffController) GetPairingsRoute(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Pairing{}))
	route.HandleFunc("", utils.DBGetFromDataParams(controller.DB, &models.Pairing{}, &[]models.Pairing{})).Methods(http.MethodGet)
}
