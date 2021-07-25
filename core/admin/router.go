package admin

import (
	"net/http"

	"github.com/gorilla/mux"
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
	route.Use(ur.DecodeUserJson)

	loginRoute := route.NewRoute().Subrouter()
	loginRoute.HandleFunc("/login", ur.Login).Methods(http.MethodGet)
	loginRoute.Use(ur.EmailCheck)
	loginRoute.Use(ur.PasswordCheck)
}

func (ur AdminRoute) CreatePrivilegedRouter(route *mux.Router) {
	route.Use(ur.ValidateJWT)
	ur.CreateStaffOptsRouter(route.PathPrefix("/staff").Subrouter())
}

func (ur AdminRoute) CreateStaffOptsRouter(route *mux.Router) {
	createStaffRoute := route.NewRoute().Subrouter()
	createStaffRoute.Use(ur.DecodeStaffJson)
	createStaffRoute.Use(ur.SanitizeStaffData)
	createStaffRoute.Use(ur.PasswordHash)
	createStaffRoute.HandleFunc("/", ur.CreateStaff).Methods(http.MethodPost)
}
