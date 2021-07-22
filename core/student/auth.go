package student

import (
	"errors"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/nus-utils/nus-peer-review/loggers"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/routes"
	"gorm.io/gorm"
)

func SignUp(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user *models.Student
		if err := routes.DecodeBody(r.Body, &user); err != nil {
			routes.HandleResponse(w, err.Error(), http.StatusBadRequest)
		}
		hash, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
		if err != nil {
			loggers.ErrorLogger.Println(err)
		}
		user.Password = hash
		result := db.Create(&user)
		if result.Error != nil {
			routes.HandleResponse(w, "Already Exists", http.StatusBadRequest)
		} else {
			routes.HandleResponse(w, "Sucess", http.StatusOK)
		}
	}
}

func Login(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqUser models.Student
		var user models.Student
		if err := routes.DecodeBody(r.Body, &reqUser); err != nil {
			routes.HandleResponse(w, err.Error(), http.StatusBadRequest)
		}
		result := db.Take(&user, "email = ?", reqUser.Email)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			routes.HandleResponse(w, "Incorrect email", http.StatusUnauthorized)
		} else {
			isEqual, _ := argon2id.ComparePasswordAndHash(reqUser.Password, user.Password)
			if isEqual {
				token, err := routes.GenerateJWT("Student", reqUser)
				if err != nil {
					routes.HandleResponse(w, "Internal Error", http.StatusInternalServerError)
				}
				routes.HandleResponse(w, token, http.StatusOK)
			}
		}
	}
}
