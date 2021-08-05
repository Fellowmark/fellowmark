package staff

import (
	"net/http"

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
	ur.CreateAssignmentRouter(route.PathPrefix("/assignment").Subrouter())
}

func (ur StaffRoute) CreateAuthRouter(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Staff{}, "user"))

	loginRoute := route.NewRoute().Subrouter()
	loginRoute.Use(ur.EmailCheck)
	loginRoute.Use(ur.PasswordCheck)
	loginRoute.HandleFunc("/login", utils.LoginHandleFunc(ur.DB, "Staff", "user")).Methods(http.MethodGet)
}

func (ur StaffRoute) CreateAssignmentRouter(route *mux.Router) {
	utils.DecodeBodyMiddleware(&models.Assignment{}, "assignment")

	assignPairingRoute := route.NewRoute().Subrouter()
	assignPairingRoute.HandleFunc("/pairing/initialize", ur.InitializePairings).Methods(http.MethodPost)
	assignPairingRoute.HandleFunc("/pairing/assign", ur.AssignPairings).Methods(http.MethodPost)
}
