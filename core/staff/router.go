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
