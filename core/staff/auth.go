package staff

import (
	"context"
	"errors"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

// Login middleware
func (ur StaffRoute) EmailCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *models.Staff
		input := r.Context().Value("user").(*models.Staff)
		result := ur.DB.Take(&user, "email = ?", input.Email)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			utils.HandleResponse(w, "Incorrect email", http.StatusUnauthorized)
		} else {
			ctxWithUser := context.WithValue(r.Context(), "user", user)
			ctxWithInputAndUser := context.WithValue(ctxWithUser, "input", input)
			next.ServeHTTP(w, r.WithContext(ctxWithInputAndUser))
		}
	})
}

func (ur StaffRoute) PasswordCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		input := r.Context().Value("input").(*models.Staff)
		user := r.Context().Value("user").(*models.Staff)
		isEqual, _ := argon2id.ComparePasswordAndHash(input.Password, user.Password)
		if !isEqual {
			utils.HandleResponse(w, "Incorrect Password", http.StatusUnauthorized)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
