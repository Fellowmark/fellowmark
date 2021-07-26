package admin

import (
	"context"
	"net/http"

	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gopkg.in/validator.v2"
	"gorm.io/gorm/clause"
)

func (ur AdminRoute) DecodeModuleJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var module models.Module
		if err := utils.DecodeBody(r.Body, &module); err != nil {
			utils.HandleResponse(w, err.Error(), http.StatusBadRequest)
		} else {
			ctxWithUser := context.WithValue(r.Context(), "module", module)
			next.ServeHTTP(w, r.WithContext(ctxWithUser))
		}
	})
}

func (ur AdminRoute) SanitizeModuleData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := validator.Validate(r.Context().Value("module"))
		if err != nil {
			utils.HandleResponse(w, err.Error(), http.StatusBadRequest)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (ur AdminRoute) CreateModule(w http.ResponseWriter, r *http.Request) {
	module := r.Context().Value("module").(models.Staff)
	result := ur.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&module)
	if result.Error != nil {
		utils.HandleResponse(w, "Something went wrong", http.StatusInternalServerError)
	} else {
		utils.HandleResponse(w, "Sucess", http.StatusOK)
	}
}

func (ur AdminRoute) DecodeEnrollmentJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var enrollment models.Enrollment
		if err := utils.DecodeBody(r.Body, &enrollment); err != nil {
			utils.HandleResponse(w, err.Error(), http.StatusBadRequest)
		} else {
			ctxWithUser := context.WithValue(r.Context(), "enrollment", enrollment)
			next.ServeHTTP(w, r.WithContext(ctxWithUser))
		}
	})
}

func (ur AdminRoute) EnrollModuleForStudent(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value("enrollment").(models.Enrollment)
	result := ur.DB.FirstOrCreate(&data)
	if result.Error != nil {
		utils.HandleResponse(w, "Something went wrong", http.StatusInternalServerError)
	} else {
		utils.HandleResponse(w, "Sucess", http.StatusOK)
	}
}
