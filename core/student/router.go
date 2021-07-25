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

	signupRoute := route.HandleFunc("/signup", ur.SignUp).Methods(http.MethodPost)
	signupRoute.Subrouter().Use(ur.SanitizeUserData)
	signupRoute.Subrouter().Use(ur.PasswordHash)

	loginRoute := route.HandleFunc("/login", ur.Login).Methods(http.MethodGet)
	loginRoute.Subrouter().Use(ur.EmailCheck)
	loginRoute.Subrouter().Use(ur.PasswordCheck)
}
