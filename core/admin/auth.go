package admin

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

// Middleware used by both Login and Signup to decode body
func (ur AdminRoute) DecodeUserJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user models.Admin
		if err := utils.DecodeBody(r.Body, &user); err != nil {
			utils.HandleResponse(w, err.Error(), http.StatusBadRequest)
		} else {
			ctxWithUser := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctxWithUser))
		}
	})
}

// Login middleware
func (ur AdminRoute) EmailCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user models.Admin
		input := r.Context().Value("user").(models.Admin)
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

func (ur AdminRoute) PasswordCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		input := r.Context().Value("input").(models.Admin)
		user := r.Context().Value("user").(models.Admin)
		isEqual, _ := argon2id.ComparePasswordAndHash(input.Password, user.Password)
		if !isEqual {
			utils.HandleResponse(w, "Incorrect Password", http.StatusUnauthorized)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (ur AdminRoute) Login(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.Admin)
	token, err := utils.GenerateJWT("Admin", user)
	if err != nil {
		utils.HandleResponse(w, "Internal Error", http.StatusInternalServerError)
	} else {
		utils.HandleResponse(w, token, http.StatusOK)
	}
}

// Middleware used for privileged operations
func (ur AdminRoute) ValidateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			utils.HandleResponse(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Split(authHeader, "Bearer ")[1]
		claims, err := utils.ParseJWT(tokenString)
		if err != nil || claims.Role != "Admin" {
			utils.HandleResponse(w, err.Error(), http.StatusUnauthorized)
		} else {
			ctxWithUser := context.WithValue(r.Context(), "claims", claims)
			next.ServeHTTP(w, r.WithContext(ctxWithUser))
		}
	})
}
