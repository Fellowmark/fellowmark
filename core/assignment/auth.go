package assignment

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
)

func (controller AssignmentController) AuthorizedToMutateAssignment(asignmentIdResolver func(r *http.Request) string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assignemtId := asignmentIdResolver(r)
			var assignment models.Assignment
			controller.DB.Model(&models.Assignment{}).Where("id = ?", assignemtId).First(&assignment)
			ctxWithModuleID := context.WithValue(r.Context(), "moduleId", assignment.Module.ID)
			next.ServeHTTP(w, r.WithContext(ctxWithModuleID))
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
