package staff

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

type StaffRoute struct {
	DB *gorm.DB
}

func (ur StaffRoute) CreateRouters(route *mux.Router) {
	ur.CreateAuthRouter(route.PathPrefix("/auth").Subrouter())
	// ur.CreatePairingsRoute(route.PathPrefix("/assignment").Subrouter())
	ur.CreatePrivilegedRouter(route.PathPrefix("/module/{moduleId}").Subrouter())
}

func (ur StaffRoute) CreateAuthRouter(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Staff{}, "user"))

	loginRoute := route.NewRoute().Subrouter()
	loginRoute.Use(ur.EmailCheck)
	loginRoute.Use(ur.PasswordCheck)
	loginRoute.HandleFunc("/login", utils.LoginHandleFunc(ur.DB, "Staff", "user")).Methods(http.MethodGet)
}

func (ur StaffRoute) CreatePrivilegedRouter(route *mux.Router) {
	if os.Getenv("RUN_ENV") == "production" {
		route.Use(utils.ValidateJWTMiddleware("Staff", "claims", &models.Staff{}))
		route.Use(utils.SupervisionCheckMiddleware(ur.DB, "claims", "moduleId"))
	}

	ur.CreatePairingsRoute(route.PathPrefix("/pairing").Subrouter())
	ur.GetPairingsRoute(route.PathPrefix("/pairing").Subrouter())
}

func (ur StaffRoute) CreatePairingsRoute(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Assignment{}, "assignment"))
	route.HandleFunc("/initialize", ur.InitializePairings).Methods(http.MethodPost)
	route.HandleFunc("/assign", ur.AssignPairings).Methods(http.MethodPost)
}

func (ur StaffRoute) GetPairingsRoute(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Pairing{}, "pairing"))
	route.HandleFunc("", utils.DBGetFromData(ur.DB, &models.Pairing{}, "pairing", &[]models.Pairing{})).Methods(http.MethodGet)
}
