package module

import (
	"net/http"

	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
)

func (mr ModuleRoute) ModuleCreateHandleFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := r.Context().Value(utils.DecodeBodyContextKey).(*models.Module)

		result := mr.DB.Model(data).Create(data)
		if result.Error != nil {
			utils.HandleResponse(w, result.Error.Error(), http.StatusBadRequest)
			return
		}
		utils.HandleResponseWithObject(w, data, http.StatusOK)
	}
}

func (mr ModuleRoute) EnrollmentCreateHandleFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := r.Context().Value(utils.DecodeBodyContextKey).(*models.Enrollment)

		result := mr.DB.Model(data).Create(data)
		if result.Error != nil {
			utils.HandleResponse(w, result.Error.Error(), http.StatusBadRequest)
			return
		}
		utils.HandleResponseWithObject(w, data, http.StatusOK)
	}
}

func (mr ModuleRoute) SupervisionCreateHandleFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := r.Context().Value(utils.DecodeBodyContextKey).(*models.Supervision)

		result := mr.DB.Model(data).Create(data)
		if result.Error != nil {
			utils.HandleResponse(w, result.Error.Error(), http.StatusBadRequest)
			return
		}
		utils.HandleResponseWithObject(w, data, http.StatusOK)
	}
}

func (mr ModuleRoute) GetStudentEnrollmentsHandleFunc() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var modules []models.Module
		user := r.Context().Value(utils.JWTClaimContextKey).(*models.User)
		mr.DB.Joins("inner join enrollments e on modules.id = e.module_id").Where("e.student_id = ?", user.ID).Find(&modules)
		utils.HandleResponseWithObject(w, modules, http.StatusOK)
	})
}

func (mr ModuleRoute) GetStaffSupervisionsHandleFunc() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var modules []models.Module
		user := r.Context().Value(utils.JWTClaimContextKey).(*models.User)
		mr.DB.Joins("inner join supervisions e on modules.id = e.module_id").Where("e.staff_id = ?", user.ID).Find(&modules)
		utils.HandleResponseWithObject(w, modules, http.StatusOK)
	})
}
