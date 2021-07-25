package admin

import (
	"context"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gopkg.in/validator.v2"
)

// Middleware used by both Login and Signup to decode body
func (ur AdminRoute) DecodeStaffJson(next http.Handler) http.Handler {
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

// Signup middleware
func (ur AdminRoute) SanitizeStaffData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := validator.Validate(r.Context().Value("user"))
		if err != nil {
			utils.HandleResponse(w, err.Error(), http.StatusBadRequest)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (ur AdminRoute) PasswordHash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(models.Staff)
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

func (ur AdminRoute) CreateStaff(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.Staff)
	result := ur.DB.Create(&user)
	if result.Error != nil {
		utils.HandleResponse(w, "Already Exists", http.StatusBadRequest)
	} else {
		utils.HandleResponse(w, "Sucess", http.StatusOK)
	}
}
