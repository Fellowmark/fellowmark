package grading

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
)

func (controller GradingController) GradeAccessPermissionCheck() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(utils.JWTClaimContextKey).(*models.User)

			data := r.Context().Value(utils.DecodeParamsContextKey).(*models.Grade)

			var pairing models.Pairing
			controller.DB.Model(&models.Pairing{}).Where("id = ?", data.PairingID).Find(&pairing)

			var assignment models.Assignment
			controller.DB.Model(&models.Assignment{}).Where("id = ?", pairing.AssignmentID).Find(&assignment)

			if claims.ID == pairing.MarkerID || claims.ID == pairing.StudentID ||
				utils.IsSupervisor(*claims, assignment.ModuleID, controller.DB) {
				next.ServeHTTP(w, r)
			} else {
				utils.HandleResponse(w, "Unauthorized", http.StatusUnauthorized)
			}
		})
	}
}

func (controller GradingController) GradeCreatePermissionCheck() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(utils.JWTClaimContextKey).(*models.User)

			data := r.Context().Value(utils.DecodeBodyContextKey).(*models.Grade)

			var pairing models.Pairing
			controller.DB.Model(&models.Pairing{}).Where("id = ?", data.PairingID).Find(&pairing)

			var assignment models.Assignment
			controller.DB.Model(&models.Assignment{}).Where("id = ?", pairing.AssignmentID).Find(&assignment)

			if claims.ID == pairing.MarkerID || utils.IsSupervisor(*claims, assignment.ModuleID, controller.DB) {
				next.ServeHTTP(w, r)
			} else {
				utils.HandleResponse(w, "Unauthorized", http.StatusUnauthorized)
			}
		})
	}
}

func (controller GradingController) DownloadPermissionCheck() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(utils.JWTClaimContextKey).(*models.User)
			submission := r.Context().Value(utils.DecodeParamsContextKey).(*models.Submission)

			var question models.Question
			controller.DB.Model(&models.Assignment{}).Where("id = ?", submission.QuestionID).Find(&question)

			if pass := utils.IsPair(*claims, question.AssignmentID, submission.StudentID, controller.DB); pass {
				next.ServeHTTP(w, r)
			} else {
				utils.HandleResponse(w, "Unauthorized", http.StatusUnauthorized)
			}
		})
	}
}
