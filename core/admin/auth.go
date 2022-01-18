package admin

import (
	"errors"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/nus-utils/nus-peer-review/loggers"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

func (ur AdminRoute) AdminLoginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *models.Admin
		input := r.Context().Value("user").(*models.Admin)
		loggers.InfoLogger.Println(input)
		result := ur.DB.Take(&user, "email = ?", input.Email)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			utils.HandleResponse(w, "Incorrect email", http.StatusUnauthorized)
			return
		}

		isEqual, _ := argon2id.ComparePasswordAndHash(input.Password, user.Password)
		if !isEqual {
			utils.HandleResponse(w, "Incorrect Password", http.StatusUnauthorized)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
