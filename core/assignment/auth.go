package assignment

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
)

func (ar AssignmentRoute) AuthorizedToMutateAssignment(asignmentIdResolver func(r *http.Request) string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assignemtId := asignmentIdResolver(r)
			var assignment models.Assignment
			ar.DB.Model(&models.Assignment{}).Where("id = ?", assignemtId).First(&assignment)
			ctxWithModuleID := context.WithValue(r.Context(), "moduleId", assignment.Module.ID)
			next.ServeHTTP(w, r.WithContext(ctxWithModuleID))
		})
	}
}

func (ar AssignmentRoute) resolveAssignmentIdFromMuxVar(r *http.Request) string {
	return mux.Vars(r)["assignmentId"]
}

func (ar AssignmentRoute) resolveAssignmentIdFromContext(r *http.Request) string {
	assignment := r.Context().Value("data").(*models.Assignment)
	return strconv.FormatUint(uint64(assignment.ID), 10)
}
