package student

import (
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type StudentRoute struct {
	DB *gorm.DB
}

func (ur StudentRoute) CreateRouters(route *mux.Router) {
	ur.CreateAuthRouter(route.PathPrefix("/auth").Subrouter())
}

func (ur StudentRoute) CreateAuthRouter(route *mux.Router) {
	route.Use(ur.DecodeUserJson)

	signUpRoute := route.NewRoute().Subrouter()
	signUpRoute.Use(ur.SanitizeUserData)
	signUpRoute.Use(ur.PasswordHash)
	signUpRoute.HandleFunc("/signup", ur.SignUp).Methods(http.MethodPost)

	loginRoute := route.NewRoute().Subrouter()
	loginRoute.HandleFunc("/login", ur.Login).Methods(http.MethodGet)
	loginRoute.Use(ur.EmailCheck)
	loginRoute.Use(ur.PasswordCheck)
}
