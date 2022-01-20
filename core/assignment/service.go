package assignment

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
)

func (controller AssignmentController) CreateAssignmentPermissionCheck() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(utils.JWTClaimContextKey).(*models.User)
			assignment := r.Context().Value(utils.DecodeBodyContextKey).(*models.Assignment)
			if pass := utils.IsSupervisor(*claims, assignment.ModuleID, controller.DB); pass {
				next.ServeHTTP(w, r)
			} else {
				utils.HandleResponse(w, "Not a supervisor", http.StatusUnauthorized)
			}
		})
	}
}

func (controller AssignmentController) resolveAssignmentIdFromMuxVar(r *http.Request) string {
	return mux.Vars(r)["assignmentId"]
}

func (controller AssignmentController) resolveAssignmentIdFromContext(r *http.Request) string {
	assignment := r.Context().Value("data").(*models.Assignment)
	return strconv.FormatUint(uint64(assignment.ID), 10)
}
