package staff

import (
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func StaffRouter(route *mux.Router, db *gorm.DB) {
	route.HandleFunc("/login", Login(db)).Methods(http.MethodGet)
}
