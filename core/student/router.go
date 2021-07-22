package student

import (
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func StudentRouter(route *mux.Router, db *gorm.DB) {
	route.HandleFunc("/signup", SignUp(db)).Methods(http.MethodPost)
	route.HandleFunc("/login", Login(db)).Methods(http.MethodGet)
}
