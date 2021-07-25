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

// Middleware used by both Login and Signup to decode body
func (ur StaffRoute) DecodeUserJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user models.Staff
		if err := utils.DecodeBody(r.Body, &user); err != nil {
			utils.HandleResponse(w, err.Error(), http.StatusBadRequest)
		} else {
			ctxWithUser := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctxWithUser))
		}
	})
}

// Login middleware
func (ur StaffRoute) EmailCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user models.Staff
		input := r.Context().Value("user").(models.Staff)
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
		input := r.Context().Value("input").(models.Staff)
		user := r.Context().Value("user").(models.Staff)
		isEqual, _ := argon2id.ComparePasswordAndHash(input.Password, user.Password)
		if !isEqual {
			utils.HandleResponse(w, "Incorrect Password", http.StatusUnauthorized)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (ur StaffRoute) Login(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.Staff)
	token, err := utils.GenerateJWT("Staff", user)
	if err != nil {
		utils.HandleResponse(w, "Internal Error", http.StatusInternalServerError)
	} else {
		utils.HandleResponse(w, token, http.StatusOK)
	}
}
