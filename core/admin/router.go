package admin

import (
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func AdminRouter(route *mux.Router, db *gorm.DB) {
	route.HandleFunc("/login", Login(db)).Methods(http.MethodGet)
}
