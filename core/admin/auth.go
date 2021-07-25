package admin

import (
	"errors"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

func SanitizeUserData(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func Login(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqUser models.Admin
		var user models.Admin
		if err := utils.DecodeBody(r.Body, &reqUser); err != nil {
			utils.HandleResponse(w, err.Error(), http.StatusBadRequest)
		}
		result := db.Take(&user, "email = ?", reqUser.Email)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			utils.HandleResponse(w, "Incorrect email", http.StatusUnauthorized)
		} else {
			isEqual, _ := argon2id.ComparePasswordAndHash(reqUser.Password, user.Password)
			if isEqual {
				token, err := utils.GenerateJWT("Admin", reqUser)
				if err != nil {
					utils.HandleResponse(w, "Internal Error", http.StatusInternalServerError)
				}
				utils.HandleResponse(w, token, http.StatusOK)
			}
		}
	}
}
