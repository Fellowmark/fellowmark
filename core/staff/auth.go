package staff

import (
	"context"
	"errors"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/nus-utils/nus-peer-review/loggers"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

func (ur StaffRoute) PasswordHash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*models.Staff)
		hash, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
		if err != nil {
			utils.HandleResponse(w, err.Error(), http.StatusInternalServerError)
		} else {
			user.Password = hash
			ctxWithUser := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctxWithUser))
		}
	})
}

func (ur StaffRoute) StaffLoginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *models.Staff
		input := r.Context().Value("user").(*models.Staff)
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
