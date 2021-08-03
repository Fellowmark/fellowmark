package staff

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/assignment"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

type StaffRoute struct {
	DB *gorm.DB
}

func (ur StaffRoute) CreateRouters(route *mux.Router) {
	ur.CreateAuthRouter(route.PathPrefix("/auth").Subrouter())
	ur.CreatePairingsRouter(route.PathPrefix("/assignment").Subrouter())
	ur.CreatePrivilegedRouter(route.PathPrefix("/module/{moduleId}").Subrouter())
}

func (ur StaffRoute) CreateAuthRouter(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Staff{}, "user"))

	loginRoute := route.NewRoute().Subrouter()
	loginRoute.Use(ur.EmailCheck)
	loginRoute.Use(ur.PasswordCheck)
	loginRoute.HandleFunc("/login", utils.LoginHandleFunc(ur.DB, "Staff", "user")).Methods(http.MethodGet)
}

func (ur StaffRoute) CreatePrivilegedRouter(route *mux.Router) {
	if os.Getenv("RUN_ENV") == "production" {
		route.Use(utils.ValidateJWTMiddleware("Staff", "claims"))
		route.Use(utils.SupervisionCheckMiddleware(ur.DB, "claims", "moduleId"))
	}

	ar := assignment.AssignmentRoute{DB: ur.DB}
	ar.CreateAssigmentRouter(route.PathPrefix("/assignment").Subrouter())
	ur.CreatePairingsRouter(route.PathPrefix("/pairing").Subrouter())
}

func (ur StaffRoute) CreatePairingsRouter(route *mux.Router) {
	utils.DecodeBodyMiddleware(&models.Assignment{}, "assignment")

	assignPairingRoute := route.NewRoute().Subrouter()
	assignPairingRoute.HandleFunc("/pairing/initialize", ur.InitializePairings).Methods(http.MethodPost)
	assignPairingRoute.HandleFunc("/pairing/assign", ur.AssignPairings).Methods(http.MethodPost)
}
