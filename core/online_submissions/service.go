package online_submissions

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
	"net/http"
)

func (controller OnlineSubmissionController) CreateOnlineSubmissionPermissionCheck() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(utils.JWTClaimContextKey).(*models.User)
			onlineSubmission := r.Context().Value(utils.DecodeBodyContextKey).(*models.OnlineSubmission)

			var question models.Question
			controller.DB.Model(&models.Question{}).Where("id = ?", onlineSubmission.QuestionID).Find(&question)

			var assignment models.Assignment
			controller.DB.Model(&models.Assignment{}).Where("id = ?", question.AssignmentID).Find(&assignment)

			if pass := utils.IsEnrolled(*claims, assignment.ModuleID, controller.DB); pass {
				next.ServeHTTP(w, r)
			} else {
				utils.HandleResponse(w, "Unauthorized", http.StatusUnauthorized)
			}
		})
	}
}

func (controller OnlineSubmissionController) UpdateOnlineSubmissionPermissionCheck() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(utils.JWTClaimContextKey).(*models.User)
			onlineSubmission := r.Context().Value(utils.DecodeBodyContextKey).(*models.OnlineSubmission)

			var question models.Question
			controller.DB.Model(&models.Question{}).Where("id = ?", onlineSubmission.QuestionID).Find(&question)

			if pass := claims.ID == onlineSubmission.StudentID || utils.IsPair(*claims, question.AssignmentID, onlineSubmission.StudentID, controller.DB); pass {
				next.ServeHTTP(w, r)
			} else {
				utils.HandleResponse(w, "Unauthorized", http.StatusUnauthorized)
			}
		})
	}
}

func (controller OnlineSubmissionController) CreateOnlineSubmissionHandleFunc() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := r.Context().Value(utils.DecodeBodyContextKey).(*models.OnlineSubmission)
			controller.DB.Model(&models.OnlineSubmission{}).Where("submitted_by = ? AND question_id = ?", data.StudentID, data.QuestionID).FirstOrCreate(data)
			ctxWithPath := context.WithValue(r.Context(), utils.DecodeBodyContextKey, data)
			next.ServeHTTP(w, r.WithContext(ctxWithPath))
		})
	}
}

func (controller OnlineSubmissionController) UpdateOnlineSubmissionHandleFunc() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := r.Context().Value(utils.DecodeBodyContextKey).(*models.OnlineSubmission)
			controller.DB.Model(&models.OnlineSubmission{}).Where("submitted_by = ? AND question_id = ?", data.StudentID, data.QuestionID).Updates(data)
			ctxWithPath := context.WithValue(r.Context(), utils.DecodeBodyContextKey, data)
			next.ServeHTTP(w, r.WithContext(ctxWithPath))
		})
	}
}

func (controller OnlineSubmissionController) SaveContentInDB(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := r.Context().Value(utils.DecodeBodyContextKey).(*models.OnlineSubmission)
		result := db.Model(&models.OnlineSubmission{}).Omit("ID").Create(data)
		if result.Error != nil {
			result = db.Model(&models.OnlineSubmission{}).Omit("ID").Where("submitted_by = ? AND question_id = ?", data.StudentID, data.QuestionID).Updates(data)
			if result.Error != nil {
				utils.HandleResponse(w, "Failed", http.StatusInternalServerError)
			}
		} else {
			utils.HandleResponseWithObject(w, data, http.StatusOK)
		}
	}
}
