package staff

import (
	"errors"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/routes"
	"gorm.io/gorm"
)

func StaffAuthRouter(route *mux.Router, db *gorm.DB) {
	route.HandleFunc("/login", Login(db)).Methods(http.MethodGet)
}

func Login(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqUser models.Staff
		var user models.Staff
		if err := routes.DecodeBody(r.Body, &reqUser); err != nil {
			routes.HandleResponse(w, err.Error(), http.StatusBadRequest)
		}
		result := db.Take(&user, "email = ?", reqUser.Email)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			routes.HandleResponse(w, "Incorrect email", http.StatusUnauthorized)
		} else {
			isEqual, _ := argon2id.ComparePasswordAndHash(reqUser.Password, user.Password)
			if isEqual {
				token, err := routes.GenerateJWT("Staff", reqUser)
				if err != nil {
					routes.HandleResponse(w, "Internal Error", http.StatusInternalServerError)
				}
				routes.HandleResponse(w, token, http.StatusOK)
			}
		}
	}
}
