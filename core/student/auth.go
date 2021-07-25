package student

import (
	"context"
	"errors"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

// Middleware used by both Login and Signup to decode body
func (ur StudentRoute) DecodeUserJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user models.Student
		if err := utils.DecodeBody(r.Body, &user); err != nil {
			utils.HandleResponse(w, err.Error(), http.StatusBadRequest)
		} else {
			ctxWithUser := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctxWithUser))
		}
	})
}

// Signup middleware
func (ur StudentRoute) SanitizeUserData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := validator.Validate(r.Context().Value("user"))
		if err != nil {
			utils.HandleResponse(w, err.Error(), http.StatusBadRequest)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (ur StudentRoute) PasswordHash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(models.Student)
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

func (ur StudentRoute) SignUp(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.Student)
	result := ur.DB.Create(&user)
	if result.Error != nil {
		utils.HandleResponse(w, "Already Exists", http.StatusBadRequest)
	} else {
		utils.HandleResponse(w, "Sucess", http.StatusOK)
	}
}

// Login middleware
func (ur StudentRoute) EmailCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user models.Student
		input := r.Context().Value("user").(models.Student)
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

func (ur StudentRoute) PasswordCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		input := r.Context().Value("input").(models.Student)
		user := r.Context().Value("user").(models.Student)
		isEqual, _ := argon2id.ComparePasswordAndHash(input.Password, user.Password)
		if !isEqual {
			utils.HandleResponse(w, "Incorrect Password", http.StatusUnauthorized)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (ur StudentRoute) Login(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.Student)
	token, err := utils.GenerateJWT("Student", user)
	if err != nil {
		utils.HandleResponse(w, "Internal Error", http.StatusInternalServerError)
	} else {
		utils.HandleResponse(w, token, http.StatusOK)
	}
}
