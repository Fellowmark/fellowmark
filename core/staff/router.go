package staff

import (
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type StaffRoute struct {
	DB *gorm.DB
}

func (ur StaffRoute) CreateRouters(route *mux.Router) {
	ur.CreateAuthRouter(route.PathPrefix("/auth").Subrouter())
}

func (ur StaffRoute) CreateAuthRouter(route *mux.Router) {
	route.Use(ur.DecodeUserJson)

	loginRoute := route.NewRoute().Subrouter()
	loginRoute.HandleFunc("/login", ur.Login).Methods(http.MethodGet)
	loginRoute.Use(ur.EmailCheck)
	loginRoute.Use(ur.PasswordCheck)
}

func (ur StaffRoute) CreateAssignmentRouter(route *mux.Router) {
	assignPairingRoute := route.NewRoute().Subrouter()
	assignPairingRoute.HandleFunc("/assignment/pairing/reset", ur.ResetPairings).Methods(http.MethodPost)
	assignPairingRoute.HandleFunc("/assignment/pairing/assign", ur.AssignPairings).Methods(http.MethodPost)
}
