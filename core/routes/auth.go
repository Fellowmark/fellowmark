package routes

import (
	"errors"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/loggers"
	"github.com/nus-utils/nus-peer-review/models"
	"gorm.io/gorm"
)

func StudentAuthRouter(route *mux.Router, db *gorm.DB) {
	route.HandleFunc("/signup", SignUp(db)).Methods(http.MethodPost)
	route.HandleFunc("/login", Login(db)).Methods(http.MethodGet)
}

func SignUp(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user *models.Student
		if err := DecodeBody(r.Body, &user); err != nil {
			HandleResponse(w, err.Error(), http.StatusBadRequest)
		}
		hash, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
		if err != nil {
			loggers.ErrorLogger.Println(err)
		}
		user.Password = hash
		db.Create(&user)
		HandleResponse(w, "Sucess", http.StatusOK)
	}
}

func Login(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqUser models.Student
		var user models.Student
		if err := DecodeBody(r.Body, &reqUser); err != nil {
			HandleResponse(w, err.Error(), http.StatusBadRequest)
		}
		result := db.Take(&user, "email = ?", reqUser.Email)
		HandleResponseWithObject(w, &result, http.StatusOK)
		return
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			HandleResponse(w, "Incorrect email", http.StatusUnauthorized)
		} else {
			isEqual, _ := argon2id.ComparePasswordAndHash(reqUser.Password, user.Password)
			if isEqual {
				token, err := GenerateJWT(reqUser)
				if err != nil {
					HandleResponse(w, "Internal Error", http.StatusInternalServerError)
				}
				HandleResponse(w, token, http.StatusOK)
			}
		}
	}
}
